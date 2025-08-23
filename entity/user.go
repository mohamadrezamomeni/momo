package entity

type User struct {
	ID               string
	Username         string
	LastName         string
	FirstName        string
	IsAdmin          bool
	Password         string
	IsSuperAdmin     bool
	TelegramID       string
	IsApproved       bool
	TelegramUsername string
	Language         Language
}

type Language int

const (
	EN Language = iota + 1
	FA
)

func LanguageString(language Language) string {
	switch language {
	case EN:
		return "en"
	case FA:
		return "fa"
	}
	return "unkhown"
}

func ConvertLanguageLabelToEnum(languageStr string) Language {
	switch languageStr {
	case "en":
		return EN
	case "fa":
		return FA
	}
	return 0
}
