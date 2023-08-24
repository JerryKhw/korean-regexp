package koreanregexp

import (
	"regexp"
	"slices"
)

func GetPhonemes(char string) (string, string, string, int, int, int) {
	initial := ""
	medial := ""
	finale := ""
	initialOffset := -1
	medialOffset := -1
	finaleOffset := -1

	match, _ := regexp.MatchString(`[ㄱ-ㅎ]`, char)
	if match {
		initial = char
		initialOffset = slices.Index(INITIALS, char)
	} else {
		match, _ = regexp.MatchString(`[가-힣]`, char)
		if match {
			tmp := int([]rune(char)[0]) - BASE

			finaleOffset = tmp % len(FINALES)
			medialOffset = ((tmp - finaleOffset) / len(FINALES)) % len(MEDIALS)
			initialOffset = ((tmp-finaleOffset)/len(FINALES) - medialOffset) / len(MEDIALS)
			initial = INITIALS[initialOffset]
			medial = MEDIALS[medialOffset]
			finale = FINALES[finaleOffset]
		}
	}

	return initial, medial, finale, initialOffset, medialOffset, finaleOffset
}
