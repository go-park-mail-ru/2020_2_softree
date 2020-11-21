package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"server/src/infrastructure/security"
	"server/src/profileService/profile"
	"time"
)

type UserDBManager struct {
	DB *sql.DB
}

func NewUserDBManager() (*UserDBManager, error) {
	db, err := sql.Open("postgres", viper.GetString("postgres.URL"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &UserDBManager{DB: db}, nil
}

func (h *UserDBManager) GetUserById(ctx context.Context, in *profile.UserID) (*profile.PublicUser, error) {
	user := profile.PublicUser{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, email, password, avatar FROM user_trade WHERE id = $1", in.Id)

	if err = row.Scan(&user.ID, &user.Email, &user.Avatar); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *UserDBManager) CheckExistence(ctx context.Context, in *profile.User) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

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

func (h *UserDBManager) CheckPassword(ctx context.Context, in *profile.User) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userPassword string
	if err := tx.QueryRow("SELECT password FROM user_trade WHERE id = $1", in.ID).Scan(&userPassword); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	check := profile.Check{Existence: bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(in.Password)) == nil}
	return &check, nil
}

func (h *UserDBManager) SaveUser(ctx context.Context, in *profile.User) (*profile.PublicUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var password string
	if password, err = security.MakeShieldedPassword(in.Password); err != nil {
		return nil, err
	}

	var lastID int64
	err = tx.
		QueryRow("INSERT INTO user_trade (email, password) VALUES ($1, $2)  RETURNING id", in.Email, password).
		Scan(&lastID)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	newUser := profile.PublicUser{ID: lastID, Email: in.Email, Avatar: ""}

	return &newUser, nil
}

func (h *UserDBManager) UpdateUserAvatar(ctx context.Context, in *profile.UpdateFields) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE user_trade SET avatar = $1 WHERE id = $2", in.User.Avatar, in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *UserDBManager) UpdateUserPassword(ctx context.Context, in *profile.UpdateFields) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newPassword, err := security.MakeShieldedPassword(in.User.NewPassword)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE user_trade SET password = $1 WHERE id = $2", newPassword, in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *UserDBManager) DeleteUser(ctx context.Context, in *profile.UserID) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM user_trade WHERE id = $1", in.Id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *UserDBManager) GetUserByLogin(ctx context.Context, in *profile.User) (*profile.PublicUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user := profile.User{Email: in.Email}
	row := tx.QueryRow("SELECT id, password, avatar FROM user_trade WHERE email = $1", user.Email)

	if err = row.Scan(&user.ID, &user.Password, &user.Avatar); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		return nil, errors.New("wrong password")
	}

	return &profile.PublicUser{ID: user.ID, Email: user.Email, Avatar: user.Avatar}, nil
}

func (h *UserDBManager) GetUserWatchlist(ctx context.Context, in *profile.UserID) (
	currencies *profile.Currencies, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query("SELECT base_title, currency_title FROM watchlist WHERE user_id = $1", in.Id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var currency profile.Currency
		if err := result.Scan(&currency.Base, &currency.Title); err != nil {
			return nil, err
		}

		currencies.Watchlist = append(currencies.Watchlist, &currency)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(currencies.Watchlist) == 0 {
		currencies.Watchlist = append(currencies.Watchlist, &profile.Currency{Base: "USD", Title: "RUB"})
	}

	return
}
