package version

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestVersion(t *testing.T) {
	Version = "v1.0.0"
	Branch = "master"
	Hash = "1234567"
	BuildTime = "2020/01/07 11:18:00"
	Comment = "build from unit test"

	v, err := New("unit-test", "v0.0.1")
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	goVersion := runtime.Version()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	expectedStrTpl := `
unit-test (build from unit test)
	VERSION          : v1.0.0
	INTERNAL_VERSION : v0.0.1
	BRANCH           : master
	HASH             : 1234567
	BUILD_TIME       : 2020/01/07 11:18:00
	GO_VERSION       : %s
	WORK_DIR         : %s
	RUN_TIME         : %s
`
	expectedStr := fmt.Sprintf(expectedStrTpl, goVersion, wd, v.RunTime)
	actualStr, err := v.String()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if expectedStr != actualStr {
		t.Fatalf("want %s, but got %s", expectedStr, actualStr)
	}

	_, err = v.JSON()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
}
