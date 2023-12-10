package download

import (
	"context"
	"dev.hackerman.me/artheon/l7-shared-launcher/errors"
	"dev.hackerman.me/artheon/l7-shared-launcher/logger"
	se "dev.hackerman.me/artheon/veverse-shared/executable"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// ErrInvalidPath Returned when the path is invalid
	ErrInvalidPath = errors.WrappedError{Message: "invalid path"}
	// ErrFileExists Returned when the file exists and the size is the same
	ErrFileExists = errors.WrappedError{Message: "file exists"}
	// ErrHttpError Returned when the HTTP request returns an error
	ErrHttpError = errors.WrappedError{Message: "http error"}
	// ErrFailedToCreateFile Returned when the file creation failed
	ErrFailedToCreateFile = errors.WrappedError{Message: "failed to create file"}
	// ErrFailedToCreateDir Returned when the directory creation failed
	ErrFailedToCreateDir = errors.WrappedError{Message: "failed to create directory"}
	// ErrFailedToWriteFile Returned when the file writing failed
	ErrFailedToWriteFile = errors.WrappedError{Message: "failed to write file"}
)

// File downloads file to the filepath from url
func File(ctx context.Context, path string, url string, size int64) (err error) {
	// Get logger from context
	l, _ := ctx.Value("logger").(logger.Logger)

	if path == "" {
		if l != nil {
			l.Error("path is empty")
		}
		return ErrInvalidPath
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		if l != nil {
			l.Error(fmt.Sprintf("failed to get absolute path: %v", err))
		}
		return errors.WrappedError{Message: ErrInvalidPath.Error(), Err: err}
	}

	// Check if file exists
	stat, err := os.Stat(absPath)
	if err == nil {
		if size > 0 && stat.Size() == size {
			if l != nil {
				l.Info(fmt.Sprintf("skipping, file exists: %s, size matches: %d", absPath, size))
			}
			return errors.WrappedError{Message: ErrFileExists.Error(), Err: err}
		}
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		if l != nil {
			l.Error(fmt.Sprintf("failed to send a HTTP GET request: %v", err))
		}
		return errors.WrappedError{Message: ErrHttpError.Error(), Err: err}
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil && l != nil {
			l.Error(fmt.Sprintf("failed to close response body: %v", err))
		}
	}(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		if l != nil {
			l.Error(fmt.Sprintf("failed to download file %s to %s: bad status: %s", url, absPath, resp.Status))
		}
		return errors.WrappedError{Message: ErrHttpError.Error(), Err: err}
	}

	// Create the dir
	dir := filepath.Dir(absPath)
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		if l != nil {
			l.Error(fmt.Sprintf("failed to create a directory %s: %v", dir, err))
		}
		return errors.WrappedError{Message: ErrFailedToCreateDir.Error(), Err: err}
	}

	// Create the file
	out, err := os.Create(absPath)
	if err != nil {
		if l != nil {
			l.Error(fmt.Sprintf("failed to create a file downloaded %s to %s: %v", url, absPath, err))
		}
		return errors.WrappedError{Message: ErrFailedToCreateFile.Error(), Err: err}
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil && l != nil {
			l.Error(fmt.Sprintf("failed to close file: %v", err))
		}
	}(out)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		if l != nil {
			l.Error(fmt.Sprintf("failed to write a file downloaded %s to %s: %v", url, absPath, err))
		}
		return errors.WrappedError{Message: ErrFailedToWriteFile.Error(), Err: err}
	}

	// Change a file mode for known executables to make them executable

	// Open the file
	f, err := os.Open(absPath)
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

	if isExe {
		err = os.Chmod(absPath, 0755)
		if err != nil && l != nil {
			l.Warning(fmt.Sprintf("failed to change file mode for %s: %v", absPath, err))
		}
	}

	return nil
}
