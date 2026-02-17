# git-credential-ghenv

Minimal Git credential helper binary for HTTPS Git operations.

It implements the Git credential helper protocol and returns credentials from
environment variables when the incoming request matches host/protocol/path
rules.

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
export GIT_CREDENTIAL_USERNAME="your-provider-username"
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"            # optional
```

`GIT_CREDENTIAL_HOST` should be the hostname only (no scheme, no port).
`GIT_CREDENTIAL_USERNAME` depends on the Git provider and token type.
The helper default is `x-access-token`, but many providers require a different
username.

## `GIT_CREDENTIAL_PATH_PREFIX` Explained

`GIT_CREDENTIAL_PATH_PREFIX` is matched against Git's credential request
`path` field.

For GitHub, `path` is typically:

- `owner/repo`
- `owner/repo.git`

Examples:

- `GIT_CREDENTIAL_PATH_PREFIX="my-org/"`: match any repo under `my-org`
- `GIT_CREDENTIAL_PATH_PREFIX="my-org/private-repo"`: match one repo
- unset: no path restriction (host/protocol checks only)

## Provider Examples

### GitHub.com

```bash
export GIT_CREDENTIAL_HOST="github.com"
export GIT_CREDENTIAL_USERNAME="x-access-token"
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"
```

### GitLab.com or Self-Managed GitLab

```bash
export GIT_CREDENTIAL_HOST="gitlab.com"      # or gitlab.example.com
export GIT_CREDENTIAL_USERNAME="oauth2"
export GIT_CREDENTIAL_PATH_PREFIX="my-group/"
```

For GitLab PATs, username is generally not evaluated. Using `oauth2` keeps the
same config working for OAuth tokens as well.

### Gitea

```bash
export GIT_CREDENTIAL_HOST="gitea.example.com"
export GIT_CREDENTIAL_USERNAME="your-username"
export GIT_CREDENTIAL_PATH_PREFIX="my-org/"
```

### Bitbucket Cloud

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

Workspace, project, or repository access token:

```bash
export GIT_CREDENTIAL_HOST="bitbucket.org"
export GIT_CREDENTIAL_USERNAME="x-token-auth"
export GIT_CREDENTIAL_PATH_PREFIX="my-workspace/"
```

### Bitbucket Server / Data Center

```bash
export GIT_CREDENTIAL_HOST="bitbucket.example.com"
export GIT_CREDENTIAL_USERNAME="your-username"
export GIT_CREDENTIAL_PATH_PREFIX="scm/PROJ/"
```

Bitbucket Server/Data Center clone paths are commonly `scm/<project>/<repo>.git`.

## Custom Token Env Var

By default the helper reads `GIT_CREDENTIAL_TOKEN`. You can point it to any
other env var name:

```bash
export GIT_CREDENTIAL_TOKEN_ENV="MY_APP_TOKEN"
export MY_APP_TOKEN="..."
```

## References

- Git credentials storage and helper protocol:
  [gitcredentials](https://git-scm.com/docs/gitcredentials),
  [git-credential](https://git-scm.com/docs/git-credential)
- GitHub tokens:
  [Managing your personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
- GitLab tokens:
  [Personal access tokens](https://docs.gitlab.com/user/profile/personal_access_tokens/)
- Gitea auth usage:
  [Repository mirror auth examples](https://docs.gitea.com/usage/repository/repo-mirror)
- Bitbucket Cloud:
  [Using API tokens](https://support.atlassian.com/bitbucket-cloud/docs/using-api-tokens/),
  [Workspace access tokens](https://support.atlassian.com/bitbucket-cloud/docs/workspace-access-tokens/),
  [Project access tokens](https://support.atlassian.com/bitbucket-cloud/docs/project-access-tokens/),
  [Repository access tokens](https://support.atlassian.com/bitbucket-cloud/docs/repository-access-tokens/)
- Bitbucket Server/Data Center:
  [Personal access tokens](https://confluence.atlassian.com/bitbucketserver/personal-access-tokens-939515499.html),
  [HTTP access tokens](https://confluence.atlassian.com/bitbucketserver/http-access-tokens-939515499.html)

## CI Example

```bash
go env -w GOPRIVATE=github.com/my-org/*
go env -w GONOSUMDB=github.com/my-org/*
export GIT_CREDENTIAL_TOKEN="$MY_APP_TOKEN"
git config --global credential.helper ghenv
```
