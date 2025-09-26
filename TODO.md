# ToDo

## **Sprint 1: Foundations & User Flow**

**Sprint Goal:** a user can create an account, verify it, log in/out, and reset password. discovery + JWKS live.

### Flow: Environment & Setup

* [ ] Add Redis & SMTP sink (MailHog) to docker-compose.
* [ ] Centralize app config (issuer URL, cookie flags, TTLs).
* [ ] Generate signing keypair (RS256/ES256), expose JWKS + discovery doc.

### Flow: User Lifecycle (Signup/Login)

* [ ] DB migrations for users, clients, sessions, auth codes, tokens, audit.
* [ ] Implement user signup with email verification (via SMTP sink).
* [ ] Implement login/logout with secure cookies (HttpOnly, Secure, SameSite).
* [ ] Implement password reset with single-use, time-limited token.
* [ ] Write audit log entries for all user lifecycle events.

---

## **Sprint 2: OAuth2/OIDC Core Flow**

**Sprint Goal:** a registered client can perform Auth Code + PKCE login and read `/userinfo`.

### Flow: Authorization

* [ ] Implement `/authorize` endpoint with PKCE support and `nonce` param.
* [ ] Build basic consent UI (show client + scopes, remember decision).

### Flow: Token Issuance

* [ ] Implement `/token` endpoint (code → ID token + access token).
* [ ] Implement refresh tokens with rotation + one-time use.
* [ ] Implement `/userinfo` endpoint.

### Flow: Client Registration

* [ ] Build minimal admin UI/API to register OIDC clients.
* [ ] Enforce exact redirect URI matching (no wildcards, HTTPS only).

---

## **Sprint 3: Integration & Hardening**

**Sprint Goal:** apps integrate with your IdP; system resists abuse and exposes metrics.

### Flow: Integration

* [ ] Build sample **web app (BFF)** using Auth Code + PKCE.
* [ ] Build sample **service-to-service** using Client Credentials flow.

### Flow: Security Hardening

* [ ] Argon2id password hashing with breach-check on signup/reset.
* [ ] Add rate limits to `/login`, `/authorize`, `/token`.
* [ ] Validate token claims (`iss`, `aud`, `exp`, `iat`, `nonce`).
* [ ] Implement lockouts/backoff for failed logins.

### Flow: Observability

* [ ] Add structured logging with request IDs.
* [ ] Add Prometheus metrics (auth latency, token issuance, error rates).
* [ ] Add OpenTelemetry tracing for `/authorize → /token`.

---

## **Stretch Sprint: Extras**

* [ ] Session dashboard (list + revoke active sessions/devices).
* [ ] MFA v1 (TOTP + backup codes).
* [ ] Key rotation playbook + runbook.
* [ ] SLOs + alerts (latency, token failures, queue depth).
