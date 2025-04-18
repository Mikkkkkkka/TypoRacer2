package model

type Play struct {
	UserId         int
	QuoteId        int
	WordsPerMinute float32
	Accuracy       float32
	Consistency    float32
}
