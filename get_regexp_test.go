package koreanregexp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegExp(t *testing.T) {
	var data = map[string]string{
		`대한민ㄱ`:  `대한민[ㄱ가-깋]`,
		`대한민구`:  `대한민[구-귛]`,
		`대한민국`:  `대한민(국|구[가-깋])`,
		`온라이`:   `온라[이-잏]`,
		`깎`:     `(깎|까[까-낗]|깍[가-깋])`,
		`뷁`:     `(뷁|뷀[가-깋])`,
		`korea`: `korea`,
	}

	for key, value := range data {
		fmt.Println(`getRegExp `, key, ` → `, value)

		assert.Equal(t, value, GetRegExp(key, GetRegExpOptions{}).String())
	}
}

func TestGetRegExpWithInitialSearch(t *testing.T) {
	fmt.Println(`options.initialSearch`)
	fmt.Println(`initialSearch: false (default)`)

	initialSearch := false
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{}), GetRegExp(`ㅎㄱ`, GetRegExpOptions{InitialSearch: &initialSearch}))
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{InitialSearch: &initialSearch}).String(), `ㅎ[ㄱ가-깋]`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{InitialSearch: &initialSearch}).String(), `\\^ㅎㄱ\\$`)

	fmt.Println(`initialSearch: true`)
	initialSearch = true
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{InitialSearch: &initialSearch}).String(), `[ㅎ하-힣][ㄱ가-깋]`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{InitialSearch: &initialSearch}).String(), `\\^[ㅎ하-힣][ㄱ가-깋]\\$`)
}

func TestGetRegExpWithStartWith(t *testing.T) {
	fmt.Println(`options.startsWith`)
	fmt.Println(`startsWith: false (default)`)

	startsWith := false
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{}), GetRegExp(`ㅎㄱ`, GetRegExpOptions{StartsWith: &startsWith}))
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{StartsWith: &startsWith}).String(), `ㅎ[ㄱ가-깋]`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{StartsWith: &startsWith}).String(), `\\^ㅎㄱ\\$`)

	fmt.Println(`startsWith: true`)
	startsWith = true
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{StartsWith: &startsWith}).String(), `^ㅎ[ㄱ가-깋]`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{StartsWith: &startsWith}).String(), `^\\^ㅎㄱ\\$`)
}

func TestGetRegExpWithEndWith(t *testing.T) {
	fmt.Println(`options.endsWith`)
	fmt.Println(`endsWith: false (default)`)

	endsWith := false
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{}), GetRegExp(`ㅎㄱ`, GetRegExpOptions{EndsWith: &endsWith}))
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{EndsWith: &endsWith}).String(), `ㅎ[ㄱ가-깋]`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{EndsWith: &endsWith}).String(), `\\^ㅎㄱ\\$`)

	fmt.Println(`endsWith: true`)
	endsWith = true
	assert.Equal(t, GetRegExp(`ㅎㄱ`, GetRegExpOptions{EndsWith: &endsWith}).String(), `ㅎ[ㄱ가-깋]$`)
	assert.Equal(t, GetRegExp(`^ㅎㄱ$`, GetRegExpOptions{EndsWith: &endsWith}).String(), `\\^ㅎㄱ\\$$`)
}

func TestGetRegExpWithIgnoreSpace(t *testing.T) {
	fmt.Println(`options.ignoreSpace`)
	fmt.Println(`ignoreSpace: false (default)`)

	ignoreSpace := false
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{}), GetRegExp(`한글날`, GetRegExpOptions{IgnoreSpace: &ignoreSpace}))
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{IgnoreSpace: &ignoreSpace}).String(), `한글(날|나[라-맇])`)

	fmt.Println(`ignoreSpace: true`)
	ignoreSpace = true
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{IgnoreSpace: &ignoreSpace}).String(), `한\\s*글\\s*(날|나[라-맇])`)
}

func TestGetRegExpWithNonCaptureGroup(t *testing.T) {
	fmt.Println(`options.nonCaptureGroup`)
	fmt.Println(`nonCaptureGroup: false (default)`)

	nonCaptureGroup := false
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{}), GetRegExp(`한글날`, GetRegExpOptions{NonCaptureGroup: &nonCaptureGroup}))
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{NonCaptureGroup: &nonCaptureGroup}).String(), `한글(날|나[라-맇])`)

	fmt.Println(`nonCaptureGroup: true`)
	nonCaptureGroup = true
	assert.Equal(t, GetRegExp(`한글날`, GetRegExpOptions{NonCaptureGroup: &nonCaptureGroup}).String(), `한글(?:날|나[라-맇])`)
}

func TestGetRegExpWithFuzzy(t *testing.T) {
	fmt.Println(`options.fuzzy`)
	fmt.Println(`fuzzy: false (default)`)

	initialSearch := true
	fuzzy := false

	pattern := GetRegExp(`ㅋㅍ`, GetRegExpOptions{InitialSearch: &initialSearch, Fuzzy: &fuzzy})
	words := []string{`카페`, `카카오페이`, `카페오레`, `카메라`, `아카펠라`}

	matched := []string{}
	for _, word := range words {
		if pattern.MatchString(word) {
			matched = append(matched, word)
		}
	}
	assert.Equal(t, matched, []string{`카페`, `카페오레`, `아카펠라`})

	fmt.Println(`fuzzy: true`)
	fuzzy = true

	pattern = GetRegExp(`ㅋㅍ`, GetRegExpOptions{InitialSearch: &initialSearch, Fuzzy: &fuzzy})
	words = []string{`카페`, `카카오페이`, `카페오레`, `카메라`, `아카펠라`}

	matched = []string{}
	for _, word := range words {
		if pattern.MatchString(word) {
			matched = append(matched, word)
		}
	}
	assert.Equal(t, matched, []string{`카페`, `카카오페이`, `카페오레`, `아카펠라`})
}
