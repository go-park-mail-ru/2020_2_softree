package entity

import profile "server/profile/pkg/profile/gen"

type Income struct {
	Id     int64 `json:"id"`
	Period string
}

func (in *Income) ConvertToGRPC() *profile.IncomeParameters {
	return &profile.IncomeParameters{Id: in.Id, Period: in.Period}
}
