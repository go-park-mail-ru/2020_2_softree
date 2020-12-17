package entity

type Description struct {
	Status    int
	Function  string
	Action    string
	Value     interface{}
	ErrorJSON ErrorJSON
}
