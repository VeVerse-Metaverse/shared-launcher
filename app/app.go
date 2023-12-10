package app

import (
	"context"
	errors "dev.hackerman.me/artheon/l7-shared-launcher/errors"
	"dev.hackerman.me/artheon/l7-shared-launcher/logger"
	se "dev.hackerman.me/artheon/veverse-shared/executable"
	"fmt"
	"github.com/gonutz/w32/v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetAppExecutablePath(path string) (string, error) {
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		if !strings.HasSuffix(path, ".exe") {
			path = path + ".exe"
		}
	}

	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if the file exists
	fi, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if the file is a regular file
	if !fi.Mode().IsRegular() {
		return "", fmt.Errorf("app executable is not a regular file")
	}

	// Open the file
	f, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	// Check if the file is executable
	var isExe bool
	isExe, err = se.IsExecutable(f)
	if err != nil {
		return "", fmt.Errorf("failed to check if file is executable: %w", err)
	}

	if !isExe {
		return "", fmt.Errorf("app executable is not executable")
	}

	return absPath, nil
}

func getExecutableDirectory() (string, error) {
	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	return executableDir, nil
}

// GetAppExecutable returns the executable path for the given app located in the given directory
func GetAppExecutable(ctx context.Context, dir string, app string) (string, error) {
	l, _ := ctx.Value("logger").(logger.Logger)

	var err error

	executableDir, err := getExecutableDirectory()
	if err != nil {
		if l != nil {
			l.Error("failed to get executable directory")
		}
		return "", errors.WrappedError{Message: "failed to get executable directory", Err: err}
	}

	appsDir := filepath.Join(executableDir, dir)
	_, err = os.Stat(appsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("apps directory does not exist: %s", appsDir)
		}
		return "", fmt.Errorf("failed to stat apps directory: %w", err)
	}

	appDir := filepath.Join(appsDir, app)
	_, err = os.Stat(appDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("app directory does not exist: %s", appDir)
		}
		return "", fmt.Errorf("failed to stat app directory: %w", err)
	}

	// Try to find the executable using the given name
	var exePath string
	if app != "" {
		exePath, err = GetAppExecutablePath(appDir)
		if err == nil {
			return exePath, nil
		} else {
			log.Printf("failed to get app executable by name: %s", err)
		}
	}

	// Walk the directory to find the executable
	err = filepath.WalkDir(appDir, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		fi, err := os.Stat(absPath)
		if err != nil {
			println(fmt.Sprintf("failed to get file info: %s", err))
			return nil
		}

		if !fi.Mode().IsRegular() {
			// Skip non-regular files (directories, symlinks, etc.)
			println(fmt.Sprintf("skipping non-regular file: %s", absPath))
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		//goland:noinspection GoUnhandledErrorResult
		defer f.Close()

		var isExe bool
		isExe, err = se.IsExecutable(f)
		if err != nil {
			return fmt.Errorf("failed to check if file is executable: %w", err)
		}

		if !isExe {
			println(fmt.Sprintf("skipping non-executable file: %s", absPath))
			return nil
		}

		// Check executable product version and publisher
		//goland:noinspection GoBoolExpressions
		if runtime.GOOS == "windows" {
			err = checkExecutableVersion(path)
			if err != nil {
				println(fmt.Sprintf("failed to check executable version: %s", err))
				return nil
			}
		}

		// Found the executable to use, we expect it to be the correct one
		exePath = path

		// Returning io.EOF to stop the walk, indicating that the walk was successful
		return io.EOF
	})

	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to find app executable: %w", err)
	}

	if exePath == "" {
		return "", fmt.Errorf("failed to find app executable")
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return "", fmt.Errorf("failed to find app executable: %w", err)
	}

	return exePath, nil
}

var UnknownProductVersionError error = errors.New("unknown product version")

func checkExecutableVersion(path string) error {
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		blockSize := w32.GetFileVersionInfoSize(path)
		if blockSize == 0 {
			return fmt.Errorf("failed to get file version info size")
		}

		block := make([]byte, blockSize)

		if !w32.GetFileVersionInfo(path, block) {
			return fmt.Errorf("failed to get file version info")
		}

		if productVersion, ok := w32.VerQueryValueString(block, "040904b0", w32.ProductVersion); ok {
			if !(strings.Contains(productVersion, "++UE") && strings.Contains(productVersion, "+Release")) {
				return UnknownProductVersionError
			}
		} else {
			return UnknownProductVersionError
		}

		if companyName, ok := w32.VerQueryValueString(block, "040904b0", w32.CompanyName); ok {
			if !(strings.Contains(companyName, "Epic Games") || strings.Contains(companyName, "LE7EL")) {
				return UnknownProductVersionError
			}
		} else {
			return UnknownProductVersionError
		}
	}

	return nil
}
