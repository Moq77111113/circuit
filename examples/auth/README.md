# Authentication Example

This example demonstrates Circuit's authentication modes.

## Quick Start - Basic Auth

The simplest way to protect your Circuit UI:

```bash
cd examples/auth

# Run with Basic Auth (default credentials: admin/secret)
go run main.go -auth=basic
```

Open http://localhost:8080 in your browser. You'll be prompted for credentials:
- Username: `admin`
- Password: `secret`

### Custom Credentials

```bash
go run main.go -auth=basic -user=alice -pass=mypassword
```

### Production: Use Argon2id Hashes

For production, **never use plaintext passwords**. Generate an argon2id hash:

```bash
# Install argon2 CLI (or use Go code)
echo -n "mypassword" | argon2 somesalt -id -t 3 -m 16 -p 4 -l 32 -e
# Output: $argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
```

Then use it:

```bash
go run main.go -auth=basic -user=admin -pass='$argon2id$v=19$...'
```

## Forward Auth (Behind a Proxy)

Use this when you **already have** OAuth2 Proxy, Traefik, or Cloudflare Access deployed.

### Test Locally with curl

```bash
# Terminal 1: Start Circuit with Forward Auth
go run main.go -auth=forward

# Terminal 2: Test with header
curl -H "X-Forwarded-User: alice@example.com" http://localhost:8080
# → Returns Circuit UI

curl http://localhost:8080
# → 401 Unauthorized (missing header)
```

### CLI Flags

```
-auth string
    Auth mode: none, basic, or forward (default "none")
-user string
    Username for basic auth (default "admin")
-pass string
    Password for basic auth - plaintext or argon2id hash (default "secret")
-port string
    HTTP server port (default "8080")
```

### Production Setup with OAuth2 Proxy

1. Deploy OAuth2 Proxy in front of your app:

```yaml
# oauth2-proxy config
upstreams:
  - http://your-app:8080
provider: github  # or google, gitlab, etc
client_id: your-github-oauth-app-id
client_secret: your-github-oauth-app-secret
cookie_secret: random-32-byte-string
email_domains:
  - yourcompany.com
```

2. Configure your app:

```go
auth := circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
    "email": "X-Forwarded-Email",
    "groups": "X-Forwarded-Groups",
})
handler, _ := circuit.From(&cfg, circuit.WithAuth(auth))
```

3. The flow:
   - User visits `/admin` → OAuth2 Proxy intercepts
   - Not logged in? → Redirect to GitHub OAuth
   - After login → Proxy sets `X-Forwarded-User` header
   - Circuit validates header → Grants access

## No Auth (Development Only)

```bash
# Default mode (no -auth flag)
go run main.go

# Or explicitly
go run main.go -auth=none
```

**⚠️ WARNING:** UI is completely open. Only use on localhost or internal networks.

## Security Notes

1. **Always use HTTPS in production** - Basic Auth sends credentials in base64
2. **Never commit passwords** - Use environment variables or secret managers
3. **Use argon2id hashes** - Not plaintext (except dev/testing)
4. **Restrict network access** - Use firewalls, VPNs, or internal-only networks

## Integration Examples

### With Echo framework

```go
e := echo.New()

auth := circuit.NewBasicAuth("admin", "secret")
circuitHandler, _ := circuit.From(&cfg, circuit.WithAuth(auth))

e.Any("/admin/*", echo.WrapHandler(circuitHandler))
e.Start(":8080")
```

### With Gin framework

```go
r := gin.Default()

auth := circuit.NewBasicAuth("admin", "secret")
circuitHandler, _ := circuit.From(&cfg, circuit.WithAuth(auth))

r.Any("/admin/*", gin.WrapH(circuitHandler))
r.Run(":8080")
```

### With Chi router

```go
r := chi.NewRouter()

auth := circuit.NewBasicAuth("admin", "secret")
circuitHandler, _ := circuit.From(&cfg, circuit.WithAuth(auth))

r.Mount("/admin", circuitHandler)
http.ListenAndServe(":8080", r)
```
