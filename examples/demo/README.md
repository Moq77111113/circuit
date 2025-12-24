# Circuit Demo (Railway-ready)

This demo is a small, production-shaped Go service that embeds Circuit as an **in-process control surface**.

It demonstrates:

- Config editing via Circuit UI (`/admin`) backed by YAML on disk
- Nested structs + slices (add/remove items in the UI)
- Validation-style tags (`required`, `min/max/step`) and rich UI types (`select`, `checkbox`, `number`, `text`)
- Live config application via `WithOnChange` (the API handler reads a snapshot)
- Pluggable auth (`none`, `basic`, `forward`) via environment variables

## Run locally

From the repo root:

```bash
go run ./examples/demo
```

Open:

- `http://localhost:8080/admin` (Circuit UI)
- `http://localhost:8080/api` (returns JSON using the current config)

Auth (local defaults):

- If `CIRCUIT_AUTH_MODE=basic` and no credentials are set, it uses `admin/admin`.

### Forward auth (behind an OAuth proxy)

If you already have an auth proxy that injects headers:

- `CIRCUIT_AUTH_MODE=forward`
- `CIRCUIT_FORWARD_SUBJECT_HEADER=X-Forwarded-User`
- `CIRCUIT_FORWARD_EMAIL_HEADER=X-Forwarded-Email`

Circuit will validate requests based on those headers; your proxy handles the OAuth login flow.

## What to try in the UI

- Change `ops.maintenance` to `true` → `/api` returns 503, but `/admin` still works.
- Add/remove entries from `backends` and `flags` slices.
- Adjust `limits.artificial_delay_ms` to simulate latency.
- Set `http.allowed_cidrs` to restrict `/api` by IP.
