# git-credential-ghenv

Minimal Git credential helper binary for HTTPS Git operations.

`git-credential-ghenv` implements the Git credential helper protocol and returns
credentials from environment variables when host/protocol/path rules match.

It is intentionally simple:

- no token refresh
- no secure storage backend
- no provider-specific API calls
- no side effects for `store`/`erase` (stateless)

## Install

```bash
go install github.com/albertocavalcante/git-credential-ghenv@latest
```

## Configure

```bash
git config --global credential.helper ghenv
```

Set env vars:

```bash
export GIT_CREDENTIAL_TOKEN="..."
export GIT_CREDENTIAL_HOST="github.com"                # default: github.com
export GIT_CREDENTIAL_PROTOCOL="https"                 # default: https
export GIT_CREDENTIAL_USERNAME="your-provider-username"
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"            # optional
```

`GIT_CREDENTIAL_HOST` must be hostname only (no scheme).
`GIT_CREDENTIAL_USERNAME` is provider/token-type dependent.
The helper default username is `x-access-token`, but this is not universal.

## Quick Start

```bash
export GIT_CREDENTIAL_TOKEN="..."
git config --global credential.helper ghenv
git config --global credential.useHttpPath true
```

Then set at least:

```bash
export GIT_CREDENTIAL_HOST="your-git-host.example.com"
export GIT_CREDENTIAL_USERNAME="your-expected-username"
```

Optionally scope credentials to a subset of repos:

```bash
export GIT_CREDENTIAL_PATH_PREFIX="team-a/"
```

## Env Vars Reference

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `GIT_CREDENTIAL_TOKEN` | yes (unless override) | none | Password/token value returned to Git |
| `GIT_CREDENTIAL_TOKEN_ENV` | no | `GIT_CREDENTIAL_TOKEN` | Name of env var that holds token |
| `GIT_CREDENTIAL_HOST` | no | `github.com` | Hostname to match |
| `GIT_CREDENTIAL_PROTOCOL` | no | `https` | Protocol to match when request includes protocol |
| `GIT_CREDENTIAL_USERNAME` | no | `x-access-token` | Username returned on match |
| `GIT_CREDENTIAL_PATH_PREFIX` | no | empty | Prefix match against Git credential `path` |

## Matching Rules and Caveats

1. Matching conditions for `get`:
- token must be non-empty
- request host must be present and match `GIT_CREDENTIAL_HOST`
- protocol must match if request includes protocol
- path must start with `GIT_CREDENTIAL_PATH_PREFIX` when prefix is set

2. Host matching:
- host comparison is case-insensitive
- incoming host port is ignored (`github.com:443` matches `github.com`)

3. Path prefix matching:
- prefix is case-sensitive
- no normalization is performed
- `.git` suffix is part of the path, so prefix should account for your remote URL style

4. `credential.useHttpPath` is critical for path scoping:
- default Git behavior may omit path when querying helpers
- if path is omitted, prefix filtering cannot work
- set `git config --global credential.useHttpPath true` when using `GIT_CREDENTIAL_PATH_PREFIX`

5. `store` and `erase` are no-ops:
- this helper does not persist credentials
- if you also use other helpers, they may still cache values

## Platform Playbooks

Use these as starting points, then validate against your provider docs and token
type.

### GitHub.com and GitHub Enterprise

Common HTTPS path:

- `owner/repo.git`

Recommended defaults:

```bash
export GIT_CREDENTIAL_HOST="github.com"         # or ghe.example.com
export GIT_CREDENTIAL_USERNAME="x-access-token" # conventional, not strictly universal
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"
```

Gotchas:

- GitHub App installation token for Git over HTTPS is documented as
  `x-access-token:TOKEN`.
- Fine-grained PATs can fail if repo access is not granted.
- `GITHUB_TOKEN` in Actions is repository-scoped and often cannot read other
  private repos unless permissions/installation boundaries allow it.

### GitLab.com and Self-Managed GitLab

Common HTTPS path:

- `group/project.git`
- `group/subgroup/project.git`

Suggested config:

