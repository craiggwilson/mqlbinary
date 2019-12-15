package internal

type Language interface {
	Length() string
}

type JavaLanguage struct{}

func (JavaLanguage) Length() string {
	return "int32"
}
