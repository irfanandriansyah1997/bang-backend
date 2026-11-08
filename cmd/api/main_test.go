package main

import (
	"bytes"
	"testing"
)

func TestStartupMessage(t *testing.T) {
	result := startupMessage("Benah API")
	expected := "Benah API starting"

	if result != expected {
		t.Errorf("startupMessage() = %q, want %q", result, expected)
	}
}

func TestStartupMessage_Table(t *testing.T) {
	tests := []struct {
		name           string
		appName        string
		expectedResult string
	}{
		{name: "api", appName: "Benah API", expectedResult: "Benah API starting"},
		{name: "worker", appName: "Benah Worker", expectedResult: "Benah Worker starting"},
		{name: "empty string", appName: "", expectedResult: ""},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := startupMessage(testCase.appName)
			if result != testCase.expectedResult {
				t.Errorf("startupMessage() = %q, want %q", result, testCase.expectedResult)
			}
		})
	}
}

func TestRun(t *testing.T) {
	var buf bytes.Buffer

	if err := run(&buf); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	expectedMessagge := "Benah API starting\n"
	if got := buf.String(); got != expectedMessagge {
		t.Errorf("run() wrote %q, want %q", got, expectedMessagge)
	}
}
