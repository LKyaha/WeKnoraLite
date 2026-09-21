//go:build windows && !bindings

package main

import (
    "os"
    "path/filepath"
    "strings"
)

// Anchor installed resources to the executable directory and keep persistent
// state outside Program Files. Diagnostics must begin before any of this work.
func init() {
    initWindowsStartupLog()
    if exe, err := os.Executable(); err == nil {
        windowsStartupStage("executable=" + exe)
        if err := os.Chdir(filepath.Dir(exe)); err != nil {
            windowsStartupStage("ERROR change working directory: " + err.Error())
        }
    } else {
        windowsStartupStage("ERROR locate executable: " + err.Error())
    }
    if cwd, err := os.Getwd(); err == nil {
        windowsStartupStage("working directory=" + cwd)
    }
    for _, resource := range []string{".env", "config/config.yaml", "web/index.html", "migrations/sqlite"} {
        if _, err := os.Stat(resource); err != nil {
            windowsStartupStage("ERROR resource " + resource + ": " + err.Error())
        } else {
            windowsStartupStage("resource OK: " + resource)
        }
    }

    configDir, err := os.UserConfigDir()
    if err != nil {
        windowsStartupStage("ERROR locate user config directory: " + err.Error())
        return
    }
    appDir := filepath.Join(configDir, "WeKnora Lite")
    dataDir := filepath.Join(appDir, "data")
    filesDir := filepath.Join(dataDir, "files")
    if err := os.MkdirAll(filesDir, 0755); err != nil {
        windowsStartupStage("ERROR create data directory: " + err.Error())
        return
    }
    windowsDesktopDefault("DB_PATH", filepath.Join(dataDir, "weknora.db"))
    windowsDesktopDefault("LOCAL_STORAGE_BASE_DIR", filesDir)
    if windowsStartupFile != nil && strings.TrimSpace(os.Getenv("LOG_PATH")) == "" {
        windowsDesktopDefault("LOG_PATH", windowsStartupFile.Name())
    } else {
        windowsDesktopDefault("LOG_PATH", filepath.Join(appDir, "weknora.log"))
    }
    windowsStartupStage("Windows bootstrap complete; DB=" + os.Getenv("DB_PATH"))
}

func windowsDesktopDefault(key, value string) {
    if strings.TrimSpace(os.Getenv(key)) == "" {
        _ = os.Setenv(key, value)
    }
}
