package koreanregexp

import (
	"math/big"
	"regexp"
	"slices"
	"strings"
)

func getInitialSearchRegExp(initial string, allowOnlyInitial bool) string {
	initialOffset := slices.Index(INITIALS, initial)

	if initialOffset != -1 {
		baseCode := initialOffset*len(MEDIALS)*len(FINALES) + BASE

		result := "["
		if allowOnlyInitial {
			result += initial
		}
		result += string(rune(baseCode)) + "-" + string(rune(baseCode+len(MEDIALS)*len(FINALES)-1))

		return result + "]"
	}

	return initial
}

type GetRegExpOptions struct {
	InitialSearch   *bool
	StartsWith      *bool
	EndsWith        *bool
	IgnoreSpace     *bool
	Fuzzy           *bool
	NonCaptureGroup *bool
}

var fuzzyInt, _ = new(big.Int).SetString("fuzzy", 36)
var FUZZY = "__" + fuzzyInt.String() + "__"

var ignoreSpaceInt, _ = new(big.Int).SetString("ignorespace", 36)
var IGNORE_SPACE = "__" + ignoreSpaceInt.String() + "__"

func GetRegExp(search string, options GetRegExpOptions) *regexp.Regexp {
	frontChars := strings.Split(search, "")
	lastChar := frontChars[len(frontChars)-1]
	lastCharPattern := ""

	initial, medial, finale, initialOffset, medialOffset, _ := GetPhonemes(lastChar)

	if initialOffset != -1 {
		frontChars = frontChars[:len(frontChars)-1]

		baseCode := initialOffset*len(MEDIALS)*len(FINALES) + BASE

		patterns := []string{}

		for {
			// case 1: 종성으로 끝나는 경우 (받침이 있는 경우)
			if finale != "" {
				// 마지막 글자
				patterns = append(patterns, lastChar)
				// 종성이 초성으로 사용 가능한 경우
				if slices.Contains(INITIALS, finale) {
					patterns = append(patterns, string(rune((baseCode+medialOffset*len(FINALES))))+getInitialSearchRegExp(finale, false))
				}
				// 종성이 복합 자음인 경우, 두 개의 자음으로 분리하여 각각 받침과 초성으로 사용
				if MIXED[finale] != nil {
					patterns = append(patterns, string(rune(baseCode+medialOffset*len(FINALES)+slices.Index(FINALES, MIXED[finale][0])))+getInitialSearchRegExp(MIXED[finale][1], false))
				}
				break
			}

			// case 2: 중성으로 끝나는 경우 (받침이 없는 경우)
			if medial != "" {
				var from int
				var to int

				// 중성이 복합 모음인 경우 범위를 확장하여 적용
				if MEDIAL_RANGE[medial] != nil {
					from = baseCode + slices.Index(MEDIALS, MEDIAL_RANGE[medial][0])*len(FINALES)
					to = baseCode + slices.Index(MEDIALS, MEDIAL_RANGE[medial][1])*len(FINALES) + len(FINALES) - 1
				} else {
					from = baseCode + medialOffset*len(FINALES)
					to = from + len(FINALES) - 1
				}

				patterns = append(patterns, "["+string(rune(from))+"-"+string(rune(to))+"]")
				break
			}

			// case 3: 초성만 입력된 경우
			if initial != "" {
				patterns = append(patterns, getInitialSearchRegExp(initial, true))
				break
			}
		}

		if len(patterns) > 1 {
			if options.NonCaptureGroup != nil && *options.NonCaptureGroup {
				lastCharPattern = "(?:" + strings.Join(patterns, "|") + ")"
			} else {
				lastCharPattern = "(" + strings.Join(patterns, "|") + ")"
			}
		} else {
			lastCharPattern = patterns[0]
		}
	}

	glue := ""

	if options.Fuzzy != nil && *options.Fuzzy {
		glue = FUZZY
	} else if options.IgnoreSpace != nil && *options.IgnoreSpace {
		glue = IGNORE_SPACE
	}

	frontCharsPattern := ""

	if options.InitialSearch != nil && *options.InitialSearch {
		tmpFrontChars := []string{}
		for _, char := range frontChars {
			var reg, _ = regexp.Compile(`[ㄱ-ㅎ]`)

			if reg.FindStringIndex(char) != nil {
				tmpFrontChars = append(tmpFrontChars, getInitialSearchRegExp(char, true))

			} else {
				tmpFrontChars = append(tmpFrontChars, escapeRegExp(char))
			}
		}

		frontCharsPattern = strings.Join(tmpFrontChars, glue)
	} else {
		frontCharsPattern = escapeRegExp(strings.Join(frontChars, glue))
	}

	pattern := ""

	if options.StartsWith != nil && *options.StartsWith {
		pattern += "^"
	}
	pattern += frontCharsPattern + glue + lastCharPattern

	if options.EndsWith != nil && *options.EndsWith {
		pattern += "$"
	}

	if glue != "" {
		fuzzyReg, _ := regexp.Compile(FUZZY)

		pattern = fuzzyReg.ReplaceAllString(pattern, `.*`)

		ignoreSpaceReg, _ := regexp.Compile(IGNORE_SPACE)

		pattern = ignoreSpaceReg.ReplaceAllString(pattern, `\\s*`)
	}

	return regexp.MustCompile(pattern)
}
