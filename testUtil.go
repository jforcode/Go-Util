package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type ITestUtil interface {
	AssertEquals(t *testing.T, expected, actual interface{}, errorMsg string)
	HandleIfTestError(t *testing.T, err error, errorMsg string)
	AssertJSONEquals(t *testing.T, expected, actual, errorMsg string)
}

type testUtil struct {
}

func (testUtil *testUtil) AssertEquals(t *testing.T, expected, actual interface{}, errorMsg string) {
	if !cmp.Equal(expected, actual) {
		t.Errorf("Error: %s\nExpected: %+v\nActual: %+v", errorMsg, expected, actual)
	}
}

func (testUtil *testUtil) HandleIfTestError(t *testing.T, err error, errorMsg string) {
	if err != nil {
		t.Errorf("Error: %s\n%s", errorMsg, err.Error())
	}
}

func (testUtil *testUtil) AssertJSONEquals(t *testing.T, expected, actual, errorMsg string) {
	if !json.Valid([]byte(expected)) {
		testUtil.HandleIfTestError(t, fmt.Errorf("Invalid expected JSON: %s", expected), errorMsg)
	}
	if !json.Valid([]byte(actual)) {
		testUtil.HandleIfTestError(t, fmt.Errorf("Invalid actual JSON: %s", actual), errorMsg)
	}

	compactExpected := bytes.NewBuffer([]byte{})
	err := json.Compact(compactExpected, []byte(expected))
	testUtil.HandleIfTestError(t, err, fmt.Sprintf("Couldn't compact expected JSON: %s", expected))

	compactActual := bytes.NewBuffer([]byte{})
	err = json.Compact(compactActual, []byte(actual))
	testUtil.HandleIfTestError(t, err, fmt.Sprintf("Couldn't compact actual JSON: %s", actual))

	testUtil.AssertEquals(t, compactExpected.String(), compactActual.String(), "JSON assertion failed")
}
