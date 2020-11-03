package entity

type Rates struct {
	Values []Rate
}

type Rate struct {
	Name  string
	Value string
}
