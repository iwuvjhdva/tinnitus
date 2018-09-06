package tinnitus

import (
	"time"
)

type ClippedTransaction struct {
	Time     time.Time
	Controls map[string]float32
}

type Tinnitus struct {
	synthID        int32
	transactionIDs map[string]*ClippedTransaction
}
