package Typ

type REC struct {
	Key           string
	Value         []byte
	Section       bool
	SectionValues []REC
}
