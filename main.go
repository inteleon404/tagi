package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	banner = `  __                 .__
_/  |______     ____ |__|
\   __\__  \   / ___\|  |
 |  |  / __ \_/ /_/  >  |
 |__| (____  /\___  /|__|
           \//_____/     
`
)

const (
	reset  = "\033[0m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			printBanner()
			printHelp()
			os.Exit(0)
		}
	}

	var (
		scope     = flag.String("scope", "", "")
		scopelist = flag.String("scopelist", "", "")
		ooscope   = flag.String("ooscope", "", "")
		ooslist   = flag.String("ooslist", "", "")
	)
	flag.Usage = func() {
		printBanner()
		printConciseUsage()
	}
	flag.Parse()

	if *scope == "" && *scopelist == "" {
		printBanner()
		fmt.Fprintln(os.Stderr, cyan+"[INF]"+reset+" no scope specified")
		fmt.Fprintln(os.Stderr)
		printConciseUsage()
		os.Exit(1)
	}

	if *scope != "" && *scopelist != "" {
		fmt.Fprintln(os.Stderr, red+"[ERR]"+reset+" --scope and --scopelist cannot be used together")
		os.Exit(1)
	}

	inScope := make(map[string]struct{})
	if *scope != "" {
		if h := normalizeHost(*scope); h != "" {
			inScope[h] = struct{}{}
		}
	}
	if *scopelist != "" {
		if err := loadHostFile(*scopelist, inScope); err != nil {
			fmt.Fprintf(os.Stderr, red+"[ERR]"+reset+" failed to read scope file: %v\n", err)
			os.Exit(1)
		}
	}

	outOfScope := make(map[string]struct{})
	if *ooscope != "" {
		if h := normalizeHost(*ooscope); h != "" {
			outOfScope[h] = struct{}{}
		}
	}
	if *ooslist != "" {
		if err := loadHostFile(*ooslist, outOfScope); err != nil {
			fmt.Fprintf(os.Stderr, red+"[ERR]"+reset+" failed to read out-of-scope file: %v\n", err)
			os.Exit(1)
		}
	}

	processInput(inScope, outOfScope)
}

func printBanner() {
	fmt.Fprint(os.Stderr, cyan+banner+reset)
	fmt.Fprintln(os.Stderr, "tagi - Fast URL Scope Filtering & Reconnaissance Utility")
	fmt.Fprintln(os.Stderr)
}

func printConciseUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  command | tagi [options]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Run:")
	fmt.Fprintln(os.Stderr, "  tagi --help")
}

func printHelp() {
	help := `FLAGs:
  --scope <host>       Single in-scope hostname
  --scopelist <file>   In-scope hostname list
  --ooscope <host>     Single out-of-scope hostname
  --ooslist <file>     Out-of-scope hostname list
  --help               Show help

Examples:
  katana -list domains.txt | tagi --scope example.com
  katana -list domains.txt | tagi --scopelist scope.txt
  katana -list domains.txt | tagi --scope example.com --ooscope admin.example.com
  katana -list domains.txt | tagi --scopelist scope.txt --ooslist oos.txt

Input:
  URLs are read from STDIN.

Output:
  Only in-scope URLs are written to STDOUT.
`
	fmt.Fprint(os.Stderr, help)
}

func loadHostFile(filename string, dest map[string]struct{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if h := normalizeHost(line); h != "" {
			dest[h] = struct{}{}
		}
	}
	return scanner.Err()
}

func normalizeHost(host string) string {
	return strings.ToLower(strings.TrimSpace(host))
}

func processInput(inScope, outOfScope map[string]struct{}) {
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		rawURL := strings.TrimSpace(scanner.Text())
		if rawURL == "" {
			continue
		}

		u, err := url.Parse(rawURL)
		if err != nil {
			continue
		}

		host := normalizeHost(u.Hostname())
		if host == "" {
			continue
		}

		if isOutOfScope(host, outOfScope) {
			continue
		}

		if !isInScope(host, inScope) {
			continue
		}

		if _, ok := seen[rawURL]; ok {
			continue
		}
		seen[rawURL] = struct{}{}

		fmt.Println(rawURL)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, yellow+"[!]"+reset+" Warning: stdin read error: %v\n", err)
		os.Exit(1)
	}
}

func isInScope(host string, scopes map[string]struct{}) bool {
	for scope := range scopes {
		if matchHost(host, scope) {
			return true
		}
	}
	return false
}

func isOutOfScope(host string, oos map[string]struct{}) bool {
	for scope := range oos {
		if matchHost(host, scope) {
			return true
		}
	}
	return false
}

func matchHost(host, scope string) bool {
	if host == scope {
		return true
	}
	return strings.HasSuffix(host, "."+scope)
}
