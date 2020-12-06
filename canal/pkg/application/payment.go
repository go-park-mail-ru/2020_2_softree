package application

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	profile "server/profile/pkg/profile/gen"
)

type PaymentApp struct {
	profile   profile.ProfileServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewPaymentApp(profile profile.ProfileServiceClient, security repository.Utils) *PaymentApp {
	return &PaymentApp{profile: profile, security: security, sanitizer: *bluemonday.UGCPolicy()}
}

func (pmt *PaymentApp) ReceiveTransactions(ctx context.Context, id int64) (entity.Description, []entity.Payment) {
	history, err := pmt.profile.GetAllPaymentHistory(ctx, &profile.UserID{Id: id})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "ReceiveTransactions",
			Action:   "GetAllPaymentHistory",
			Err:      err,
		}, []entity.Payment{}
	}

	return entity.Description{}, entity.ConvertToPayment(history)
}
