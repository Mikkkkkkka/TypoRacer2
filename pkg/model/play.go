package model

type Play struct {
	UserId         uint
	QuoteId        uint
	WordsPerMinute float32
	Accuracy       float32
	Consistency    float32
}
