package translate

type Translator interface {
	Translate(content string) (string, error)
}
