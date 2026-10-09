// Command gen adds Material Symbols Rounded icons to symbols.go.
package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	version = "v0.27.0"
	baseURL = "https://raw.githubusercontent.com/marella/material-symbols/" +
		version + "/svg/500/rounded/"
)

var declaration = regexp.MustCompile(`(?m)^\s*([A-Za-z][A-Za-z0-9]*)\s*=\s*mustParse\("`)
var namePattern = regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)

type svg struct {
	ViewBox string `xml:"viewBox,attr"`
	Path    []struct {
		D string `xml:"d,attr"`
	} `xml:"path"`
}

func main() {
	var names string
	flag.StringVar(&names, "icons", "", "comma-separated Material Symbols names")
	flag.Parse()

	list, err := requested(names, flag.Args())
	if err != nil {
		fail(err)
	}
	if len(list) == 0 {
		var err error
		list, err = readList("icons.txt")
		if err != nil {
			fail(err)
		}
	}
	if len(list) == 0 {
		return
	}

	path := "symbols.go"
	data, err := os.ReadFile(path)
	if err != nil {
		fail(err)
	}
	known := make(map[string]bool)
	for _, match := range declaration.FindAllSubmatch(data, -1) {
		known[string(match[1])] = true
	}

	added := make([]string, 0, len(list))
	for _, name := range list {
		goName, err := identifier(name)
		if err != nil {
			fail(err)
		}
		if known[goName] {
			continue
		}
		d, err := download(name)
		if err != nil {
			fail(err)
		}
		added = append(added, fmt.Sprintf("\t%s = mustParse(%q)", goName, d))
		known[goName] = true
	}
	if len(added) == 0 {
		return
	}
	sort.Strings(added)
	end := bytes.LastIndex(data, []byte("\n)"))
	if end < 0 {
		fail(errors.New("symbols.go: could not find the end of the var block"))
	}
	var out bytes.Buffer
	out.Write(data[:end])
	out.WriteByte('\n')
	for _, line := range added {
		out.WriteString(line)
		out.WriteByte('\n')
	}
	out.Write(data[end:])
	if err := os.WriteFile(path, out.Bytes(), 0o644); err != nil {
		fail(err)
	}
	fmt.Printf("added %d icon(s) to %s\n", len(added), filepath.ToSlash(path))
}

func requested(commaSeparated string, args []string) ([]string, error) {
	var names []string
	if commaSeparated != "" {
		names = append(names, strings.Split(commaSeparated, ",")...)
	}
	names = append(names, args...)
	seen := make(map[string]bool)
	out := names[:0]
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if !namePattern.MatchString(name) {
			return nil, fmt.Errorf("invalid icon name %q", name)
		}
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}

func readList(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if line != "" {
			names = append(names, line)
		}
	}
	return requested("", names)
}

func identifier(name string) (string, error) {
	parts := strings.Split(name, "_")
	var out strings.Builder
	for _, part := range parts {
		if part == "" {
			return "", fmt.Errorf("invalid icon name %q", name)
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		out.WriteString(string(runes))
	}
	if out.Len() == 0 {
		return "", fmt.Errorf("invalid icon name %q", name)
	}
	return out.String(), nil
}

func download(name string) (string, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(baseURL + name + ".svg")
	if err != nil {
		return "", fmt.Errorf("download %s: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: HTTP %s", name, resp.Status)
	}
	var doc svg
	if err := xml.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return "", fmt.Errorf("parse %s.svg: %w", name, err)
	}
	if doc.ViewBox != "0 -960 960 960" {
		return "", fmt.Errorf("%s.svg: unexpected viewBox %q", name, doc.ViewBox)
	}
	if len(doc.Path) != 1 || doc.Path[0].D == "" {
		return "", fmt.Errorf("%s.svg: expected exactly one non-empty path", name)
	}
	if strings.ContainsAny(doc.Path[0].D, "\r\n") {
		return "", fmt.Errorf("%s.svg: path contains a newline", name)
	}
	return doc.Path[0].D, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "icon generator:", err)
	os.Exit(1)
}
