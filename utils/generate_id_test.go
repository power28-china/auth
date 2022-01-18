package utils

import (
	"testing"

	"github.com/power28-china/auth/utils/logger"
)

// BenchmarkGetID test benchmark for GetID.
func BenchmarkGetID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetID()
	}
}

func TestGetID(t *testing.T) {
	for i := 0; i < 100; i++ {
		logger.Sugar.Debugf("ID: %d", GetID())
	}
}
