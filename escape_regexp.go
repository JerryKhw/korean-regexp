package koreanregexp

import (
	"regexp"
)

const reRegExpChar = `[\\^$.*+?()[\]{}|]`

var reHasRegExpChar, _ = regexp.Compile(reRegExpChar)

func escapeRegExp(string string) string {
	match := reHasRegExpChar.MatchString(string)

	if match {
		return reHasRegExpChar.ReplaceAllString(string, `\\$0`)
	} else {
		return string
	}
}
