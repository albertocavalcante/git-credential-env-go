# git-credential-ghenv

Minimal Git credential helper binary for private GitHub module fetches.

It implements the Git credential helper protocol and serves a token from
environment variables when the request matches configured host/protocol/path.

## Install

```bash
go install github.com/albertocavalcante/git-credential-ghenv@latest
```

## Configure

```bash
git config --global credential.helper ghenv
```

Set credentials and matching rules:

```bash
export GIT_CREDENTIAL_TOKEN="..."
export GIT_CREDENTIAL_HOST="github.com"                # default: github.com
export GIT_CREDENTIAL_PROTOCOL="https"                 # default: https
export GIT_CREDENTIAL_USERNAME="x-access-token"        # default: x-access-token
export GIT_CREDENTIAL_PATH_PREFIX="albertocavalcante/" # optional
```

Token env var can be redirected with:

```bash
export GIT_CREDENTIAL_TOKEN_ENV="EUKIA_TOKEN"
export EUKIA_TOKEN="..."
```

## GitHub Actions Example

```bash
go env -w GOPRIVATE=github.com/albertocavalcante/*
go env -w GONOSUMDB=github.com/albertocavalcante/*
git config --global credential.helper ghenv
```
