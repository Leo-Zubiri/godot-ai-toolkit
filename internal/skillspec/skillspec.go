// Package skillspec parses and validates Agent Skills SKILL.md files.
//
// It implements a deliberately constrained YAML-frontmatter subset (top-level
// scalars, block scalars, and a single nested "metadata" map) and FAILS LOUDLY
// on anything outside that subset rather than guessing. See
// docs/authoring-skills.md for the supported format.
package skillspec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Spec limits (from https://agentskills.io/specification).
const (
	MaxNameLen          = 64
	MaxDescriptionLen   = 1024
	MaxCompatibilityLen = 500
)

var nameRE = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Frontmatter holds the recognized SKILL.md fields.
type Frontmatter struct {
	Name          string
	Description   string
	License       string
	Compatibility string
	AllowedTools  string
	Metadata      map[string]string
}

// Parse splits a SKILL.md into frontmatter and body.
func Parse(content []byte) (Frontmatter, string, error) {
	fm := Frontmatter{Metadata: map[string]string{}}
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	text = strings.TrimPrefix(text, string([]byte{0xEF, 0xBB, 0xBF})) // strip UTF-8 BOM
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], " \t") != "---" {
		return fm, "", fmt.Errorf("missing frontmatter: file must start with a '---' line")
	}
	closing := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], " \t") == "---" {
			closing = i
			break
		}
	}
	if closing < 0 {
		return fm, "", fmt.Errorf("missing frontmatter: no closing '---' line")
	}
	body := strings.TrimPrefix(strings.Join(lines[closing+1:], "\n"), "\n")
	if err := parseFrontmatter(lines[1:closing], &fm); err != nil {
		return fm, "", err
	}
	return fm, body, nil
}

func parseFrontmatter(lines []string, fm *Frontmatter) error {
	i := 0
	for i < len(lines) {
		raw := lines[i]
		if strings.TrimSpace(raw) == "" || strings.HasPrefix(strings.TrimSpace(raw), "#") {
			i++
			continue
		}
		if raw[0] == ' ' || raw[0] == '\t' {
			return fmt.Errorf("unexpected indentation in frontmatter: %q", raw)
		}
		key, val, ok := splitKeyVal(strings.TrimRight(raw, " \t"))
		if !ok {
			return fmt.Errorf("invalid frontmatter line (expected 'key: value'): %q", raw)
		}
		tv := strings.TrimSpace(val)
		if key == "metadata" {
			if tv != "" {
				return fmt.Errorf("'metadata' must be a nested block of 'key: value' pairs, not an inline value")
			}
			ni, err := parseMetadata(lines, i+1, fm)
			if err != nil {
				return err
			}
			i = ni
			continue
		}
		if tv == ">" || tv == ">-" || tv == "|" || tv == "|-" {
			value, ni, err := readBlockScalar(lines, i+1, strings.HasPrefix(tv, ">"))
			if err != nil {
				return err
			}
			assignField(fm, key, value)
			i = ni
			continue
		}
		if strings.HasPrefix(tv, "|") || strings.HasPrefix(tv, ">") {
			return fmt.Errorf("unsupported block scalar indicator %q for key %q", tv, key)
		}
		if strings.HasPrefix(tv, "[") || strings.HasPrefix(tv, "{") || strings.HasPrefix(tv, "&") || strings.HasPrefix(tv, "*") {
			return fmt.Errorf("unsupported YAML construct for key %q: %q", key, tv)
		}
		assignField(fm, key, unquote(tv))
		i++
	}
	return nil
}

func parseMetadata(lines []string, start int, fm *Frontmatter) (int, error) {
	i := start
	for i < len(lines) {
		ml := lines[i]
		if strings.TrimSpace(ml) == "" {
			i++
			continue
		}
		if ml[0] != ' ' && ml[0] != '\t' {
			break // end of nested block
		}
		mkey, mval, ok := splitKeyVal(strings.TrimSpace(ml))
		if !ok {
			return i, fmt.Errorf("invalid metadata line (expected 'key: value'): %q", ml)
		}
		fm.Metadata[mkey] = unquote(strings.TrimSpace(mval))
		i++
	}
	return i, nil
}

