# tagi

> Fast URL Scope Filtering & Reconnaissance Utility

```text
  __                 .__
_/  |______     ____ |__|
\   __\__  \   / ___\|  |
 |  |  / __ \_/ /_/  >  |
 |__| (____  /\___  /|__|
           \//____/
```

`tagi` is a fast, lightweight command-line utility written in Go for filtering URL streams against an authorized bug-bounty scope.

It is designed for reconnaissance pipelines and works cleanly with tools such as `katana`, `waybackurls`, `gau`, `httpx`, and other tools that output URLs or hosts to STDOUT.

## Features

* **Fast & Lightweight** — Go standard library only.
* **Pipeline Friendly** — Reads from STDIN and writes matched URLs to STDOUT.
* **Scope Matching** — Matches an exact hostname and its valid subdomains.
* **Multiple Scopes** — Load multiple in-scope domains from a file.
* **Boundary Safe** — Prevents false matches such as `evil-example.com` matching `example.com`.
* **Case Insensitive** — Hostnames are normalized before matching.
* **Schemeless Support** — Handles inputs such as `example.com/path`.
* **Port Aware** — Correctly handles URLs such as `example.com:8080`.
* **Optional Deduplication** — Disabled by default; enable with `-d`.
* **Quiet Mode** — Keeps pipeline output clean.
* **Verbose Statistics** — Optional statistics are written to STDERR.
* **No External Dependencies** — Standard library only.

## Installation

### Go

```bash
go install github.com/INTELEON404/tagi@latest
```

Make sure your Go binary directory is in your `$PATH`.

### Build from Source

```bash
git clone https://github.com/INTELEON404/tagi.git
cd tagi
go build -o tagi .
```

## Usage

```bash
command | tagi -s <host|file> [options]
```

## Options

| Flag                       | Description                                  |
| -------------------------- | -------------------------------------------- |
| `-s, --scope <host\|file>` | In-scope hostname or file containing hosts   |
| `-d, --dedupe`             | Remove duplicate URLs                        |
| `-q, --quiet`              | Suppress non-essential output and statistics |
| `-v, --verbose`            | Print processing statistics to STDERR        |
| `-h, --help`               | Show help and examples                       |
| `--version`                | Show version                                 |

> `--version` is the only version flag. There is no `-V` option.

## Examples

### Single Scope

```bash
cat urls.txt | tagi -s example.com
```

### Scope File

```bash
cat urls.txt | tagi -s scope.txt
```

### Deduplicate Results

```bash
cat urls.txt | tagi -s example.com -d
```

### Quiet Pipeline

```bash
cat urls.txt | tagi -s example.com -q
```

### Verbose Statistics

```bash
cat urls.txt | tagi -s example.com -v
```

### Bug Bounty Recon Pipeline

```bash
katana -list domains.txt | tagi -s scope.txt -d
```

```bash
gau example.com | tagi -s example.com -d
```

```bash
waybackurls example.com | tagi -s example.com -d
```

## Scope File

The scope file contains one hostname per line.

```text
# Main target
example.com

# Additional authorized target
target.org
```

Blank lines and lines beginning with `#` are ignored.

## Matching Behavior

If the scope is:

```text
example.com
```

These are matched:

```text
example.com
www.example.com
api.example.com
a.b.example.com
```

These are **not** matched:

```text
evil-example.com
notexample.com
example.com.au
example.org
```

Hostnames are matched case-insensitively.

For example:

```text
https://EXAMPLE.COM/path
```

matches:

```text
example.com
```

Ports are handled correctly:

```text
https://example.com:8443/api
```

matches:

```text
example.com
```

## Input & Output

`tagi` is designed for Unix-style pipelines.

**STDIN**

Receives URLs or hosts:

```text
https://example.com/
https://api.example.com/login
https://external.example.org/
```

**STDOUT**

Only matched input lines are printed:

```text
https://example.com/
https://api.example.com/login
```

Statistics and errors are written to **STDERR**, keeping STDOUT safe for further pipeline processing.

## Deduplication

Deduplication is **disabled by default**.

Enable it with:

```bash
tagi -s example.com -d
```

The original input URL is preserved when printed; `tagi` does not rewrite or normalize the output URL.

## Version

```bash
tagi --version
```

Example:

```text
tagi 1.0.0
```

## Verification

Run the following before releasing:

```bash
gofmt -w .
go vet ./...
go test ./...
go build -o tagi .
```

## Disclaimer

> [!WARNING]
> `tagi` is intended for authorized security research, bug-bounty reconnaissance, and educational use.
>
> Only process targets that are explicitly within the scope of the relevant bug-bounty program or for which you have authorization to perform security testing.
>
> The author and contributors are not responsible for misuse of the software.

## License

Distributed under the MIT License. See [`LICENSE`](LICENSE) for details.
