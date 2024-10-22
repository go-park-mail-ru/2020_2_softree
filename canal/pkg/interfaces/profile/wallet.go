package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (p *Profile) GetWallets(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetWallets")

	desc, wallets, err := p.paymentLogic.ReceiveWallets(r.Context(), r.Context().Value(entity.UserIdKey).(int64))
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(wallets)
	if err != nil {
		desc = entity.Description{Function: "GetWallets", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err = w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetWallets", Action: "Write"}, err)
	}
}

func (p *Profile) SetWallet(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "SetWallet")

	wallet, desc, err := entity.GetWalletFromBody(r.Body)
	if err != nil {
		desc.Function = "SetWallet"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}
	wallet.UserId = r.Context().Value(entity.UserIdKey).(int64)

	if desc, err = p.paymentLogic.SetWallet(r.Context(), wallet); err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	w.WriteHeader(http.StatusCreated)
	metric.RecordHitMetric(http.StatusCreated, r.URL.Path)
}
