package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	profile "server/profile/pkg/profile/gen"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserDBManager struct {
	DB *sql.DB
}

func NewUserDBManager(DB *sql.DB) *UserDBManager {
	return &UserDBManager{DB}
}

func (managerDB *UserDBManager) GetUserById(c context.Context, in *profile.UserID) (*profile.PublicUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("GetUserById: %v", err))
		}
	}()

	row := tx.QueryRow("SELECT id, email, avatar FROM user_trade WHERE id = $1", in.Id)

	user := profile.PublicUser{}
	if err = row.Scan(&user.Id, &user.Email, &user.Avatar); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (managerDB *UserDBManager) CheckExistence(ctx context.Context, in *profile.User) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("CheckExistence: %v", err))
		}
	}()

	row := tx.QueryRow("SELECT COUNT(id) FROM user_trade WHERE email = $1", in.Email)

	var exists int
	if err = row.Scan(&exists); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &profile.Check{Existence: exists != 0}, nil
}

func (managerDB *UserDBManager) CheckPassword(ctx context.Context, in *profile.User) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("CheckPassword: %v", err))
		}
	}()

	var userPassword string
	if err := tx.QueryRow("SELECT password FROM user_trade WHERE id = $1", in.Id).Scan(&userPassword); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	check := profile.Check{
		Existence: bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(in.OldPassword)) == nil,
	}
	return &check, nil
}

func (managerDB *UserDBManager) SaveUser(ctx context.Context, in *profile.User) (*profile.PublicUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("SaveUser: %v", err))
		}
	}()

	var lastID int64
	err = tx.
		QueryRow("INSERT INTO user_trade (email, password) VALUES ($1, $2)  RETURNING id", in.Email, in.Password).
		Scan(&lastID)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &profile.PublicUser{Id: lastID, Email: in.Email}, nil
}

func (managerDB *UserDBManager) UpdateUserAvatar(ctx context.Context, in *profile.UpdateFields) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("UpdateUserAvatar: %v", err))
		}
	}()

	_, err = tx.Exec("UPDATE user_trade SET avatar = $1 WHERE id = $2", in.User.Avatar, in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (managerDB *UserDBManager) UpdateUserPassword(ctx context.Context, in *profile.UpdateFields) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("UpdateUserPassword: %v", err))
		}
	}()

	_, err = tx.Exec("UPDATE user_trade SET password = $1 WHERE id = $2", in.User.NewPassword, in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (managerDB *UserDBManager) DeleteUser(ctx context.Context, in *profile.UserID) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("DeleteUser: %v", err))
		}
	}()

	_, err = tx.Exec("DELETE FROM user_trade WHERE id = $1", in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (managerDB *UserDBManager) GetUserByLogin(ctx context.Context, in *profile.User) (*profile.PublicUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("GetUserByLogin: %v", err))
		}
	}()

	var password string
	row := tx.QueryRow("SELECT id, password, avatar FROM user_trade WHERE email = $1", in.Email)

	if err = row.Scan(&in.Id, &password, &in.Avatar); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(password), []byte(in.Password)) != nil {
		return nil, errors.New("wrong password")
	}

	return &profile.PublicUser{Id: in.Id, Email: in.Email, Avatar: in.Avatar}, nil
}

func (managerDB *UserDBManager) GetUserWatchlist(ctx context.Context, in *profile.UserID) (*profile.Currencies, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("GetUserWatchlist: %v", err))
		}
	}()

	result, err := tx.Query("SELECT base_title, currency_title FROM watchlist WHERE user_id = $1", in.Id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var currencies profile.Currencies
	for result.Next() {
		var currency profile.Currency
		if err := result.Scan(&currency.Base, &currency.Title); err != nil {
			return nil, err
		}

		currencies.Currencies = append(currencies.Currencies, &currency)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(currencies.Currencies) == 0 {
		currencies.Currencies = append(currencies.Currencies, &profile.Currency{Base: "USD", Title: "RUB"})
	}

	return &currencies, nil
}
