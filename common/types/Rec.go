package types

type Rec struct {
	Key           string
	Value         []byte
	Section       bool
	SectionValues []Rec
}
