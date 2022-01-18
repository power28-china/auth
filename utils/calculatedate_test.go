package utils

import (
	"testing"
	"time"

	"github.com/power28-china/auth/utils/logger"
)

func TestGetFirstAndLastDateOfMonth(t *testing.T) {
	d := time.Now()
	logger.Sugar.Debugf("Now: %v", d)
	logger.Sugar.Debugf("First day: %v(%d)", GetFirstDateOfMonth(d), GetFirstDateOfMonth(d).UnixMilli())
	logger.Sugar.Debugf("Last day: %v(%d)", GetLastDateOfMonth(d), GetLastDateOfMonth(d).UnixMilli())
	logger.Sugar.Debugf("First day of this year: %v(%d)", GetFirstDateOfYear(d), GetFirstDateOfYear(d).UnixMilli())
	logger.Sugar.Debugf("Last day of this year: %v(%d)", GetLastDateOfYear(d), GetLastDateOfYear(d).UnixMilli())
}