// readBlockScalar collects an indented block starting at index start. folded
// (">") joins lines with spaces; literal ("|") joins with newlines.
func readBlockScalar(lines []string, start int, folded bool) (string, int, error) {
	var block []string
	i := start
	for i < len(lines) {
		l := lines[i]
		if strings.TrimSpace(l) == "" {
			block = append(block, "")
			i++
			continue
		}
		if l[0] != ' ' && l[0] != '\t' {
			break
		}
		block = append(block, l)
		i++
	}
	for len(block) > 0 && strings.TrimSpace(block[len(block)-1]) == "" {
		block = block[:len(block)-1]
	}
	minIndent := -1
	for _, l := range block {
		if strings.TrimSpace(l) == "" {
			continue
		}
		ind := len(l) - len(strings.TrimLeft(l, " \t"))
		if minIndent < 0 || ind < minIndent {
			minIndent = ind
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}
	for j, l := range block {
		if len(l) >= minIndent {
			block[j] = l[minIndent:]
		} else {
			block[j] = strings.TrimLeft(l, " \t")
		}
	}
	var out string
	if folded {
		var sb strings.Builder
		prevBlank := false
		for k, l := range block {
			if l == "" {
				sb.WriteString("\n")
				prevBlank = true
				continue
			}
			if k > 0 && !prevBlank {
				sb.WriteString(" ")
			}
			sb.WriteString(l)
			prevBlank = false
		}
		out = sb.String()
	} else {
		out = strings.Join(block, "\n")
	}
	return strings.TrimSpace(out), i, nil
}

func splitKeyVal(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", false
	}
	return key, line[idx+1:], true
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func assignField(fm *Frontmatter, key, value string) {
	switch key {
	case "name":
		fm.Name = value
	case "description":
		fm.Description = value
	case "license":
		fm.License = value
	case "compatibility":
		fm.Compatibility = value
	case "allowed-tools":
		fm.AllowedTools = value
	default:
		// Unknown top-level keys are permitted by the spec; ignored here.
	}
}

// Issue is a single validation problem.
type Issue struct {
	Path    string
	Message string
}

func (i Issue) String() string { return fmt.Sprintf("%s: %s", i.Path, i.Message) }

// ValidateDir validates a single skill directory against the spec.
func ValidateDir(dir string) []Issue {
	var issues []Issue
	add := func(msg string) {
		issues = append(issues, Issue{Path: filepath.Join(dir, "SKILL.md"), Message: msg})
	}
	skillFile := filepath.Join(dir, "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		add("cannot read SKILL.md: " + err.Error())
		return issues
	}
	fm, body, err := Parse(content)
	if err != nil {
		add(err.Error())
		return issues
	}
	folder := filepath.Base(dir)

	if fm.Name == "" {
		add("frontmatter 'name' is required")
	} else {
		if utf8.RuneCountInString(fm.Name) > MaxNameLen {
			add(fmt.Sprintf("'name' exceeds %d characters", MaxNameLen))
		}
		if !nameRE.MatchString(fm.Name) {
			add("'name' must match ^[a-z0-9]+(-[a-z0-9]+)*$")
		}
		if fm.Name != folder {
			add(fmt.Sprintf("'name' (%q) must equal the folder name (%q)", fm.Name, folder))
		}
	}

	if strings.TrimSpace(fm.Description) == "" {
		add("frontmatter 'description' is required")
	} else if n := utf8.RuneCountInString(fm.Description); n > MaxDescriptionLen {
		add(fmt.Sprintf("'description' is %d characters; max is %d", n, MaxDescriptionLen))
	}

	if utf8.RuneCountInString(fm.Compatibility) > MaxCompatibilityLen {
		add(fmt.Sprintf("'compatibility' exceeds %d characters", MaxCompatibilityLen))
	}

	for _, ref := range extractRelativeLinks(body) {
		target := filepath.Join(dir, filepath.FromSlash(ref))
		if _, err := os.Stat(target); err != nil {
			add(fmt.Sprintf("broken reference: %q does not resolve", ref))
		}
	}
	return issues
}

var linkRE = regexp.MustCompile(`\]\(([^)]+)\)`)

// extractRelativeLinks returns relative file paths referenced by markdown links,
// skipping URLs, anchors, mailto and absolute paths.
func extractRelativeLinks(body string) []string {
	var out []string
	for _, m := range linkRE.FindAllStringSubmatch(body, -1) {
		target := strings.TrimSpace(m[1])
		if sp := strings.IndexAny(target, " \t"); sp >= 0 {
			target = target[:sp] // drop optional link title
		}
		if target == "" || strings.Contains(target, "://") ||
			strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") ||
			strings.HasPrefix(target, "/") {
			continue
		}
		if h := strings.IndexAny(target, "#?"); h >= 0 {
			target = target[:h]
		}
		if target != "" {
			out = append(out, target)
		}
	}
	return out
}
