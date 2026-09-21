//go:build windows && !bindings

package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime/debug"
    "syscall"
    "time"
)

// Windows GUI executables have no visible console. Open a diagnostics file before
// main() so startup errors can be recovered without building another installer.
// Keep it independent from WeKnora's logger and .env/config loading.
var windowsStartupFile *os.File

func initWindowsStartupLog() {
    base := os.Getenv("LOCALAPPDATA")
    if base == "" {
        base, _ = os.UserCacheDir()
    }
    if base == "" {
        base = os.TempDir()
    }
    logDir := filepath.Join(base, "WeKnora Lite", "logs")
    if err := os.MkdirAll(logDir, 0700); err != nil {
        fmt.Fprintf(os.Stderr, "WeKnora Lite: cannot create startup log directory %s: %v\n", logDir, err)
        return
    }
    path := filepath.Join(logDir, "startup.log")
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
    if err != nil {
        fmt.Fprintf(os.Stderr, "WeKnora Lite: cannot open startup log %s: %v\n", path, err)
        return
    }
    windowsStartupFile = file
    os.Stdout = file
    os.Stderr = file
    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

    // Go's runtime prints unrecovered panics to its OS-level crash output,
    // not necessarily to the os.Stderr variable in a windowsgui binary.
    if err := debug.SetCrashOutput(file, debug.CrashOptions{}); err != nil {
        log.Printf("diagnostics: SetCrashOutput failed: %v", err)
    }
    kernel32 := syscall.NewLazyDLL("kernel32.dll")
    setStdHandle := kernel32.NewProc("SetStdHandle")
    for _, which := range []uintptr{0xfffffff5, 0xfffffff4} { // STD_OUTPUT_HANDLE, STD_ERROR_HANDLE
        result, _, callErr := setStdHandle.Call(which, file.Fd())
        if result == 0 {
            log.Printf("diagnostics: SetStdHandle(%x) failed: %v", which, callErr)
        }
    }
    log.Printf("========== WeKnora Lite launch %s ==========", time.Now().Format(time.RFC3339Nano))
    log.Printf("diagnostics: startup log=%s pid=%d", path, os.Getpid())
    windowsStartupStage("logging initialized")
}

func windowsStartupStage(stage string) {
    if windowsStartupFile == nil {
        return
    }
    log.Printf("startup: %s", stage)
    _ = windowsStartupFile.Sync()
}
