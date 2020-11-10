package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"server/src/domain/entity"
	"server/src/infrastructure/config"
	"server/src/infrastructure/security"
	"time"
)

type UserDBManager struct {
	DB   *sql.DB
}

func NewUserDBManager() (*UserDBManager, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.UserDatabaseConfig.Host,
		5432,
		config.UserDatabaseConfig.User,
		config.UserDatabaseConfig.Password,
		config.UserDatabaseConfig.Schema)

	db, err := sql.Open("postgres", psqlInfo)

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &UserDBManager{DB: db}, nil
}

func (h *UserDBManager) GetUserById(id uint64) (entity.User, error) {
	user := entity.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, email, password FROM user_trade WHERE id = $1", id)

	err = tx.Commit()
	if err != nil {
		return entity.User{}, err
	}

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return entity.User{}, err
	}

	// #TODO ADD AVATAR

	return user, nil
}

func (h *UserDBManager) SaveUser(user entity.User) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT COUNT(id) FROM user_trade WHERE email = $1", user.Email)

	var exists int
	if err = row.Scan(&exists); err != nil {
		return entity.User{}, err
	}
	if err = tx.Commit(); err != nil {
		return entity.User{}, err
	}

	if exists != 0 {
		return entity.User{}, errors.New("user already exists")
	}

	tx, err = h.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}
	password, err := security.MakeShieldedPassword(user.Password)
	if err != nil {
		return entity.User{}, err
	}

	result, err := tx.Exec(
		"INSERT INTO user_trade (`email`, `password`) VALUES ($1, $2)",
		user.Email,
		password,
	)

	err = tx.Commit()
	if err != nil {
		return entity.User{}, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return entity.User{}, err
	}
	newUser := entity.User{
		ID: uint64(lastID),
		Email: user.Email,
		Password: password,
	}

	return newUser, nil
}

func (h *UserDBManager) UpdateUser(id uint64, user entity.User) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}
	defer tx.Rollback()

	currentUser, err := h.GetUserById(id)
	if err != nil {
		return entity.User{}, err
	}

	if !govalidator.IsNull(user.Avatar) {
		// #TODO
	}

	var newPassword string
	if !govalidator.IsNull(user.OldPassword) && !govalidator.IsNull(user.NewPassword) {
		if bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(user.OldPassword)) == nil {
			newPassword, err = security.MakeShieldedPassword(user.NewPassword)
			if err != nil {
				return entity.User{}, err
			}
			_, err = tx.Exec(
				"UPDATE user_trade SET password = $1 WHERE id = $2",
				newPassword,
				id,
			)
			err = tx.Commit()
			if err != nil {
				return entity.User{}, err
			}
		}
	}
	resUser := entity.User{
		ID: id,
		Email: user.Email,
		Password: newPassword,
		// ADD AVATAR
	}

	return resUser, nil
}

func (h *UserDBManager) DeleteUser(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"DELETE FROM user_trade WHERE id = $1",
		id,
	)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (h *UserDBManager) GetUserByLogin(email string, password string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.User{}, err
	}
	defer tx.Rollback()

	user := entity.User{Email: email}

	row := tx.QueryRow("SELECT id, password FROM user_trade WHERE email = $1", email)

	err = tx.Commit()
	if err != nil {
		return entity.User{}, err
	}

	err = row.Scan(&user.ID, &user.Password)
	if govalidator.IsNull(user.Password) {
		return entity.User{}, errors.New("user does not exist")
	}
	if err != nil {
		return entity.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return entity.User{}, errors.New("wrong password")
	}

	// #TODO ADD AVATAR

	return user, nil
}

func (h *UserDBManager) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, _ := tx.Query(
		"SELECT base_title, currency_title FROM watchlist WHERE user_id = $1",
		id,
	)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	defer result.Close()
	if err := result.Err(); err != nil {
		return nil, err
	}

	currencies := make([]entity.Currency, 0)
	for result.Next() {
		var currency entity.Currency
		if err := result.Scan(&currency.Base, &currency.Title); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}
