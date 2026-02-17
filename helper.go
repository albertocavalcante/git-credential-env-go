package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	actionGet   = "get"
	actionStore = "store"
	actionErase = "erase"

	argCount = 2

	keyProtocol = "protocol"
	keyHost     = "host"
	keyPath     = "path"

	defaultProtocol = "https"
	defaultHost     = "github.com"
	defaultUsername = "x-access-token"
	// #nosec G101 -- This is an environment variable name, not a credential value.
	defaultTokenEnv = "GIT_CREDENTIAL_TOKEN"

	// #nosec G101 -- This is an environment variable name, not a credential value.
	tokenEnvOverrideVar = "GIT_CREDENTIAL_TOKEN_ENV"
	hostEnvVar          = "GIT_CREDENTIAL_HOST"
	protocolEnvVar      = "GIT_CREDENTIAL_PROTOCOL"
	usernameEnvVar      = "GIT_CREDENTIAL_USERNAME"
	pathPrefixEnvVar    = "GIT_CREDENTIAL_PATH_PREFIX"
)

type credentialRequest struct {
	protocol string
	host     string
	path     string
}

type helperConfig struct {
	protocol   string
	host       string
	username   string
	token      string
	pathPrefix string
}

func run(args []string, in io.Reader, out io.Writer) error {
	if len(args) != argCount {
		return fmt.Errorf("usage: %s <%s|%s|%s>", filepath.Base(args[0]), actionGet, actionStore, actionErase)
	}

	switch args[1] {
	case actionGet:
		req, err := parseRequest(in)
		if err != nil {
			return err
		}

		cfg := loadConfig()
		if !cfg.matches(req) {
			return nil
		}

		return writeResponse(out, cfg.username, cfg.token)
	case actionStore, actionErase:
		_, err := parseRequest(in)
		return err
	default:
		return fmt.Errorf("unsupported action %q", args[1])
	}
}

func parseRequest(in io.Reader) (credentialRequest, error) {
	req := credentialRequest{}
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		switch key {
		case keyProtocol:
			req.protocol = value
		case keyHost:
			req.host = value
		case keyPath:
			req.path = value
		}
	}

	if err := scanner.Err(); err != nil {
		return credentialRequest{}, fmt.Errorf("read request: %w", err)
	}

	return req, nil
}

func loadConfig() helperConfig {
	tokenEnv := getenvOrDefault(tokenEnvOverrideVar, defaultTokenEnv)

	return helperConfig{
		protocol:   getenvOrDefault(protocolEnvVar, defaultProtocol),
		host:       getenvOrDefault(hostEnvVar, defaultHost),
		username:   getenvOrDefault(usernameEnvVar, defaultUsername),
		token:      os.Getenv(tokenEnv),
		pathPrefix: os.Getenv(pathPrefixEnvVar),
	}
}

func getenvOrDefault(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}

func (cfg helperConfig) matches(req credentialRequest) bool {
	if cfg.token == "" || req.host == "" {
		return false
	}

	if req.protocol != "" && !strings.EqualFold(req.protocol, cfg.protocol) {
		return false
	}

	if !hostMatches(req.host, cfg.host) {
		return false
	}

	if cfg.pathPrefix != "" && !strings.HasPrefix(req.path, cfg.pathPrefix) {
		return false
	}

	return true
}

func hostMatches(rawHost string, expectedHost string) bool {
	hostOnly := rawHost
	if host, _, found := strings.Cut(rawHost, ":"); found {
		hostOnly = host
	}

	return strings.EqualFold(hostOnly, expectedHost)
}

func writeResponse(out io.Writer, username string, token string) error {
	writer := bufio.NewWriter(out)
	if _, err := fmt.Fprintf(writer, "username=%s\npassword=%s\n\n", username, token); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush response: %w", err)
	}

	return nil
}
