package utility

import (
	"bytes"
	"html/template"
	"regexp"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dongri/phonenumber"
)

func NormalizePhoneNumber(phoneNum string, countryCode string) string {
	code := countryCode

	number := phonenumber.Parse(phoneNum, code)
	if number != "" {
		return number
	}

	// default MY
	number = phonenumber.Parse(phoneNum, "MY")
	if number != "" {
		return number
	}

	country := phonenumber.GetISO3166ByNumber(phoneNum, true)
	number = phonenumber.ParseWithLandLine(phoneNum, country.CountryName)
	if number != "" {
		return number
	}

	country = phonenumber.GetISO3166ByNumber(phoneNum, false)
	number = phonenumber.ParseWithLandLine(phoneNum, country.CountryName)
	if number != "" {
		return number
	}

	country = phonenumber.GetISO3166ByNumber("+"+phoneNum, true)
	number = phonenumber.ParseWithLandLine("+"+phoneNum, country.CountryName)
	if number != "" {
		return number
	}

	country = phonenumber.GetISO3166ByNumber("+"+phoneNum, false)
	number = phonenumber.ParseWithLandLine("+"+phoneNum, country.CountryName)
	if number != "" {
		return number
	}

	return number
}

func NormalizeID(id string) string {
	re := regexp.MustCompile(`[^0-9a-zA-Z]`)
	return strings.ToUpper(re.ReplaceAllString(id, ""))
}

func NormalizeDate(date string) (string, error) {
	re := regexp.MustCompile(`[^0-9]`)
	d := strings.ToUpper(re.ReplaceAllString(date, ""))
	if len(d) < 8 {
		return "", constants.InvalidDateError
	}
	return d, nil
}

func NormalizeRole(role string) string {
	re := regexp.MustCompile(`[^a-zA-Z]`)
	return strings.ToLower(re.ReplaceAllString(role, ""))
}

func NormalizeName(name string) string {
	return strings.Trim(name, " ")
}

func ValidateEmail(email string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(email)
}

// MalaysiaTime gets Malaysia time
func MalaysiaTime(t time.Time) time.Time {
	// Load required location
	location, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		return t
	}

	return t.In(location)
}

// TimeToMilli converts time to millisecond
func TimeToMilli(time time.Time) int64 {
	return MalaysiaTime(time).UnixNano() / 1000000
}

// MilliToTime converts millisecond to time
func MilliToTime(milli int64) time.Time {
	return MalaysiaTime(time.Unix(0, milli*int64(time.Millisecond)))
}

// DateStringToTime converts date string to time
func DateStringToTime(date string) (time.Time, error) {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return time.Now(), err
	}
	t = t.Add(-8 * time.Hour)

	return MalaysiaTime(t), nil
}

// TimeToDateString timestamp to date string (yyyyMMdd)
func TimeToDateString(t time.Time) string {
	return MalaysiaTime(t).Format("20060102")
}

// RemoveZeroWidth removes zero width characters from string
func RemoveZeroWidth(t string) string {
	rslt := strings.Map(func(r rune) rune {
		if r == '↵' || r == '\n' || unicode.IsGraphic(r) &&
			r != '\u2000' &&
			r != '\u2001' &&
			r != '\u2002' &&
			r != '\u2003' &&
			r != '\u2004' &&
			r != '\u2005' &&
			r != '\u2006' &&
			r != '\u2007' &&
			r != '\u2008' &&
			r != '\u2009' &&
			r != '\u200a' &&
			r != '\u202f' &&
			r != '\u205f' &&
			r != '⠀' &&
			r != '\u3000' {
			return r
		}
		return -1
	}, t)

	// for weird characters like zalgo
	if utf8.RuneCountInString(rslt) > 500 {
		return ""
	}

	rslt = strings.Trim(rslt, " ")

	return rslt
}

func UserInUsers(users []*dto.User, user *dto.User) bool {
	for _, u := range users {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

func ZoneInZones(zones []*dto.Zone, zone *dto.Zone) bool {
	for _, z := range zones {
		if z.ID == zone.ID {
			return true
		}
	}
	return false
}

func ParseHTMLTemplate(templateFilename string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetZoneRisk(radius float64, capacity, usersWithin int64) int64 {
	score := float64(usersWithin) / radius / float64(capacity) * 100
	// score = usersWithin / radius / capacity * 100%
	if score > 100 {
		return constants.MaximumRisk
	}
	if score > 80 {
		return constants.HighRisk
	}
	if score > 60 {
		return constants.MediumRisk
	}
	if score > 40 {
		return constants.LowRisk
	}
	return constants.MinimumRisk
}
