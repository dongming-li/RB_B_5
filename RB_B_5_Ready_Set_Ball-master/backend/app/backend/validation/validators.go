package validation

import (
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"
)

// IsValidString returns true if the string contains only unicode alphabets
func IsValidString(str string) bool {
	for _, c := range str {
		if !unicode.IsLetter(c) {
			return false
		}
	}

	return true
}

//IsValidJoinCode returns true if the string could be a join code
func IsValidJoinCode(str string) bool {
	if len(str) != constraints.JoinCodeLen {
		return false
	}

	for _, c := range str {
		if !isAlpha(c) && !unicode.IsNumber(c) {
			return false
		}
	}

	return true
}

// IsValidGameName returns true if the string meets the requirements for a game name
func IsValidGameName(str string) bool {
	l := utf8.RuneCountInString(str)
	if l > constraints.GameNameMax || l < constraints.GameNameMin {
		return false
	}

	for i, c := range str {
		isInvalidChar := !unicode.IsLetter(c) && c != ' ' && c != '_' && !unicode.IsNumber(c)
		isConsecutiveChar := (i-1 >= 0 && ((str[i-1] == ' ' && c == ' ') || (str[i-1] == '_' && c == '_')))
		//don't allow consecutive spaces, allow underscores but don't allow consecutive
		if isInvalidChar || isConsecutiveChar {
			return false
		}
	}

	return true
}

// IsValidAgeRange returns true if the ages are > [constraints.MinPlayerAge] and < [constraints.MaxPlayerAge]
func IsValidAgeRange(minAge float64, maxAge float64) bool {
	return minAge <= maxAge &&
		minAge >= constraints.MinPlayerAge &&
		maxAge <= constraints.MaxPlayerAge
}

// IsValidLocation returns true if the longitude and latitude floats are valid
func IsValidLocation(lat float64, lng float64) bool {
	if lat == 0 && lng == 0 {
		return false
	}
	if lat < constraints.MinLatitude || lat > constraints.MaxLatitude {
		return false
	}
	if lng < constraints.MinLongitude || lng > constraints.MaxLongitude {
		return false
	}
	return true
}

//IsValidTimeRange returns true if the time range is valid
func IsValidTimeRange(start, end time.Time) bool {
	return start.Before(end)
}

// IsValidName returns true if the [name] meets the name contraints and is a valid string
func IsValidName(name string) bool {
	return utf8.RuneCountInString(name) >= constraints.MinNameLength &&
		utf8.RuneCountInString(name) <= constraints.MaxNameLength && IsValidString(name)
}

// IsValidUserName returns true if the string meets the username requirements
func IsValidUserName(str string) bool {
	l := len(str)
	if l < 1 || str[0] == '_' || unicode.IsNumber(rune(str[0])) { //cannot start with number or underscore
		return false
	}
	for i, c := range str {
		if !unicode.IsLetter(c) &&
			(c != '_' || (i-1 >= 0 && str[i-1] == '_')) && //underscores are allowed but not consecutive
			(!unicode.IsNumber(c)) { //numbers are allowed
			return false
		}
	}

	return l <= constraints.MaxUsernameLength && l >= constraints.MinUsernameLength
}

// IsValidPassword returns true if the password meets the password requirements
func IsValidPassword(passwd []rune) bool {
	l := len(passwd)
	if l < constraints.MinPasswordLength || l > constraints.MaxPasswordength {
		return false
	}

	var upper, lower, digit, punct bool
	for i, c := range passwd {
		if i-1 >= 0 && c == passwd[i-1] { //TODO: this doesn't allow even two consecutive characters. Maybe change that
			return false
		}
		switch true {
		case unicode.IsDigit(c):
			digit = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c):
			punct = true
		}
	}

	return upper && lower && digit && punct
}

// IsValidEmail returns true if the email is valid and false otherwise
// It is generally RFC-5322 compliant except that it doesn't allow '?' in the DNS
func IsValidEmail(email string) bool {
	// basic length check and check by golang standard (is not RFC-5322 compliant)
	if _, err := mail.ParseAddress(email); err != nil || len(email) > constraints.MaxEmailLength {
		return false
	}

	// hypthens don't start or end domain names or end emails
	domain := email[strings.LastIndex(email, "@")+1:]
	if domain[0] == '-' || domain[strings.LastIndex(domain, ".")-1] == '-' || domain[len(domain)-1] == '-' {
		return false
	}

	// only hyphens, digits and periods are allowed in domain names
	for _, c := range domain {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-' && c != '.' {
			return false
		}
	}
	return true
}

// CleanEmail returns a safe email version that can be stored uniquely in the database
func CleanEmail(email string) string {
	at := strings.LastIndex(email, "@")
	return email[:at+1] + strings.ToLower(email[at+1:])
}

// HashPassword returns the hashed version of a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err //NOTE: This line is hard to test because the error would have from bcrypt. See bcrypt Go manual
	}
	return string(hash), nil
}

// CheckPassword returns true if [hashedPassword] is the [bcrypt] hash of [password]
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// IsValidFile returns true if the [typeOfFile] matches the type of file
func IsValidFile(file []byte, typeOfFile string) bool {
	acceptedType, _ := map[string]string{
		"profilepic": "image/png",
	}[typeOfFile]

	return http.DetectContentType(file) == acceptedType
}

// IsValidImageName returns true if the image name has an allowed extension
func IsValidImageName(fileName string) bool {
	name := strings.ToLower(fileName)
	return strings.HasSuffix(name, ".png") ||
		strings.HasSuffix(name, ".jpg") ||
		strings.HasSuffix(name, ".jpeg")
}

// IsValidGameRating returns true if the rating, [r] is between the rating constraints
func IsValidGameRating(r float64) bool {
	return r >= constraints.MinGameRating && r <= constraints.MaxGameRating
}

func isAlpha(c rune) bool {
	return c >= 'A' && c <= 'Z' || c <= 'z' && c >= 'a'
}
