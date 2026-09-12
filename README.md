# tagi

<video src="https://github.com/user-attachments/assets/44b2a54b-b959-44fa-9911-b5d700cfc9e4" controls width="800"></video>

`tagi` is a lightweight, fast command-line utility for filtering URL streams against an authorized bug-bounty scope. Built with Go and designed for reconnaissance pipelines, it integrates seamlessly with tools like `katana`, `waybackurls`, `gau`, and `httpx`.

## Features

- **Fast & Minimal** — Go standard library only, no external dependencies
- **Pipeline Friendly** — Reads STDIN, writes matched URLs to STDOUT
- **Exact & Subdomain Matching** — Match full domain trees or exact hostnames only
- **Multiple Scopes** — Load domains from a file or pass a single hostname
- **Boundary Safe** — Prevents false positives like `evil-example.com` matching `example.com`
- **Case Insensitive** — Hostnames normalized before matching
- **Schemeless Support** — Accepts `example.com/path` as input
- **Port Aware** — Correctly handles `example.com:8080`
- **Deduplication** — Optional URL deduplication via `-d` flag
- **Clean Output** — Statistics written to STDERR, STDOUT reserved for data

## Installation

### Pre-built Binary

```bash
go install github.com/inteleon404/tagi@latest
```

Ensure your Go binary directory is in your `$PATH`.

### Build from Source

```bash
git clone https://github.com/inteleon404/tagi.git
cd tagi
go build -o tagi .
```

## Usage

```bash
cat urls.txt | tagi -s <host|file> [options]
```

## Options

| Flag | Description |
|------|-------------|
| `-s, --scope <host\|file>` | In-scope hostname or path to scope file |
| `-d, --dedupe` | Remove duplicate URLs from output |
| `--no-sub` | Match exact hostname only (exclude all subdomains) |
| `-q, --quiet` | Suppress statistics and non-critical messages |
| `-v, --verbose` | Print processing statistics to STDERR |
| `-h, --help` | Show usage and examples |
| `--version` | Display version information |

## Scope Matching

### Default Behavior (Subdomain Matching)

By default, `tagi` matches a scope hostname and all its subdomains.

**Scope:**
```
example.com
```

**Matched:**
```
example.com
www.example.com
api.example.com
blog.example.com
deep.nested.example.com
```

**Not Matched:**
```
evil-example.com
notexample.com
example.com.au
example.org
```

### Exact Hostname Matching (`--no-sub`)

With the `--no-sub` flag, `tagi` matches the scope hostname exactly and excludes all subdomains.

**Scope:**
```
example.com
```

**Matched:**
```
example.com
```

**Not Matched:**
```
www.example.com
api.example.com
blog.example.com
deep.nested.example.com
evil-example.com
notexample.com
example.com.au
example.org
```

### Case Sensitivity

Hostnames are normalized to lowercase before matching.

```
https://EXAMPLE.COM/path     → matches example.com
https://WWW.Example.Com/page → matches example.com (default)
```

### Ports

Ports are correctly parsed and do not affect matching.

```
https://example.com:8443/api → matches example.com
https://example.com:3000     → matches example.com
```

## Input & Output

**STDIN:** Accepts URLs or hostnames, one per line.

```
https://example.com/
https://api.example.com/login
https://external.org/
```

**STDOUT:** Matched input lines only.

```
https://example.com/
https://api.example.com/login
```

**STDERR:** Statistics, errors, and diagnostic messages (does not interfere with pipeline data).

## Deduplication

Deduplication is disabled by default. Enable with `-d`.

```bash
cat urls.txt | tagi -s example.com -d
```

Output preserves original URL formatting; `tagi` does not normalize URLs.

## Scope File Format

Scope files contain one hostname per line. Blank lines and lines starting with `#` are ignored.

```
# Primary target
example.com

# Secondary target
target.org

# Out of scope
exclude.com
```

## Examples

### Basic Usage

```bash
cat urls.txt | tagi -s example.com
```

### Scope File

```bash
cat urls.txt | tagi -s scope.txt
```

### Exact Hostname Matching

```bash
cat urls.txt | tagi -s example.com --no-sub
```

### With Deduplication

```bash
cat urls.txt | tagi -s example.com -d
```

### Quiet Output

```bash
cat urls.txt | tagi -s example.com -q
```

### Statistics

```bash
cat urls.txt | tagi -s example.com -v
```

Output example:
```
tagi: total=1000 matched=487
```

### Bug Bounty Pipelines

```bash
katana -list domains.txt | tagi -s scope.txt -d
```

```bash
gau example.com | tagi -s example.com -d
```

```bash
waybackurls example.com | tagi -s example.com -d
```

## Version

```bash
tagi --version
```

Output:
```
tagi 1.3
```

## Development

Before releasing, run:

```bash
gofmt -w .
go vet ./...
go test ./...
go build -o tagi .
```

## Disclaimer

> [!WARNING]
> `tagi` is for authorized security research, bug-bounty testing, and educational use only.
>
> Only process targets that are explicitly in scope or for which you have written authorization.
>
> Unauthorized access to computer systems is illegal. The author and contributors are not responsible for misuse.

## License

MIT License. See [`LICENSE`](LICENSE) for details.
