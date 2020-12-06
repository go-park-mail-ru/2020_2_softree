package entity

type Description struct {
	Status    int
	Err       error
	Function  string
	Action    string
	Value     interface{}
	ErrorJSON ErrorJSON
}
