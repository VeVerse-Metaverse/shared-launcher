package tests

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"dev.hackerman.me/artheon/l7-shared-launcher/app"
)

func TestGetAppExecutablePath(t *testing.T) {
	name := filepath.FromSlash("apps/example/example.exe")

	path, err := app.GetAppExecutablePath(name)
	if err != nil {
		t.Fatal(err)
	}

	if path == "" {
		t.Fatal("path is empty")
	}

	if !strings.HasSuffix(path, name) {
		t.Fatal("path does not end with name")
	}
}

func TestGetAppExecutablePathMissing(t *testing.T) {
	name := filepath.FromSlash("this/path/doesnt/exist.exe")

	_, err := app.GetAppExecutablePath(name)
	if err == nil {
		t.Fatal(err)
	}
}

func TestGetAppExecutableName(t *testing.T) {
	name := "example"

	//launcher := &testLauncher{}

	path, err := app.GetAppExecutable("../apps", name)
	if err != nil {
		t.Fatal(err)
	}

	if path == "" {
		t.Fatal("path is empty")
	}

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(path, name+".exe") {
			t.Fatal("path does not end with .exe")
		}
	} else {
		if !strings.HasSuffix(path, name) {
			t.Fatal("path does not end with name")
		}
	}

	println(path)
}

func TestGetAppExecutableGuid(t *testing.T) {
	name := "F72C7005-B0A3-4E89-9B13-2BF1CDAB73FA"

	//launcher := &testLauncher{}

	path, err := app.GetAppExecutable("../apps", name)
	if err != nil {
		t.Fatal(err)
	}

	if path == "" {
		t.Fatal("path is empty")
	}

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(path, ".exe") {
			println(path)
			t.Fatal("path does not end with .exe")
		}
	} else {
		if !strings.HasSuffix(path, name) {
			t.Fatal("path does not end with name")
		}
	}

	println(path)
}
