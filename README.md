# envdiff

> CLI tool to diff and reconcile `.env` files across environments, with secret masking and CI integration support.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envdiff/releases).

---

## Usage

Compare two `.env` files and highlight differences:

```bash
envdiff .env.development .env.production
```

Mask secret values in output (keys containing `SECRET`, `KEY`, `TOKEN`, etc.):

```bash
envdiff --mask .env.local .env.staging
```

Reconcile missing keys from one file into another:

```bash
envdiff reconcile --source .env.example --target .env.local
```

Use in CI pipelines to fail on unexpected drift:

```bash
envdiff --ci --strict .env.example .env.production
```

**Example output:**

```
~ DATABASE_URL   changed
+ NEW_FEATURE_FLAG   only in production
- DEBUG_MODE   only in development
  API_TIMEOUT   match
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--mask` | Redact sensitive values in diff output |
| `--ci` | Exit with non-zero code if differences are found |
| `--strict` | Treat missing keys as errors |
| `--json` | Output results as JSON |

---

## License

[MIT](LICENSE) © 2024 yourusername