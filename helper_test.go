package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunGetMatchReturnsCredentials(t *testing.T) {
	t.Setenv(defaultTokenEnv, "token-123")

	input := strings.NewReader("protocol=https\nhost=github.com\npath=albertocavalcante/groovy-parser-go\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	expected := "username=x-access-token\npassword=token-123\n\n"
	if output.String() != expected {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunGetWithHostPortMatch(t *testing.T) {
	t.Setenv(defaultTokenEnv, "token-123")

	input := strings.NewReader("protocol=https\nhost=github.com:443\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	expected := "username=x-access-token\npassword=token-123\n\n"
	if output.String() != expected {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunGetNoTokenReturnsNoCredentials(t *testing.T) {
	input := strings.NewReader("protocol=https\nhost=github.com\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if output.Len() != 0 {
		t.Fatalf("expected empty output, got %q", output.String())
	}
}

func TestRunGetHostMismatchReturnsNoCredentials(t *testing.T) {
	t.Setenv(defaultTokenEnv, "token-123")

	input := strings.NewReader("protocol=https\nhost=gitlab.com\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if output.Len() != 0 {
		t.Fatalf("expected empty output, got %q", output.String())
	}
}

func TestRunGetPathPrefixFilter(t *testing.T) {
	t.Setenv(defaultTokenEnv, "token-123")
	t.Setenv(pathPrefixEnvVar, "albertocavalcante/")

	input := strings.NewReader("protocol=https\nhost=github.com\npath=another/repo\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if output.Len() != 0 {
		t.Fatalf("expected empty output, got %q", output.String())
	}
}

func TestRunGetCustomTokenEnv(t *testing.T) {
	t.Setenv(tokenEnvOverrideVar, "MY_APP_TOKEN")
	t.Setenv("MY_APP_TOKEN", "token-abc")

	input := strings.NewReader("protocol=https\nhost=github.com\n\n")
	var output bytes.Buffer

	err := run([]string{"git-credential-ghenv", actionGet}, input, &output)
	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	expected := "username=x-access-token\npassword=token-abc\n\n"
	if output.String() != expected {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunStoreAndEraseNoop(t *testing.T) {
	input := strings.NewReader("protocol=https\nhost=github.com\n\n")

	err := run([]string{"git-credential-ghenv", actionStore}, input, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("store returned error: %v", err)
	}

	err = run([]string{"git-credential-ghenv", actionErase}, strings.NewReader("host=github.com\n\n"), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("erase returned error: %v", err)
	}
}

func TestRunUnsupportedAction(t *testing.T) {
	err := run([]string{"git-credential-ghenv", "invalid"}, strings.NewReader(""), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestRunInvalidArgumentCount(t *testing.T) {
	err := run([]string{"git-credential-ghenv"}, strings.NewReader(""), bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected usage error")
	}
}
