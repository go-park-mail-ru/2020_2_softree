package entity

import (
	"github.com/golang/protobuf/ptypes"
	json "github.com/mailru/easyjson"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"net/http"
	profile "server/profile/pkg/profile/gen"
	"strconv"
	"time"
)

type Payment struct {
	Currency  string          `json:"currency"`
	Base      string          `json:"base"`
	Amount    decimal.Decimal `json:"amount"`
	Value     decimal.Decimal `json:"value"`
	Sell      bool            `json:"sell"`
	UpdatedUp time.Time       `json:"updated_up"`
	UserId    int64
}

type Payments struct {
	Payments []Payment
}

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
	payments := make([]Payment, 0, len(history.History))
	for _, pay := range history.History {
		updated, _ := ptypes.Timestamp(pay.UpdatedAt)
		payments = append(payments, Payment{
			Currency:  pay.Currency,
			Base:      pay.Base,
			Amount:    decimal.NewFromFloat(pay.Amount),
			Value:     decimal.NewFromFloat(pay.Value),
			Sell:      !(pay.Sell == "false"),
			UpdatedUp: updated,
		})
	}

	return Payments{payments}
}

func (pay *Payment) ConvertToGRPC() *profile.PaymentHistory {
	amount, _ := pay.Amount.Float64()
	value, _ := pay.Value.Float64()
	updated, _ := ptypes.TimestampProto(pay.UpdatedUp)
	return &profile.PaymentHistory{
		Currency:  pay.Currency,
		Base:      pay.Base,
		Amount:    amount,
		Value:     value,
		Sell:      strconv.FormatBool(pay.Sell),
		UpdatedAt: updated,
	}
}
