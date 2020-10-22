package profile

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/infrastructure/log"
)

func (p *Profile) UpdateUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := p.auth.CheckAuth(cookie.Value)
	if err != nil {
		log.GlobalLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user entity.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.GlobalLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, _ = p.userApp.UpdateUser(id, user)

	res, err := json.Marshal(user.MakePublicUser())
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
