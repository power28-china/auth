package domain

import (
	"testing"

	"github.com/power28-china/auth/utils/logger"
)

func TestAppAuth(t *testing.T) {
	auth := &AuthApp{}
	err := auth.Auth()
	if err != nil {
		t.Errorf("TestAppAuth failed: %v", err)
	}
	logger.Sugar.Debugf("AppAuth: %#v", auth)
}
