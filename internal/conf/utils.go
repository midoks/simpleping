package conf

import (
	"os"
	"os/user"
	"os/exec"
	"path/filepath"
	"sync"
	"runtime"
)

// IsWindowsRuntime returns true if the current runtime in Windows.
func IsWindowsRuntime() bool {
	return runtime.GOOS == "windows"
}


// CurrentUsername returns the username of the current user.
func CurrentUsername() string {
	username := os.Getenv("USER")
	if len(username) > 0 {
		return username
	}

	username = os.Getenv("USERNAME")
	if len(username) > 0 {
		return username
	}

	if user, err := user.Current(); err == nil {
		username = user.Username
	}
	return username
}

var (
	workDir     string
	workDirOnce sync.Once
)

// WorkDir returns the absolute path of work directory. It reads the value of environment
// variable IMAIL_WORK_DIR. When not set, it uses the directory where the application's
// binary is located.
func WorkDir() string {
	workDirOnce.Do(func() {
		workDir = os.Getenv("SP_WORK_DIR")
		if workDir != "" {
			return
		}

		workDir = filepath.Dir(AppPath())
	})

	return workDir
}

var (
	appPath     string
	appPathOnce sync.Once
)

// AppPath returns the absolute path of the application's binary.
func AppPath() string {
	appPathOnce.Do(func() {
		var err error
		appPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			panic("look executable path: " + err.Error())
		}

		appPath, err = filepath.Abs(appPath)
		if err != nil {
			panic("get absolute executable path: " + err.Error())
		}
	})

	return appPath
}

// ensureAbs prepends the WorkDir to the given path if it is not an absolute path.
func ensureAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(WorkDir(), path)
}

// CheckRunUser returns false if configured run user does not match actual user that
// runs the app. The first return value is the actual user name. This check is ignored
// under Windows since SSH remote login is not the main method to login on Windows.
func CheckRunUser(runUser string) (string, bool) {
	if IsWindowsRuntime() {
		return "", true
	}

	currentUser := CurrentUsername()
	return currentUser, runUser == currentUser
}
