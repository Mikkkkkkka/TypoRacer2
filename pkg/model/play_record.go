package model

import (
	"time"

	"github.com/eiannone/keyboard"
)

type PlayRecord struct {
	UserId    uint
	QuoteId   uint
	KeyStream []KeyPress
}

type KeyPress struct {
	KeyEvent    keyboard.KeyEvent
	ElapsedTime time.Duration
}
