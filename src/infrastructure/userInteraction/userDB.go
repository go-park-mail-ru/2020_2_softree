package userInteraction

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"server/src/domain/entity"
	"server/src/infrastructure/config"
	"server/src/infrastructure/security"
)

type UserDBManager struct {
	DB   *sql.DB
}

func NewUserDBManager() (*UserDBManager, error) {
	/*dsn := config.UserDatabaseConfig.User +
		":" + config.UserDatabaseConfig.Password +
		"@" + config.UserDatabaseConfig.Port +
		"/" + config.UserDatabaseConfig.Schema
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"*/

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.UserDatabaseConfig.Host,
		5432,
		config.UserDatabaseConfig.User,
		config.UserDatabaseConfig.Password,
		config.UserDatabaseConfig.Schema)

	/*dsn := "root:1234@tcp(localhost:3306)/tech?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"*/

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

	row := h.DB.QueryRow("SELECT id, email, password FROM user WHERE id = $1", id)

	err := row.Scan(user.ID, user.Email, user.Password)
	if err != nil {
		return entity.User{}, err
	}

	// ADD AVATAR

	return user, nil
}

func (h *UserDBManager) SaveUser(user entity.User) (entity.User, error) {
	row := h.DB.QueryRow("SELECT COUNT(*) FROM user WHERE email = $1", user.Email)

	var exists int
	err := row.Scan(exists)
	if err != nil {
		return entity.User{}, err
	}
	if exists != 0 {
		return entity.User{}, errors.New("user already exists")
	}

	password, err := security.MakeShieldedPassword(user.Password)
	if err != nil {
		return entity.User{}, err
	}

	result, err := h.DB.Exec(
		"INSERT INTO user (`email`, `password`) VALUES ($1, $2)",
		user.Email,
		password,
	)

	//affected, err := result.RowsAffected()
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
	currentUser, err := h.GetUserById(id)
	if err != nil {
		return entity.User{}, err
	}

	if !govalidator.IsNull(user.Avatar) {
		// TO DO
	}

	var newPassword string
	if !govalidator.IsNull(user.OldPassword) && !govalidator.IsNull(user.NewPassword) {
		if bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(user.OldPassword)) == nil {
			newPassword, err = security.MakeShieldedPassword(user.NewPassword)
			if err != nil {
				return entity.User{}, err
			}
			_, err = h.DB.Exec(
				"UPDATE user SET password = $1 WHERE id = $2",
				newPassword,
				id,
			)
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
	_, err := h.DB.Exec(
		"DELETE FROM user WHERE id = $1",
		id,
	)

	return err
}

func (h *UserDBManager) GetUserByLogin(email string, password string) (entity.User, error) {
	user := entity.User{Email: email}

	row := h.DB.QueryRow("SELECT id, password FROM user WHERE email = $1", email)

	err := row.Scan(user.ID, user.Password)
	if err != nil {
		return entity.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return entity.User{}, errors.New("user does not exist in DB")
	}

	// ADD AVATAR

	return user, nil
}

func (h *UserDBManager) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	result, err := h.DB.Query(
		"SELECT base_title, currency_title FROM watchlist WHERE user_id = $1",
		id,
	)
	defer result.Close()
	if err != nil {
		return nil, err
	}

	currencies := make([]entity.Currency, 0)
	for result.Next() {
		var currency entity.Currency
		if err := result.Scan(currency.Base, currency.Title); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}
