package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

const version = "1.2"

func main() {
	os.Args = preprocessArgs(os.Args)

	var (
		scope   = flag.String("scope", "", "")
		quiet   = flag.Bool("quiet", false, "")
		dedupe  = flag.Bool("dedupe", false, "")
		verbose = flag.Bool("verbose", false, "")
		help    = flag.Bool("help", false, "")
		ver     = flag.Bool("version", false, "")
	)
	flag.Usage = func() { usage() }
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}
	if *ver {
		fmt.Println("tagi", version)
		os.Exit(0)
	}
	if *scope == "" {
		fmt.Fprintln(os.Stderr, "tagi: error: no scope specified")
		if !*quiet {
			usage()
		}
		os.Exit(1)
	}

	scopes, err := loadScope(*scope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tagi: error: %v\n", err)
		os.Exit(1)
	}

	if err := process(os.Stdin, os.Stdout, os.Stderr, scopes, *dedupe, *verbose, *quiet); err != nil {
		fmt.Fprintf(os.Stderr, "tagi: error: %v\n", err)
		os.Exit(1)
	}
}

func preprocessArgs(args []string) []string {
	if len(args) <= 1 {
		return args
	}

	shorts := map[string]string{
		"s": "scope",
		"q": "quiet",
		"d": "dedupe",
		"v": "verbose",
		"h": "help",
	}

	out := make([]string, 0, len(args))
	out = append(out, args[0])

	for i, arg := range args[1:] {
		if arg == "--" {
			out = append(out, arg)
			out = append(out, args[i+2:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			out = append(out, "-"+arg[2:])
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			body := arg[1:]
			if idx := strings.Index(body, "="); idx >= 0 {
				key, val := body[:idx], body[idx:]
				if long, ok := shorts[key]; ok {
					out = append(out, "-"+long+val)
					continue
				}
			} else {
				if long, ok := shorts[body]; ok {
					out = append(out, "-"+long)
					continue
				}
			}
		}
		out = append(out, arg)
	}

	return out
}

func usage() {
	fmt.Fprintln(os.Stderr, "tagi - URL scope filter")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage: tagi -s <host|file> [options]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -s, --scope <host|file>  In-scope hostname or file containing hosts")
	fmt.Fprintln(os.Stderr, "  -d, --dedupe             Remove duplicate URLs")
	fmt.Fprintln(os.Stderr, "  -q, --quiet              Suppress usage and statistics")
	fmt.Fprintln(os.Stderr, "  -v, --verbose            Print statistics to stderr")
	fmt.Fprintln(os.Stderr, "  -h, --help               Show this help message")
	fmt.Fprintln(os.Stderr, "      --version            Show version")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  cat urls.txt | tagi -s example.com")
	fmt.Fprintln(os.Stderr, "  cat urls.txt | tagi -s scope.txt -d")
	fmt.Fprintln(os.Stderr, "  cat urls.txt | tagi -s scope.txt -q")
}


func loadScope(arg string) (map[string]struct{}, error) {
	scopes := make(map[string]struct{})

	if strings.ContainsAny(arg, `/\`) {
		f, err := os.Open(arg)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return parseScopeFile(f, scopes)
	}

	if fi, err := os.Stat(arg); err == nil && fi.Mode().IsRegular() {
		f, err := os.Open(arg)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return parseScopeFile(f, scopes)
	}

	host := normalizeHost(arg)
	if host == "" {
		return nil, fmt.Errorf("invalid scope: %s", arg)
	}
	scopes[host] = struct{}{}
	return scopes, nil
}

func parseScopeFile(r io.Reader, dest map[string]struct{}) (map[string]struct{}, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if h := normalizeHost(line); h != "" {
			dest[h] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(dest) == 0 {
		return nil, fmt.Errorf("no valid hosts in scope file")
	}
	return dest, nil
}

func normalizeHost(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimLeft(s, ".")
	s = strings.TrimRight(s, ".")
	return s
}

func inScope(host string, scopes map[string]struct{}) bool {
	if _, ok := scopes[host]; ok {
		return true
	}
	for i := 0; i < len(host); i++ {
		if host[i] == '.' {
			if _, ok := scopes[host[i+1:]]; ok {
				return true
			}
		}
	}
	return false
}

func process(r io.Reader, w io.Writer, errW io.Writer, scopes map[string]struct{}, dedupe, verbose, quiet bool) error {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var seen map[string]struct{}
	if dedupe {
		seen = make(map[string]struct{})
	}

	var total, matched, dupes int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		total++

		parseURL := line
		if !strings.Contains(line, "://") && !strings.HasPrefix(line, "//") {
			parseURL = "http://" + line
		}

		u, err := url.Parse(parseURL)
		if err != nil {
			continue
		}

		host := normalizeHost(u.Hostname())
		if host == "" {
			continue
		}

		if !inScope(host, scopes) {
			continue
		}

		if dedupe {
			if _, ok := seen[line]; ok {
				dupes++
				continue
			}
			seen[line] = struct{}{}
		}

		matched++
		fmt.Fprintln(w, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stdin error: %v", err)
	}

	if verbose && !quiet {
		fmt.Fprintf(errW, "tagi: total=%d matched=%d", total, matched)
		if dedupe {
			fmt.Fprintf(errW, " dupes=%d", dupes)
		}
		fmt.Fprintln(errW)
	}
	return nil
}
