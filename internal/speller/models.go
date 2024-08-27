package speller

var (
	spellerURL  = "http://speller.yandex.net/services/spellservice.json/checkText"
	spellerOpt  = "6"
	spellerLang = "ru,en"
)

type SpellerRequest struct {
	Text    []string `json:"text"`
	Options string   `json:"options"`
}

type SpellerResponse struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
}
