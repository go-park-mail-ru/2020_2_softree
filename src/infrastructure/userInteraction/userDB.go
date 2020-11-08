package userInteraction

import (
	"database/sql"
	"errors"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
	_ "github.com/lib/pq"
)

type Handler struct {
	DB   *sql.DB
}

func (h *Handler) GetUserById(id uint64) (entity.User, error) {
	user := entity.User{}

	row := h.DB.QueryRow("SELECT id, email, password FROM user WHERE id = ?", id)

	err := row.Scan(user.ID, user.Email, user.Password)
	if err != nil {
		return entity.User{}, err
	}

	// ADD AVATAR

	return user, nil
}

func (h *Handler) SaveUser(user entity.User) (entity.User, error) {
	row := h.DB.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", user.Email)

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
		"INSERT INTO user (`email`, `password`) VALUES (?, ?)",
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

func (h *Handler) UpdateUser(id uint64, user entity.User) (entity.User, error) {
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
				"UPDATE user SET password = ? WHERE id = ?",
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

func (h *Handler) DeleteUser(id uint64) error {
	_, err := h.DB.Exec(
		"DELETE FROM user WHERE id = ?",
		id,
	)

	return err
}

func (h *Handler) GetUserByLogin(email string, password string) (entity.User, error) {
	user := entity.User{Email: email}

	row := h.DB.QueryRow("SELECT id, password FROM user WHERE email = ?", email)

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