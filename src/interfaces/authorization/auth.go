package authorization

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/src/domain/entity"
)

func (a *Authenticate) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user entity.User
	if user, err = extractUserFromSession(cookie, a); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func extractUserFromSession(c *http.Cookie, a *Authenticate) (entity.User, error) {
	id, err := a.auth.CheckAuth(c.Value)
	if err != nil {
		return entity.User{}, errors.New("CheckAuth in extractUserFromSession")
	}

	user, err := a.userApp.GetUserById(id)
	if err != nil {
		return entity.User{}, errors.New("GetUserById in extractUserFromSession")
	}

	return user, nil
}
