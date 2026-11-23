A personal fork of [glanceapp/glance](https://github.com/glanceapp/glance).

## Custom changes

### OpenID Connect (SSO) login

Optional OIDC authentication alongside the existing local username/password login. Authorization code flow with PKCE, signed state cookie, nonce validation and ID token verification; a successful login issues the same session cookie as a local login. Supports OIDC-only setups, combined setups, `auto-login`, `disable-local-login` and an explicit `redirect-url`.

```yaml
auth:
  secret-key: ${secret:auth-secret}
  oidc:
    issuer: https://auth.example.com/realms/home
    client-id: glance
    client-secret: ${secret:oidc-client-secret}
```

See [authentication docs](docs/configuration.md#openid-connect-oidc).

### `Auto` theme

The theme picker has an `Auto` preset that follows the operating system's dark/light preference, switching between `default-dark` and `default-light`. See [theme docs](docs/configuration.md#theme).

### Releases widget: Gitea and self-hosted GitLab/Gitea

- New `gitea:` prefix with a Gitea source icon.
- `gitea:` and `gitlab:` repositories accept a full URL, so self-hosted instances work.

```yaml
- type: releases
  repositories:
    - gitea:gitea/runner
    - gitea:https://gitea.example.com:3000/owner/repo
    - gitlab:https://gitlab.example.com/group/project
```

## Run it

```bash
docker pull ghcr.io/felix-schindler/glance:latest
```

`latest` is the only tag (no versions). The image is rebuilt from `main` on every push and supports `linux/amd64` and `linux/arm64`. There are no release binaries.
