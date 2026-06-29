package util

import (
	"testing"
	"time"
)

func TestRandomNumeric(t *testing.T) {
	t.Log(RandomNumeric(6))
}

func TestEndOfDay(t *testing.T) {
	t.Log(EndOfDay(time.Now()))
}
