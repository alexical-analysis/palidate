package palidate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TODO: need a better name than parsedTag
type parsedTag struct {
	required bool
	min      *int
	max      *int
	pattern  *regexp.Regexp
}

// TODO: this probably needs to be a struct so it can hold some state. Right now too much of the logic
// is just held inside the function
func parseTag(rawTag string) parsedTag {
	if rawTag == "" {
		return parsedTag{}
	}

	parsed := parsedTag{}
	for offset := 0; offset < len(rawTag); {
		frag := rawTag[offset:]

		fieldName := ""
		switch {
		case frag[0] == ' ' || frag[0] == '\t':
			// just skip whitespace
			offset++
			continue
		case strings.HasPrefix(frag, "required"):
			fieldName = "required"
			parsed.required = true
			offset += 8
		case strings.HasPrefix(frag, "max"):
			fieldName = "max"

			if frag[3] != '=' {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, missing '=' after 'max' field")
			}

			maxInt, count, err := parseInt(frag[4:])
			if err != nil {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, 'max' value must be a valid int")
			}
			offset += count + 4 // 4 for 'max' + '='
			parsed.max = &maxInt

		case strings.HasPrefix(frag, "min"):
			fieldName = "min"

			if frag[3] != '=' {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, missing '=' after 'min' field")
			}

			minInt, count, err := parseInt(frag[4:])
			if err != nil {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, 'min' value must be a valid int")
			}
			offset += count + 4 // 4 for min and '='
			parsed.min = &minInt

		case strings.HasPrefix(frag, "regex"):
			fieldName = "regex"
			if frag[5] != '=' {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, missing '=' after 'min' field")
			}

			if frag[6] != '\'' {
				// malformed field, panic so it's easy to track down
				panic("invalid tag, regex patterns must be contained inside single quotes")
			}

			regex, count, err := parsePattern(frag[7:])
			if err != nil {
				// malformed field, panic so it's easy to track down
				pattern := frag[7 : count+7]
				panic(fmt.Sprintf("invalid regex, failed to compile '%s' into a valid pattern", pattern))
			}
			offset += count + 7 // 7 for 'regex' + '=' + leading single quote
			parsed.pattern = regex
		default:
			// malformed field, panic so it's easy to track down
			panic(fmt.Sprintf("unknown field '%s'", frag))
		}

		// TODO: this eats any trailing white space but I know there's a better way to do this
		for {
			if offset >= len(rawTag) {
				break
			}

			if rawTag[offset] == ' ' || rawTag[offset] == '\t' {
				offset++
			} else {
				break
			}
		}

		if offset >= len(rawTag) {
			break
		}

		if rawTag[offset] != ',' {
			panic(fmt.Sprintf("invalid tag, missing ',' after '%s' field", fieldName))
		}
		// skip the , so we can start parsing the next field
		offset++
	}

	return parsed
}

// parseInt parses an int from a tag fragment and returns that int and the number of characters that
// were parsed to create it.
func parseInt(frag string) (int, int, error) {
	count := 0
	for i, r := range frag {
		count++
		if r >= '0' && r <= '9' {
			continue
		}

		// the first character can be negative
		if i == 0 && r == '-' {
			continue
		}

		// take one back since this rune was not actually a match
		count--
		break
	}

	i, err := strconv.Atoi(frag[:count])
	if err != nil {
		// need to make sure we return count here so we don't get stuck in a parsing loop
		return 0, count, err
	}

	return i, count, nil
}

// parsePattern parses a regex from a tag fragment and returns that regex and the number of characters
// that were parsed to create it. Compiling the regexp can fail so the error is returned if it does.
func parsePattern(frag string) (*regexp.Regexp, int, error) {
	count := 0
	escape := false
	for _, r := range frag {
		count++
		if r == '\'' && escape {
			escape = false
			continue
		}
		if escape {
			count--
			break
		}
		if r == '\'' {
			escape = true
		}
	}

	// TODO: we should actually build the string using a string.Builder instead of doing this hacky fix
	pattern := frag[0 : count-1]
	pattern = strings.ReplaceAll(pattern, "''", "'")
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, count, err
	}

	return regex, count, nil
}
