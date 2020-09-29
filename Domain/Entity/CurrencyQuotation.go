package Entity

type CurrencyQuotation struct {
	Id     uint64 `json:"id"`
	Value  int64  `json:"value"`
	Title  string `json:"title"`
	Change int64  `json:"change"`
}

type Quotations []CurrencyQuotation

func (q *Quotations) Add(qs ...CurrencyQuotation) {
	*q = append(*q, qs...)
}

func FindById(id uint64) CurrencyQuotation {
	return CurrencyQuotation{
		Id:     id,
		Value:  0,
		Title:  "",
		Change: 0,
	}
}
