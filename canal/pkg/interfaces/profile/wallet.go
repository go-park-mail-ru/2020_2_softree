package profile

import (
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) GetWallets(w http.ResponseWriter, r *http.Request) {
	desc, wallets, err := p.paymentLogic.ReceiveWallets(r.Context(), r.Context().Value(entity.UserIdKey).(int64))
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := wallets.MarshalJSON()
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetWallets", Action: "Marshal", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetTransactions", Action: "Write"}, err)
	}
}

func (p *Profile) SetWallet(w http.ResponseWriter, r *http.Request) {
	wallet, desc, err := entity.GetWalletFromBody(r.Body)
	if err != nil {
		desc.Function = "SetWallet"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	wallet.UserId = r.Context().Value(entity.UserIdKey).(int64)

	if desc, err := p.paymentLogic.SetWallet(r.Context(), wallet); err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	p.recordHitMetric(http.StatusCreated)
}
