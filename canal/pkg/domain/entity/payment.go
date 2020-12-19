package entity

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	json "github.com/mailru/easyjson"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"net/http"
	profile "server/profile/pkg/profile/gen"
	"strconv"
)

//easyjson:json
type (
	Payment struct {
		Currency  string               `json:"currency"`
		Base      string               `json:"base"`
		Amount    decimal.Decimal      `json:"amount"`
		Value     decimal.Decimal      `json:"value"`
		Sell      bool                 `json:"sell"`
		UpdatedUp *timestamp.Timestamp `json:"updated_up"`
		UserId    int64
	}

	Payments []Payment
)

func GetTransactionFromBody(body io.ReadCloser) (Payment, Description, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return Payment{}, Description{Action: "ReadAll", Status: http.StatusInternalServerError}, err
	}
	defer body.Close()

	var pay Payment
	err = json.Unmarshal(data, &pay)
	if err != nil {
		return Payment{}, Description{Action: "Unmarshal", Status: http.StatusInternalServerError}, err
	}
	return pay, Description{}, nil
}

func ConvertToPayment(history *profile.AllHistory) Payments {
	payments := make(Payments, 0, len(history.History))
	for _, pay := range history.History {
		payments = append(payments, Payment{
			Currency:  pay.Currency,
			Base:      pay.Base,
			Amount:    decimal.NewFromFloat(pay.Amount),
			Value:     decimal.NewFromFloat(pay.Value),
			Sell:      !(pay.Sell == "false"),
			UpdatedUp: pay.UpdatedAt,
		})
	}

	return payments
}

func (pay *Payment) ConvertToGRPC() *profile.PaymentHistory {
	amount, _ := pay.Amount.Float64()
	value, _ := pay.Value.Float64()
	return &profile.PaymentHistory{
		Currency: pay.Currency,
		Base:     pay.Base,
		Amount:   amount,
		Value:    value,
		Sell:     strconv.FormatBool(pay.Sell),
	}
}