```bash
export GIT_CREDENTIAL_HOST="gitlab.com"  # or gitlab.example.com
export GIT_CREDENTIAL_USERNAME="oauth2"  # safe default for OAuth token flows
export GIT_CREDENTIAL_PATH_PREFIX="my-group/"
```

Gotchas:

- GitLab documents that PAT auth over HTTPS ignores username value, but requires
  it to be non-empty.
- If your GitLab instance is hosted under a relative URL (for example
  `example.com/gitlab`), include that leading segment in path prefix:
  `GIT_CREDENTIAL_PATH_PREFIX="gitlab/group/"`
- Token scopes must include `read_repository` or `write_repository` as needed.

### Gitea

Common HTTPS path:

- `owner/repo.git`

Suggested config:

```bash
export GIT_CREDENTIAL_HOST="gitea.example.com"
export GIT_CREDENTIAL_USERNAME="your-gitea-username"
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"
```

Gotchas:

- With MFA/2FA enabled, password auth for Git over HTTP is not viable; use an
  access token as password.
- Some deployments run under a subpath; include it in prefix when needed.

### Bitbucket Cloud

Common HTTPS path:

- `workspace/repo.git`

App password:

```bash
export GIT_CREDENTIAL_HOST="bitbucket.org"
export GIT_CREDENTIAL_USERNAME="your-bitbucket-username"
export GIT_CREDENTIAL_PATH_PREFIX="my-workspace/"
```

API token:

```bash
export GIT_CREDENTIAL_HOST="bitbucket.org"
export GIT_CREDENTIAL_USERNAME="x-bitbucket-api-token-auth"
export GIT_CREDENTIAL_PATH_PREFIX="my-workspace/"
```

Workspace/project/repository access token:

```bash
export GIT_CREDENTIAL_HOST="bitbucket.org"
export GIT_CREDENTIAL_USERNAME="x-token-auth"
export GIT_CREDENTIAL_PATH_PREFIX="my-workspace/"
```

Gotchas:

- Bitbucket Cloud token systems differ by token type, and required username may
  differ.
- If you see username/password errors, re-check username convention for your
  token type and clear stale local credentials.

### Bitbucket Server / Data Center

Common HTTPS path is often:

- `scm/PROJECT/repo.git`

Example:

```bash
export GIT_CREDENTIAL_HOST="bitbucket.example.com"
export GIT_CREDENTIAL_USERNAME="your-username"
export GIT_CREDENTIAL_PATH_PREFIX="scm/PROJECT/"
```

Important caveat:

- `scm/` is common, not universal. Reverse proxies and context paths can change
  URL shape.
- Always derive prefix from your actual remote URL path.

If your remote is:

- `https://bitbucket.example.com/scm/PROJ/repo.git`
  use `GIT_CREDENTIAL_PATH_PREFIX="scm/PROJ/"`
- `https://bitbucket.example.com/bitbucket/scm/PROJ/repo.git`
  use `GIT_CREDENTIAL_PATH_PREFIX="bitbucket/scm/PROJ/"`

## CI/CD Caveats and Workarounds

### General

- Prefer short-lived tokens where possible.
- Inject token through CI secret manager; never commit it.
- Avoid embedding token in remote URL.
- If multiple helpers are configured, order can cause unexpected credential
  source selection.

### GitHub Actions

- `GITHUB_TOKEN` may be insufficient for cross-repo private module access.
- For cross-repo/private org boundaries, use GitHub App installation token or
  scoped PAT.

### GitLab CI / Other CI

- Ensure non-interactive git auth. Go and Git frequently run with prompts
  disabled in automation.
- Rotate tokens and monitor expiry.

## Go Private Modules Caveats

For private modules, this helper solves Git HTTPS credentials only. You still
need Go module privacy settings.

Typical direct-VCS setup:

```bash
go env -w GOPRIVATE=your.git.host/your-org/*
go env -w GONOSUMDB=your.git.host/your-org/*
git config --global credential.helper ghenv
```

Notes:

- `GOPRIVATE` marks modules as private and defaults `GONOPROXY`/`GONOSUMDB`
  behavior.
- You generally cannot use `sum.golang.org` for private modules matched by
  `GOPRIVATE`/`GONOSUMDB`.
