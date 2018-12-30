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
		t.Fatalf("Error: %s\nExpected: %+v\nActual: %+v", errorMsg, expected, actual)
	}
}

func (testUtil *testUtil) HandleIfTestError(t *testing.T, err error, errorMsg string) {
	if err != nil {
		t.Fatalf("Error: %s\n%s", errorMsg, err.Error())
	}
}

func (testUtil *testUtil) AssertJSONEquals(t *testing.T, expected, actual, errorMsg string) {
	if !json.Valid([]byte(expected)) {
		testUtil.HandleIfTestError(t, fmt.Errorf("Invalid expected JSON: %s", expected), errorMsg)
	}
	if !json.Valid([]byte(actual)) {
		testUtil.HandleIfTestError(t, fmt.Errorf("Invalid actual JSON: %s", actual), errorMsg)
	}

	byteExpected := []byte(expected)
	byteActual := []byte(actual)

	compactExpected := bytes.NewBuffer([]byte{})
	compactActual := bytes.NewBuffer([]byte{})
	indentedExpected := bytes.NewBuffer([]byte{})
	indentedActual := bytes.NewBuffer([]byte{})

	// ignoring errors here, as already checked for validity.
	// will refactor, if found some other causes of error
	json.Compact(compactExpected, byteExpected)
	json.Compact(compactActual, byteActual)
	json.Indent(indentedExpected, byteExpected, "", "  ")
	json.Indent(indentedActual, byteActual, "", "  ")

	if compactExpected.String() != compactActual.String() {
		t.Fatalf("Error: JSON Assertion failed\nExpected JSON: \n%s\nActual JSON: \n%s\nExpected compact JSON: \n%s\nActual compact JSON: \n%s\n", indentedExpected, indentedActual, compactExpected, compactActual)
	}
}
