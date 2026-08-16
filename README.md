# tagi
> Fast URL Scope Filtering & Reconnaissance Utility

```text
  __                 .__
_/  |______     ____ |__|
\   __\__  \   / ___\|  |
 |  |  / __ \_/ /_/  >  |
 |__| (____  /\___  /|__|
           \//_____/     

```

`tagi` is a fast, lightweight command-line utility written in Go designed for filtering large streams of URLs based on in-scope and out-of-scope domains. It reads URLs from standard input (STDIN), processes them to ensure they match your defined scope (including subdomains), automatically removes duplicates, and prints the filtered results to standard output (STDOUT).

It is built for bug bounty reconnaissance pipelines, integrating seamlessly with tools like `katana`, `waybackurls`, `gau`, and `hakrawler`.

## Features

* **Pipeline Ready:** Designed to accept STDIN from other recon tools and output clean results to STDOUT.
* **Smart Scope Matching:** Automatically matches both the exact domain and its subdomains (e.g., scoping `example.com` will match `api.example.com`).
* **Out-of-Scope (OOS) Filtering:** Exclude specific domains or subdomains from your final output to avoid out-of-bounds testing.
* **Built-in Deduplication:** Automatically tracks and drops duplicate URLs on the fly to keep your output clean and efficient.
* **Highly Performant:** Utilizes buffered I/O and optimized mapping for rapid processing of massive text streams.

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed on your system. You can install `tagi` directly using `go install`:

```bash
go install github.com/INTELEON404/tagi@latest
```

*Ensure your `$GOPATH/bin` directory is added to your system's `$PATH` to run the command globally.*

## Usage

```bash
command | tagi [options]

```

### Flags

| Flag | Description |
| --- | --- |
| `--scope <host>` | Define a single in-scope hostname (e.g., `example.com`). |
| `--scopelist <file>` | Path to a file containing a list of in-scope hostnames. |
| `--ooscope <host>` | Define a single out-of-scope hostname to exclude (e.g., `admin.example.com`). |
| `--ooslist <file>` | Path to a file containing a list of out-of-scope hostnames. |
| `--help` | Display the help menu and examples. |

*Note: You cannot use `--scope` and `--scopelist` simultaneously.*

## Examples

### 1. Basic Single Scope

Filter URLs for a single domain and its subdomains:

```bash
cat urls.txt | tagi --scope example.com

```

### 2. Using a Scope List

Filter URLs against a text file containing multiple in-scope domains:

```bash
katana -list domains.txt | tagi --scopelist scope.txt

```

### 3. Scope with a Specific Exclusion

Keep everything in `example.com`, but strictly exclude `admin.example.com`:

```bash
waybackurls example.com | tagi --scope example.com --ooscope admin.example.com

```

### 4. Advanced: Scope and Out-of-Scope Lists

Filter using comprehensive lists for both allowed and excluded scopes, writing the output to a file:

```bash
cat all_urls.txt | tagi --scopelist scope.txt --ooslist out_of_scope.txt > valid_targets.txt

```

## Input/Output File Formatting

When providing a file to `--scopelist` or `--ooslist`, list one hostname per line. You can also use `#` for comments or leave blank lines for structural readability. The tool will parse and ignore these automatically.

**Example `scope.txt`:**

```text
# Main target
example.com

# Acquisitions
target-acquisition.com

```

## Disclaimer

> [!WARNING]
> **This tool is for educational and authorized security testing purposes only.**
> The author and contributors assume no liability and are not responsible for any misuse, damage, or legal consequences resulting from the use of this tool.
> Users must ensure they have explicit, written authorization from the target system's owners prior to conducting any reconnaissance or security testing operations.

## License

Distributed under the MIT License. See the  [LICENSE](LICENSE) file for more information.
