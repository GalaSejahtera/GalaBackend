package utility

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestToDateString ...
func TestTimeToDateString(t *testing.T) {
	// Load required location
	location, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		t.Error("failed to load location.")
	}

	tests := []struct {
		name           string
		inputTime      time.Time
		expectedResult string
	}{
		{
			name:           "12 am, should return same date",
			inputTime:      time.Date(1996, 10, 10, 0, 0, 0, 0, location),
			expectedResult: "19961010",
		},
		{
			name:           "normal time, should return same date",
			inputTime:      time.Date(2020, 3, 31, 10, 0, 0, 0, location),
			expectedResult: "20200331",
		},
		{
			name:           "11.36, should return same date",
			inputTime:      time.Date(2020, 3, 31, 23, 59, 59, 0, location),
			expectedResult: "20200331",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, TimeToDateString(test.inputTime))
	}
}

// DateStringToTime ...
func TestDateStringToTime(t *testing.T) {
	// Load required location
	location, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		t.Error("failed to load location.")
	}

	tests := []struct {
		name           string
		date           string
		expectedResult time.Time
		expectedErr    bool
	}{
		{
			name:           "valid date, should return Malaysia time 12 am",
			date:           "20200331",
			expectedResult: time.Date(2020, 3, 31, 0, 0, 0, 0, location),
		},
		{
			name:        "invalid date, should return err",
			date:        "20205555",
			expectedErr: true,
		},
	}

	for _, test := range tests {
		actualResult, err := DateStringToTime(test.date)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, test.expectedResult, actualResult)
		}
	}
}

// TestNormalizeEmail ...
func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		expectedResult string
	}{
		{
			name:           "email with spaces behind, should be trimmed",
			date:           "asda@jiaosd.jsaioda   ",
			expectedResult: "asda@jiaosd.jsaioda",
		},
		{
			name:           "email with spaces in between, should be trimmed",
			date:           "asda@jiaosd   .jsaioda   ",
			expectedResult: "asda@jiaosd.jsaioda",
		},
		{
			name:           "email with upper cases, should be lowered",
			date:           "JJ@jiaosd   .jsaioda   ",
			expectedResult: "jj@jiaosd.jsaioda",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, NormalizeEmail(test.date))
	}
}

// TestNormalizePhoneNumber ...
func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumber    string
		countryCode    string
		expectedResult string
	}{
		{
			name:           "valid malaysia phone number",
			phoneNumber:    "60127153013",
			countryCode:    "",
			expectedResult: "60127153013",
		},
		{
			name:           "valid malaysia phone number without country code",
			phoneNumber:    "0127153013",
			countryCode:    "",
			expectedResult: "60127153013",
		},
		{
			name:           "valid malaysia phone number without 60",
			phoneNumber:    "127153013",
			countryCode:    "",
			expectedResult: "60127153013",
		},
		{
			name:           "valid malaysia phone number with spaces and dash",
			phoneNumber:    "012-715 3013  ",
			countryCode:    "",
			expectedResult: "60127153013",
		},
		{
			name:           "valid singapore phone number",
			phoneNumber:    "+6512342134",
			countryCode:    "",
			expectedResult: "6512342134",
		},
		{
			name:           "valid Korea phone number",
			phoneNumber:    "821072814030",
			countryCode:    "",
			expectedResult: "821072814030",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, NormalizePhoneNumber(test.phoneNumber, test.countryCode))
	}
}

// TestNormalizeName ...
func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name           string
		n              string
		expectedResult string
	}{
		{
			name:           "lower case name",
			n:              "jia xiong",
			expectedResult: "jia xiong",
		},
		{
			name:           "lower case name with spaces",
			n:              "   jia xiong    ",
			expectedResult: "jia xiong",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, NormalizeName(test.n))
	}
}

func TestRemoveZeroWidth(t *testing.T) {

	tests := []struct {
		name           string
		text           string
		expectedResult string
	}{
		{
			name:           "normal text, should return same text",
			text:           "jiaxiong",
			expectedResult: "jiaxiong",
		},
		{
			name:           "with zero width, should return without zero width",
			text:           "a\u200B\u200C\u200D\uFEFFb",
			expectedResult: "ab",
		},
		{
			name:           "Chinese, should return chinese",
			text:           "嘉雄",
			expectedResult: "嘉雄",
		},
		{
			name:           "Japanese, should return japanese",
			text:           "なに",
			expectedResult: "なに",
		},
		{
			name:           "Tamil character, should return tamil",
			text:           "சுதந்திரமாகவே பிறக்கின்",
			expectedResult: "சுதந்திரமாகவே பிறக்கின்",
		},
		{
			name:           "U+2000 En Quad",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2001 Em Quad",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+00A0 No-Break Space",
			text:           " ",
			expectedResult: "",
		},
		{
			name:           "U+2002 En Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2003 Em Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2004 Three-Per-Em Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2005 Four-Per-Em Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2006 Six-Per-Em Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2007 Figure Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2008 Punctuation Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2009 Thin Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+200A Hair Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+205F Medium Mathematical Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+202F Narrow No-Break Space",
			text:           "< >",
			expectedResult: "<>",
		},
		{
			name:           "U+2800 Braille Pattern Blank",
			text:           "<⠀>",
			expectedResult: "<>",
		},
		{
			name:           "U+3000 Ideographic Space",
			text:           "<　>",
			expectedResult: "<>",
		},
		{
			name:           "multiline text",
			text:           "<↵>",
			expectedResult: "<↵>",
		},
		{
			name:           "nextline text",
			text:           "<\n>",
			expectedResult: "<\n>",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedResult, RemoveZeroWidth(test.text), test.name)
	}
}
