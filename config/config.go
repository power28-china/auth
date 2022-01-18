package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config returns a Config object by specified key.
func Config(key string) string {
	if os.Getenv("APP_ENV") == "test" { // in test mode.
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env.test)")
			os.Exit(1)
		}

		// return the root of the project.
		basepath := filepath.Dir(filepath.Dir(file))
		// logger.Sugar.Debug("Project root: ", basepath)

		godotenv.Load(filepath.Join(basepath, ".env.common"), filepath.Join(basepath, ".env.development.local"))
		return os.Getenv(key)
	}

	if os.Getenv("APP_ENV") == "development" {
		godotenv.Load(".env.common", ".env.development.local")
	} else {
		godotenv.Load(".env.common", ".env.production.local")
	}
	return os.Getenv(key)
}
