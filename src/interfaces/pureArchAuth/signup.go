package pureArchAuth

import (
	"context"
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
	"server/src/infrastructure/log"
)

func (a *Authenticate) PrepareUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user entity.User
		var err error
		if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.GlobalLogger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errors := user.Validate("signup")
		if errors.NotEmpty {
			createInternalServerError(&errors, w)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authenticate) SaveUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(entity.User)
		var err error
		user, err = a.userApp.SaveUser(user)
		if err != nil {
			log.GlobalLogger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authenticate) Signup(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entity.User)

	cookie, err := a.cookie.CreateCookie()
	if err != nil {
		log.GlobalLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	if err := a.auth.CreateAuth(user.ID, cookie.Value); err != nil {
		log.GlobalLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func createInternalServerError(errors *jsonRealisation.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		log.GlobalLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
