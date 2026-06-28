# Run the server:
```bash
go run ./cmd/server
```

# Run the client:
```bash
go run ./cmd/client
```

# Dev script for faster testing:
```bash
./scripts/dev.sh
```

# Env variables:

| Var             | Used by | Default                  | Values / notes              |
|-----------------|---------|--------------------------|-----------------------------|
| `FYLINE_LOG`    | both    | `info`                   | `debug` \| `info` \| `warn` |
| `FYLINE_USER`   | client  | `Catalin`                | username (or `-user` flag)  |
| `FYLINE_SERVER` | client  | `ws://localhost:8080/ws` | server to connect to        |
| `FYLINE_ADDR`   | server  | `:8080`                  | listen address              |

Example:
```bash
FYLINE_LOG=debug FYLINE_USER=Bob go run ./cmd/client
```
```bash
FYLINE_LOG=debug go run ./cmd/server
```

# Flags

| Flag    | Used by | What it does                          |
|---------|---------|---------------------------------------|
| `-user` | client  | username (overrides `FYLINE_USER`)    |

Example:
```bash
go run ./cmd/client -user Bob
```
