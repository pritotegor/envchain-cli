# envchain-cli

> Manage and inject environment variable sets per project with optional secret store integration.

---

## Installation

```bash
go install github.com/yourname/envchain-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain-cli.git
cd envchain-cli && go build -o envchain .
```

---

## Usage

**Define an environment chain for a project:**

```bash
envchain set myproject AWS_REGION=us-east-1 APP_ENV=production
```

**Inject variables when running a command:**

```bash
envchain run myproject -- ./start-server
```

**List all configured chains:**

```bash
envchain list
```

**Remove a chain or a specific variable from a chain:**

```bash
# Remove a single variable from a chain
envchain unset myproject AWS_REGION

# Remove an entire chain
envchain unset myproject
```

**Integrate with a secret store (e.g., Vault, AWS SSM):**

```bash
envchain set myproject --source vault secret/myproject/config
```

Variables are stored per named chain and injected into the subprocess environment at runtime. Secret store references are resolved lazily on each `run`.

---

## Configuration

By default, chains are stored in `~/.envchain/chains.json`. Override the path with:

```bash
export ENVCHAIN_CONFIG=/path/to/config.json
```

---

## License

MIT © yourname
