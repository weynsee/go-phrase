package cli

import (
	"fmt"
	"github.com/weynsee/go-phrase/phrase"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type formatProperties struct {
	localeAware     bool
	extensions      []string
	directoryFormat string
	filenameFormat  string
	targetDirectory string
	localeExtension bool
}

type format interface {
	properties() *formatProperties
	directoryForLocale(*Config, *phrase.Locale) string
	filenameForLocale(*Config, *phrase.Locale) string
	extractLocaleFromPath(*phrase.Client, string) (string, error)
}

type defaultFormat struct {
	*formatProperties
}

func (f *defaultFormat) properties() *formatProperties {
	return f.formatProperties
}

func (f *defaultFormat) directoryForLocale(c *Config, l *phrase.Locale) string {
	return replacePlaceholders(f.directoryFormat, c, l)
}

func (f *defaultFormat) filenameForLocale(c *Config, l *phrase.Locale) string {
	return replacePlaceholders(f.filenameFormat, c, l)
}

func (f *defaultFormat) extractLocaleFromPath(c *phrase.Client, path string) (string, error) {
	return "", nil
}

type xmlFormat struct {
	*formatProperties
}

func (f *xmlFormat) properties() *formatProperties {
	return f.formatProperties
}

func (f *xmlFormat) directoryForLocale(c *Config, l *phrase.Locale) string {
	if l.Default {
		return "values"
	}
	name := l.Code
	if name == "" {
		name = l.Name
	}
	return fmt.Sprintf("values-%s", f.formatName(name))
}

func (f *xmlFormat) formatName(s string) string {
	if strings.Index(s, "-") == -1 {
		return s
	}
	parts := strings.Split(s, "-")
	return fmt.Sprintf("%s-r%s", parts[0], strings.ToUpper(parts[len(parts)-1]))
}

func (f *xmlFormat) filenameForLocale(c *Config, l *phrase.Locale) string {
	return "strings.xml"
}

var xmlPathNoLocaleFormat = regexp.MustCompile(`(?i)/values/strings.xml`)
var xmlLocaleFromPathFormat = regexp.MustCompile(`(?i)/values-([a-zA-Z\-_]*)/strings.xml`)

func (f *xmlFormat) extractLocaleFromPath(c *phrase.Client, path string) (string, error) {
	if xmlPathNoLocaleFormat.MatchString(path) {
		return findDefaultLocaleName(c)
	}
	res := xmlLocaleFromPathFormat.FindStringSubmatch(path)
	if res == nil || len(res) < 1 {
		return "", nil
	}
	if strings.Index(res[1], "-r") != -1 {
		localePart := strings.Split(res[1], "-r")
		return fmt.Sprintf("%s-%s", localePart[0], localePart[len(localePart)-1]), nil
	} else {
		return res[1], nil
	}
}

type stringsFormat struct {
	*formatProperties
}

func (f *stringsFormat) properties() *formatProperties {
	return f.formatProperties
}

func (f *stringsFormat) directoryForLocale(c *Config, l *phrase.Locale) string {
	name := l.Code
	if name == "" {
		name = l.Name
	}
	return fmt.Sprintf("%s.lproj", f.formatName(name))
}

func (f *stringsFormat) formatName(s string) string {
	if strings.Index(s, "-") == -1 {
		return s
	}
	parts := strings.Split(s, "-")
	if strings.Index(strings.ToLower(parts[0]), "zh") != -1 {
		last := parts[len(parts)-1]
		r, n := utf8.DecodeRuneInString(last)
		last = string(unicode.ToUpper(r)) + last[n:]
		return fmt.Sprintf("%s-%s", strings.ToLower(parts[0]), last)
	} else {
		return fmt.Sprintf("%s-%s", parts[0], strings.ToUpper(parts[len(parts)-1]))
	}
}

func (f *stringsFormat) filenameForLocale(c *Config, l *phrase.Locale) string {
	return "Localizable.strings"
}

var stringsLocaleFromPathFormat = regexp.MustCompile(`(?i)/([a-zA-Z\-_]*).lproj/`)

func (f *stringsFormat) extractLocaleFromPath(_ *phrase.Client, path string) (string, error) {
	res := stringsLocaleFromPathFormat.FindStringSubmatch(path)
	if res == nil || len(res) < 1 {
		return "", nil
	}
	return res[1], nil
}

type stringsdictFormat struct {
	*stringsFormat
}

func (f *stringsdictFormat) filenameForLocale(c *Config, l *phrase.Locale) string {
	return "Localizable.stringsdict"
}

func newDefaultFormat(f string, aware bool) format {
	props := &formatProperties{
		localeAware:     aware,
		extensions:      []string{f},
		filenameFormat:  fmt.Sprintf("phrase.<locale.name>.%s", f),
		directoryFormat: "./",
	}
	return &defaultFormat{props}
}

var formats map[string]format = map[string]format{
	"json": newDefaultFormat("json", false),
	"csv":  newDefaultFormat("csv", false),
	"gettext": &defaultFormat{
		formatProperties: &formatProperties{
			localeAware:     true,
			extensions:      []string{"po"},
			targetDirectory: "locales/",
			directoryFormat: "./<locale.name>/",
			filenameFormat:  "<domain>.po",
		},
	},
	"gettext_template": &defaultFormat{
		formatProperties: &formatProperties{
			localeAware:     false,
			extensions:      []string{"pot"},
			directoryFormat: "./",
			filenameFormat:  "phrase.pot",
		},
	},
	"ini":               newDefaultFormat("ini", false),
	"properties":        newDefaultFormat("properties", true),
	"properties_xml":    newDefaultFormat("xml", false),
	"plist":             newDefaultFormat("plist", true),
	"qph":               newDefaultFormat("qph", true),
	"ts":                newDefaultFormat("ts", true),
	"resx":              newDefaultFormat("resx", false),
	"resx_windowsphone": newDefaultFormat("resx", false),
	"windows8_resource": newDefaultFormat("resw", false),
	"simple_json":       newDefaultFormat("json", false),
	"nested_json":       newDefaultFormat("json", false),
	"node_json": &defaultFormat{
		formatProperties: &formatProperties{
			extensions:      []string{"js"},
			filenameFormat:  "<locale.name>.js",
			targetDirectory: "locales/",
			directoryFormat: "./",
		},
	},
	"strings": &stringsFormat{
		formatProperties: &formatProperties{
			localeAware:     true,
			extensions:      []string{"strings"},
			targetDirectory: "./",
		},
	},
	"stringsdict": &stringsdictFormat{
		&stringsFormat{
			formatProperties: &formatProperties{
				localeAware:     true,
				extensions:      []string{"stringsdict"},
				targetDirectory: "./",
			},
		},
	},
	"xml": &xmlFormat{
		formatProperties: &formatProperties{
			localeAware:     true,
			extensions:      []string{"xml"},
			targetDirectory: "res/",
		},
	},
	"xlf": &defaultFormat{
		formatProperties: &formatProperties{
			localeAware:     true,
			extensions:      []string{"xlf", "xliff"},
			filenameFormat:  "phrase.<locale.name>.xlf",
			directoryFormat: "./",
		},
	},
	"tmx":                newDefaultFormat("tmx", false),
	"yml":                newDefaultFormat("yml", true),
	"yml_symfony":        newDefaultFormat("yml", false),
	"yml_symfony2":       newDefaultFormat("yml", false),
	"php_array":          newDefaultFormat("php", false),
	"angular_translate":  newDefaultFormat("json", false),
	"laravel":            newDefaultFormat("php", false),
	"mozilla_properties": newDefaultFormat("properties", true),
	"go_i18n": &defaultFormat{
		formatProperties: &formatProperties{
			extensions:      []string{"json"},
			targetDirectory: "locales/",
			filenameFormat:  "<locale.name>.all.json",
			directoryFormat: "./",
		},
	},
	"play_properties": &defaultFormat{
		formatProperties: &formatProperties{
			filenameFormat:  "messages.<locale.code>",
			directoryFormat: "./",
			localeExtension: true,
		},
	},
}
