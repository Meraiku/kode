package speller

var (
	spellerURL  = "https://speller.yandex.net/services/spellservice.json/checkText"
	spellerOpt  = "6"
	spellerLang = "ru,en"
)

type SpellerResponse struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
}
