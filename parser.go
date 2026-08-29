package palidate

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Rule is the structured representation of the validation rules set by the palidate struct tag
type Rule struct {
	required bool
	min      *int
	max      *int
	pattern  *regexp.Regexp
}

// Parser is a parser for palidate tags, not safe for concurent use. If needed, concurent parsing is
// needed create a new parser per-thread.
type Parser struct {
	offset int
	target []byte
}

// Parse a palidate struct tag into a Rule struct, all errors encountered durring parsing are returned
// If no errors exist the slice will be nil so checking for nil is still valid.
func (p *Parser) Parse(rawTag string) (*Rule, []error) {
	// reset the parser before parsing since this same parser can be used for more than one tag
	p.offset = 0
	p.target = []byte(rawTag)

	// make sure we eat any leading whitespace before we try to start parsing
	p.eatWhiteSpace()

	rule := &Rule{}
	var errs []error
	for p.offset < len(p.target) {
		err := p.parseRule(rule)
		if err != nil {
			errs = append(errs, err)
		}

		// parse the trailing comma and any white space so we're ready to parse the next rule at the
		// top of the loop
		err = p.checkComma()
		if err != nil {
			errs = append(errs, err)
			// this is a malformed tag at this point so we just need to bail
			break
		}

		p.eatWhiteSpace()
	}

	if errs != nil {
		return nil, errs
	}

	return rule, nil
}

// eatWhiteSpace will progress the parser past any spaces or tabs. Tags can not contain new lines so
// those are not considered
func (p *Parser) eatWhiteSpace() {
	for p.offset < len(p.target) {
		r, i := utf8.DecodeRune(p.target[p.offset:])
		if r != ' ' && r != '\t' {
			break
		}

		p.offset += i
	}
}

// checkComma checks if the parser is not at the end of the target, ensure there is a comma and progress
// the parser. All whitespace before the comma will be skipped
func (p *Parser) checkComma() error {
	p.eatWhiteSpace()
	if p.offset >= len(p.target) {
		return nil
	}

	if p.target[p.offset] == ',' {
		p.offset++
		return nil
	}

	return errors.New("malformed tag, missing comma after rule")
}

// check if the substring at the current parser position has the given prefix
func (p *Parser) hasPrefix(prefix []byte) bool {
	check := p.target[p.offset:]
	return bytes.HasPrefix(check, prefix)
}

// parseRule parses a single rule from the given tag and sets that field on the given Rule struct
func (p *Parser) parseRule(rule *Rule) error {
	switch {
	case p.hasPrefix([]byte("required")):
		rule.required = true
		p.offset += len("required")
	case p.hasPrefix([]byte("max=")):
		p.offset += len("max=")
		num, err := p.parseInt()
		if err != nil {
			return err
		}

		// need to trigger this after parsing to make sure we don't get stuck mid rule
		if rule.max != nil {
			return errors.New("max is defined more than once")
		}
		rule.max = &num
	case p.hasPrefix([]byte("min=")):
		p.offset += len("min=")
		num, err := p.parseInt()
		if err != nil {
			return err
		}

		// need to trigger this after parsing to make sure we don't get stuck mid rule
		if rule.min != nil {
			return errors.New("min is defined more than once")
		}
		rule.min = &num
	case p.hasPrefix([]byte("regex=")):
		p.offset += len("regex=")
		pat, err := p.parsePattern()
		if err != nil {
			return err
		}

		// need to trigger this after parsing to make sure we don't get stuck mid rule
		if rule.pattern != nil {
			return errors.New("regex is defined more than once")
		}
		rule.pattern = pat
	default:
		s := p.recover()
		return fmt.Errorf("unknown rule '%s'", s)
	}

	return nil
}

// recover eats tokens until it finds a single comma which marks the end of rule. If it reaches the
// end of the target, it stops looking. It returns the substring it found while searching
func (p *Parser) recover() string {
	s := strings.Builder{}
	for p.offset < len(p.target) {
		r, i := utf8.DecodeRune(p.target[p.offset:])
		if r == ',' {
			return s.String()
		}

		s.WriteRune(r)
		p.offset += i
	}

	return s.String()
}

// parseInt parses the int at the current parser location. If no int can be parsed an error is returned
func (p *Parser) parseInt() (int, error) {
	if p.offset >= len(p.target) {
		return 0, errors.New("missing number for validation rule")
	}

	numStr := strings.Builder{}

	// the first character can be a '-' since we support negative numbers
	if p.target[p.offset] == '-' {
		p.offset++
		numStr.WriteRune('-')
	}

	for p.offset < len(p.target) {
		r, i := utf8.DecodeRune(p.target[p.offset:])
		if r < '0' || r > '9' {
			// this is not a valid numeric rune so we should stop parsing
			break
		}

		p.offset += i
		numStr.WriteRune(r)
	}

	i, err := strconv.Atoi(numStr.String())
	if err != nil {
		return 0, err
	}

	return i, nil
}

// parsePattern parses a regex out of the tag. The regex must be enclosed by single quotes. Repeated
// single quotes can be used to escape a single quote literal.
func (p *Parser) parsePattern() (*regexp.Regexp, error) {
	if p.offset >= len(p.target) {
		return nil, errors.New("missing regex for validation rule")
	}

	if p.target[p.offset] != '\'' {
		return nil, errors.New("regex must be enclosed by single quotes")
	}
	p.offset++

	closed := false
	pat := strings.Builder{}
	for p.offset < len(p.target) {
		if p.hasPrefix([]byte("''")) {
			// this is an escaped single quote so we need to add the unescaped value to the pattern
			p.offset += 2
			pat.WriteRune('\'')
			continue
		}

		r, i := utf8.DecodeRune(p.target[p.offset:])
		p.offset += i

		if r == '\'' {
			closed = true
			break
		}

		pat.WriteRune(r)
	}

	if !closed {
		return nil, errors.New("regex must be enclosed by single quotes")
	}

	regex, err := regexp.Compile(pat.String())
	if err != nil {
		return nil, err
	}

	return regex, nil
}