- If your workflow uses a private proxy, configure `GOPROXY` accordingly.

## FAQ

### Do I need `credential.useHttpPath=true`?

Yes, if you rely on `GIT_CREDENTIAL_PATH_PREFIX`. Without it, Git may omit path
and your prefix rule will not apply.

### Is `x-access-token` always correct?

No. It is only this helper's default. Use the username required by your
provider/token type.

### Can I target multiple hosts with one process?

Not with one static env set. Use separate jobs/processes or wrap the helper to
switch env vars.

### Can this helper refresh expired OAuth tokens?

No. This helper returns a static token from env. Token refresh must happen
outside this helper.

### Why does authentication still fail after updating env vars?

Another credential helper may be returning cached/stale credentials first.
Inspect `git config --show-origin --get-all credential.helper` and clear stale
entries in system keychain/helper.

## Troubleshooting

### Symptom: helper returns nothing

Checks:

1. `GIT_CREDENTIAL_TOKEN` (or override var) is set and non-empty.
2. Request host matches `GIT_CREDENTIAL_HOST`.
3. Protocol is `https` or matches your configured protocol.
4. Path prefix, if set, actually matches request path.

### Symptom: wrong username for provider

Set explicitly:

```bash
export GIT_CREDENTIAL_USERNAME="provider-required-username"
```

### Symptom: prefix filter seems ignored

Enable:

```bash
git config --global credential.useHttpPath true
```

Then confirm your remote URL path and prefix align.

### Symptom: still using old credentials

Likely cached elsewhere. Clear credentials from:

- `osxkeychain` / Windows Credential Manager / libsecret helper
- previously configured custom helpers

### Symptom: works locally, fails in CI

- token missing in CI env
- different host value in CI
- token expired or insufficient scopes
- helper not configured in CI step where Git runs

## Security Recommendations

- Prefer least-privilege token scopes.
- Prefer short-lived tokens where supported.
- Do not print token env values in logs.
- Avoid storing tokens in remote URLs or shell history.
- Rotate tokens periodically.

## References (Official Docs)

- Git helper protocol and matching:
  [gitcredentials](https://git-scm.com/docs/gitcredentials),
  [git-credential](https://git-scm.com/docs/git-credential)
- GitHub:
  [Managing personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens),
  [Authenticating as a GitHub App installation](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation),
  [Use GITHUB_TOKEN for authentication](https://docs.github.com/en/actions/how-tos/security-for-github-actions/security-guides/automatic-token-authentication),
  [About GITHUB_TOKEN](https://docs.github.com/en/actions/concepts/security/github_token)
- GitLab:
  [Personal access tokens](https://docs.gitlab.com/user/profile/personal_access_tokens/),
  [Clone with token](https://docs.gitlab.com/topics/git/clone/),
  [OAuth token for Git over HTTPS](https://docs.gitlab.com/api/oauth2/)
- Gitea:
  [MFA and Git over HTTP token usage](https://docs.gitea.com/usage/user-setting/multi-factor-authentication),
  [Repository mirror auth examples](https://docs.gitea.com/usage/repository/repo-mirror)
- Bitbucket Cloud:
  [Using App passwords](https://support.atlassian.com/bitbucket-cloud/docs/using-app-passwords/),
  [Using API tokens](https://support.atlassian.com/bitbucket-cloud/docs/using-api-tokens/),
  [Repository access tokens](https://support.atlassian.com/bitbucket-cloud/docs/using-access-tokens/),
  [Workspace access tokens](https://support.atlassian.com/bitbucket-cloud/docs/workspace-access-tokens/),
  [Project access tokens](https://support.atlassian.com/bitbucket-cloud/docs/project-access-tokens/)
- Bitbucket Data Center / Server:
  [Personal access tokens](https://confluence.atlassian.com/bitbucketserver/personal-access-tokens-939515499.html),
  [HTTP access tokens](https://confluence.atlassian.com/bitbucketserver/http-access-tokens-939515499.html)
- Go modules/private modules:
  [Go Modules Reference](https://go.dev/ref/mod)
