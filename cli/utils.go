package cli

import (
	"bytes"
	"fmt"
	"github.com/weynsee/go-phrase/phrase"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func replacePlaceholders(s string, c *Config, l *phrase.Locale) string {
	translations := map[string]string{
		"<domain>":      c.Domain,
		"<format>":      c.Format,
		"<locale.name>": l.Name,
		"<locale.code>": l.Code,
		"<locale>":      l.Name,
	}
	return replace(s, translations)
}

func replace(s string, translations map[string]string) string {
	if s == "" {
		return s
	}
	for placeholder, t := range translations {
		s = strings.Replace(s, placeholder, t, -1)
	}
	return s
}

func findDefaultLocaleName(c *phrase.Client) (string, error) {
	locales, err := c.Locales.ListAll()
	if err != nil {
		return "", err
	}
	for _, locale := range locales {
		if locale.Default {
			return locale.Name, nil
		}
	}
	return "", nil
}

func isUTF16(b []byte) bool {
	return (b[0] == 0xfe && b[1] == 0xff) || (b[0] == 0xff && b[1] == 0xfe)
}

func decodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}
