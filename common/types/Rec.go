package types

//Rec Type that represents each Record
type Rec struct {
	Key           string
	Value         []byte
	Section       bool
	SectionValues []Rec
}
