package profile

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (p *Profile) createErrorJSON(err error) (errs entity.ErrorJSON) {
	if err .Error() == "wrong old password" {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Введен неверно старый пароль")
	}

	return
}

func (p *Profile) createServerError(errs *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errs)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
}

func (p *Profile) checkWalletFrom(w http.ResponseWriter, id uint64, transaction entity.PaymentHistory) bool {
	var exist bool
	var err error
	if exist, err = p.userApp.CheckWallet(id, transaction.From); err != nil {
		p.log.Info("func: checkWallets, with error while CheckWallet: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (p *Profile) checkWalletTo(w http.ResponseWriter, id uint64, transaction entity.PaymentHistory) bool {
	var exist bool
	var err error
	if exist, err = p.userApp.CheckWallet(id, transaction.To); err != nil {
		p.log.Info("func: checkWallets, with error while CheckWallet: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !exist {
		if err = p.userApp.CreateWallet(id, transaction.To); err != nil {
			p.log.Info("func: checkWallets, with error while CreateWallet: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
	}

	return true
}
