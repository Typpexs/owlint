package utilities

import (
	"github.com/bregydoc/gtranslate"
)

func TranslateText(from string, to string, text string) (string, error) {
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: from,
			To:   to,
		},
	)
	if err != nil {
		return "", err
	}
	return translated, nil
}
