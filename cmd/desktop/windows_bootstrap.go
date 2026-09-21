//go:build windows && !bindings

package main

import (
	"os"
	"path/filepath"
	"strings"
)

// The installed Windows shortcut may start with a working directory other than
// the installation directory. Config, SQLite migrations and the proxied web
// frontend are currently loaded from relative paths, so anchor them to the EXE.
// Keep persistent state in the user's config directory, never Program Files.
func init() {
	if exe, err := os.Executable(); err == nil {
		_ = os.Chdir(filepath.Dir(exe))
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	appDir := filepath.Join(configDir, "WeKnora Lite")
	dataDir := filepath.Join(appDir, "data")
	filesDir := filepath.Join(dataDir, "files")
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return
	}

	windowsDesktopDefault("DB_PATH", filepath.Join(dataDir, "weknora.db"))
	windowsDesktopDefault("LOCAL_STORAGE_BASE_DIR", filesDir)
	windowsDesktopDefault("LOG_PATH", filepath.Join(appDir, "weknora.log"))
}

func windowsDesktopDefault(key, value string) {
	if strings.TrimSpace(os.Getenv(key)) == "" {
		_ = os.Setenv(key, value)
	}
}
