package validation

import (
	"strings"
	"testing"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"
	"golang.org/x/crypto/bcrypt"
)

func TestIsValidString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		//Test cases
		{name: "Normal string", args: args{str: "johnson"}, want: true},
		{name: "string with non-ascii", args: args{str: "汉字"}, want: true},
		{name: "string with numbers", args: args{str: "johnson123"}, want: false},
		{name: "string with special characters", args: args{str: "johnson."}, want: false},
		{name: "string with special characters", args: args{str: "johnson-"}, want: false},
		{name: "string with special characters", args: args{str: "johnson*"}, want: false},
		{name: "string with special characters", args: args{str: "john son"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidString(tt.args.str); got != tt.want {
				t.Errorf("IsValidString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// Test cases.
		{name: "Safe emails are untouched", args: args{"john.doe@smartnotes.com"}, want: "john.doe@smartnotes.com"},
		{name: "Local is untouched", args: args{"j0HN.doe@smartnotes.com"}, want: "j0HN.doe@smartnotes.com"},
		{name: "Domain is lowercased", args: args{"j0HN.doe@sMArTnoTes.com"}, want: "j0HN.doe@smartnotes.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanEmail(tt.args.email); got != tt.want {
				t.Errorf("CleanEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		//Test cases for local part of the email
		{name: "Normal emails are valid", args: args{"john.doe@smartnotes.com"}, want: true},
		{name: "Empty emails are invalid", args: args{""}, want: false},
		{name: "Emails containing two @ runes are invalid", args: args{"john.doe@smart@notes.com"}, want: false},
		{name: "Comments at the beginning are not allowed", args: args{"(jo)hn.doe@smart@notes.com"}, want: false},
		{name: "Comments at the end are not allowed", args: args{"john.d(oe)@smart@notes.com"}, want: false},
		{name: "Consequtive dots are not allowed without quotes", args: args{"john..doe@smartnotes.com"}, want: false},
		{name: "Consequtive dots are allowed with quotes", args: args{`"john..doe"@smartnotes.com`}, want: true},
		{name: "Dots cannot start the local part", args: args{".john.doe@smartnotes.com"}, want: false},
		{name: "Dots cannot end the local part", args: args{".john.doe.@smartnotes.com"}, want: false},
		{name: "!#$%&'*+-/=?^_`{|}~ are allowed without quotes", args: args{"john!#$%&'*+-/=?^_`{|}~@smartnotes.com"}, want: true},
		{name: "\"(),:;<>@[\\] are allowed with quotes", args: args{`"john.doe\"(),:;<>@[\\]"@smartnotes.com`}, want: true},
		{name: "\"(),:;<>@[\\] are not allowed without quotes", args: args{`john.doe"(),:;<>@[\\]@smartnotes.com`}, want: false},
		// Test cases for domain part of the email
		{name: "hyphens don't start domain names", args: args{"john.doe@-smartnotes.com"}, want: false},
		{name: "hyphens don't end domain names", args: args{"john.doe@smartnotes.com-"}, want: false},
		{name: "hyphens don't come just before the period", args: args{"john.doe@smartnotes-.com"}, want: false},
		{name: "question marks not allowed besides hyphens and dots", args: args{"john.doe@smart?notes.com"}, want: false},
		{name: "special chars not allowed besides hyphens and dots", args: args{"john.doe@smart$%notes.com"}, want: false},
		{name: "Consequtive dots are not allowed", args: args{"john.doe@smartnotes..com"}, want: false},
		{name: "Digits are allowed", args: args{"john.doe@9smartnotes.com"}, want: true},
		{name: "hyphens are allowed", args: args{"john.doe@smart-notes.com"}, want: true},
		{name: "spaces are not allowed", args: args{"john.doe@smart notes.com"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.args.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidUserName(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		//Test cases
		{name: "Normal username", args: args{str: "johnson"}, want: true},
		{name: "Normal username containing underscores", args: args{str: "john_son"}, want: true},
		{name: "username with non-ascii", args: args{str: "汉字"}, want: true},
		{name: "username with numbers", args: args{str: "johnson123"}, want: true},
		{name: "Normal username containing consecutive underscores", args: args{str: "john__son"}, want: false},
		{name: "username starting with numbers", args: args{str: "1johnson123"}, want: false},
		{name: "username with special characters", args: args{str: "johnson."}, want: false},
		{name: "username with special characters", args: args{str: "johnson-"}, want: false},
		{name: "username with special characters", args: args{str: "johnson*"}, want: false},
		{name: "username with special characters", args: args{str: "john son"}, want: false},
		{name: "username starting with an underscore", args: args{str: "_johnson"}, want: false},
		{name: "username at max length", args: args{str: strings.Repeat("i", constraints.MaxUsernameLength)}, want: true},
		{name: "username beyond max length", args: args{str: strings.Repeat("i", constraints.MaxUsernameLength+1)}, want: false},
		{name: "username at min length", args: args{str: strings.Repeat("i", constraints.MinUsernameLength)}, want: true},
		{name: "username below min length", args: args{str: strings.Repeat("i", constraints.MinUsernameLength-1)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUserName(tt.args.str); got != tt.want {
				t.Errorf("IsValidUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	type args struct {
		passwd []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		//Test cases
		{name: "Containing all parameters", args: args{passwd: []rune("Ab.4567890")}, want: true},
		{name: "Just alphabets", args: args{passwd: []rune("johnson")}, want: false},
		{name: "No digit", args: args{passwd: []rune("Ab.defghij")}, want: false},
		{name: "No punctuation", args: args{passwd: []rune("Ab34567890")}, want: false},
		{name: "No lower rune", args: args{passwd: []rune("AB.4567890")}, want: false},
		{name: "No uppercase rune", args: args{passwd: []rune("ab.4567890")}, want: false},
		{name: "No consecutive characters", args: args{passwd: []rune("Abb.4567890")}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.args.passwd); got != tt.want {
				t.Errorf("IsValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Test cases.
		{
			name: "valid name",
			args: args{"Walt"},
			want: true,
		},
		{
			name: "name less than minimum allowed",
			args: args{strings.Repeat("a", constraints.MinNameLength-1)},
			want: false,
		},
		{
			name: "name greater than maximum allowed",
			args: args{strings.Repeat("a", constraints.MaxNameLength+1)},
			want: false,
		},
		{
			name: "name is invalid string",
			args: args{"_invalid"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidName(tt.args.name); got != tt.want {
				t.Errorf("IsValidName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	rawPass := "Don't try this at home"
	hash, err := HashPassword(rawPass)
	if err != nil {
		t.Errorf("Unexpected error after HashPassword: %v", err)
	}

	isCorrect := CheckPassword(hash, rawPass)
	if !isCorrect {
		t.Errorf("checkPassword should return true on same passwords")
	}

	isCorrect = CheckPassword(hash, rawPass+" ")
	if isCorrect {
		t.Errorf("checkPassword should return false on different passwords")
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Errorf("Unexpectd error when getting cost from bcrypt: %v", err)
	}
	if cost != 10 {
		t.Errorf("bcrypt cost should be default but cost was %d", cost)
	}

	//TODO test hashpassword returning an error
	// hash, err = HashPassword("")
	// if err != nil {
	// 	t.Errorf("Unexpected error after HashPassword: %v", err)
	// }
}

func TestIsValidAgeRange(t *testing.T) {
	type args struct {
		minAge float64
		maxAge float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid age range",
			args: args{minAge: 18, maxAge: 60},
			want: true,
		},
		{
			name: "Minimum age > maximum age",
			args: args{minAge: 20, maxAge: 7},
			want: false,
		},
		{
			name: "Age below constraint",
			args: args{minAge: constraints.MinPlayerAge - 1, maxAge: 25},
			want: false,
		},
		{
			name: "Age above contraint",
			args: args{minAge: 20, maxAge: constraints.MaxPlayerAge + 1},
			want: false,
		},
		{
			name: "Same age",
			args: args{minAge: 25, maxAge: 25},
			want: true,
		},
		{
			name: "Same as contraints",
			args: args{minAge: constraints.MinPlayerAge, maxAge: constraints.MaxPlayerAge},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAgeRange(tt.args.minAge, tt.args.maxAge); got != tt.want {
				t.Errorf("IsValidAgeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidLocation(t *testing.T) {
	type args struct {
		lat float64
		lng float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid location",
			args: args{lat: 0, lng: 50},
			want: true,
		},
		{
			name: "Latitude and Longitude are zero",
			args: args{lat: 0, lng: 0},
			want: false,
		},
		{
			name: "Latitude too low",
			args: args{lat: constraints.MinLatitude - 1, lng: 50},
			want: false,
		},
		{
			name: "Latitude too high",
			args: args{lat: constraints.MaxLatitude + 1, lng: 50},
			want: false,
		},
		{
			name: "Longitude too low",
			args: args{lat: 50, lng: constraints.MinLongitude - 1},
			want: false,
		},
		{
			name: "Longitude too high",
			args: args{lat: 50, lng: constraints.MaxLongitude + 1},
			want: false,
		},
		{
			name: "Coordinates on lower constraint",
			args: args{lat: constraints.MinLatitude, lng: constraints.MinLongitude},
			want: true,
		},
		{
			name: "Coordinates on higer contriant",
			args: args{lat: constraints.MaxLatitude, lng: constraints.MaxLongitude},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidLocation(tt.args.lat, tt.args.lng); got != tt.want {
				t.Errorf("IsValidLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidTimeRange(t *testing.T) {
	smallTime := time.Now()
	bigTime := smallTime.Add(time.Hour)
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Duration",
			args: args{start: smallTime, end: bigTime},
			want: true,
		},
		{
			name: "Invalid Duration",
			args: args{start: bigTime, end: smallTime},
			want: false,
		},
		{
			name: "Start time equal to end time",
			args: args{start: smallTime, end: smallTime},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTimeRange(tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("IsValidTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidJoinCode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Valid join code with mixed charcters", args: args{str: "ABCDEF7"}, want: true},
		{name: "Valid join code with only numbers", args: args{str: "1234567"}, want: true},
		{name: "Valid join code with only letters", args: args{str: "ABCDEFG"}, want: true},
		{name: "Invalid join code with special charcacters (.)", args: args{str: "ABCDEF."}, want: false},
		{name: "Invalid join code with special characters (-)", args: args{str: "ABCDEF-"}, want: false},
		{name: "Invalid join code with special characters (*)", args: args{str: "ABCDEF*"}, want: false},
		{name: "Invalid join code with special characters ( )", args: args{str: "ABCD EF"}, want: false},
		{name: "Invalid join code with underscore", args: args{str: "ABCDEF_"}, want: false},
		{name: "Invalid join code with non-ascii aplhabets", args: args{str: "汉字EFG"}, want: false},
		{name: "Invalid join code with too many charcaters", args: args{str: strings.Repeat("a", constraints.JoinCodeLen-1)}, want: false},
		{name: "Invalid join code with too few charcaters", args: args{str: strings.Repeat("a", constraints.JoinCodeLen+1)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidJoinCode(tt.args.str); got != tt.want {
				t.Errorf("IsValidJoinCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidGameName(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Valid with mixed charcters", args: args{str: "Superstars7"}, want: true},
		{name: "Valid with only numbers", args: args{str: "1234567"}, want: true},
		{name: "Valid with spaces", args: args{str: "Supers tars7"}, want: true},
		{name: "Valid with only letters", args: args{str: "ABCDEFG"}, want: true},
		{name: "Valid with non-ascii aplhabets", args: args{str: "汉字EFG"}, want: true},
		{name: "Invalid with special charcacters (.)", args: args{str: "ABCDEF."}, want: false},
		{name: "Invalid with special characters (-)", args: args{str: "ABCDEF-"}, want: false},
		{name: "Invalid with special characters (*)", args: args{str: "ABCDEF*"}, want: false},
		{name: "Invalid with two consecutive spaces", args: args{str: "AB  EFGHI"}, want: false},
		{name: "Invalid with first two consecutive spaces", args: args{str: "  ABEFGHI"}, want: false},
		{name: "Invalid with last two consecutive spaces", args: args{str: "ABEFGHI  "}, want: false},
		{name: "Valid with two nonconsecutive spaces", args: args{str: "A B EFGHI"}, want: true},

		{name: "Invalid with two consecutive underscores", args: args{str: "AB__EFGHI"}, want: false},
		{name: "Invalid with first two consecutive underscores", args: args{str: "__ABEFGHI"}, want: false},
		{name: "Invalid with last two consecutive underscores", args: args{str: "ABEFGHI__"}, want: false},
		{name: "Valid with two nonconsecutive underscores", args: args{str: "A_B_EFGHI"}, want: true},

		{name: "Invalid join code with too many charcaters", args: args{str: strings.Repeat("a", constraints.GameNameMin-1)}, want: false},
		{name: "Invalid join code with too few charcaters", args: args{str: strings.Repeat("a", constraints.GameNameMax+1)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidGameName(tt.args.str); got != tt.want {
				t.Errorf("IsValidGameName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidGameRating(t *testing.T) {
	type args struct {
		r float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Minimim game rating is valid", args: args{constraints.MinGameRating}, want: true},
		{name: "Maximum game rating is valid", args: args{constraints.MaxGameRating}, want: true},
		{name: "float rating is valid", args: args{3.5}, want: true},
		{name: "integer rating is valid", args: args{4}, want: true},
		{name: "Below minimim game rating is invalid", args: args{constraints.MinGameRating - 0.1}, want: false},
		{name: "Above maximum game rating is invalid", args: args{constraints.MaxGameRating + 0.1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidGameRating(tt.args.r); got != tt.want {
				t.Errorf("IsValidGameRating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidImageName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Valid jpg", args: args{"pic.jpg"}, want: true},
		{name: "Valid png", args: args{"pic.png"}, want: true},
		{name: "Valid jpeg", args: args{"pic.jpeg"}, want: true},
		{name: "Valid JPG", args: args{"pic.JPG"}, want: true},
		{name: "Valid PNG", args: args{"pic.PNG"}, want: true},
		{name: "Valid JPEF", args: args{"pic.JPEG"}, want: true},
		{name: "File with two extensions", args: args{"pic.png.jpg"}, want: true},
		{name: "File with no extension", args: args{"pic"}, want: false},
		{name: "File with invalid extension", args: args{"pic.txt"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidImageName(tt.args.fileName); got != tt.want {
				t.Errorf("IsValidImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
