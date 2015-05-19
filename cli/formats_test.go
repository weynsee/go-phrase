package cli

import (
	"fmt"
	"github.com/weynsee/go-phrase/phrase"
	"net/http"
	"reflect"
	"testing"
)

func testExtensions(t *testing.T, name string, props *formatProperties, extensions []string) {
	if !reflect.DeepEqual(props.extensions, extensions) {
		t.Errorf("%s format supports %s extensions", name, extensions)
	}
}

func testLocaleAware(t *testing.T, name string, props *formatProperties, expected bool) {
	if got := props.localeAware; got != expected {
		t.Errorf("%s format locale aware should be %v", name, expected)
	}
}

func testFilenameForLocale(t *testing.T, name, filename string, f format) {
	if got := f.filenameForLocale(&Config{Domain: "testing"}, &phrase.Locale{Name: "de", Code: "de"}); got != filename {

		t.Errorf("%s format filename expects %s, got %s", name, filename, got)
	}
}

func testDirectoryForLocale(t *testing.T, name, directory string, f format) {
	if got := f.directoryForLocale(&Config{}, &phrase.Locale{Name: "fr"}); got != directory {
		t.Errorf("%s format directory expects %s, got %s", name, directory, got)
	}
}

func testTargetDirectory(t *testing.T, name string, props *formatProperties, target string) {
	if got := props.targetDirectory; got != target {
		t.Errorf("%s format target directory should be %v", name, target)
	}
}

func testSimpleFormat(t *testing.T, name, format string, localeAware bool) {
	f := formats[name]
	props := f.properties()
	testExtensions(t, name, props, []string{format})
	testLocaleAware(t, name, props, localeAware)
	filename := fmt.Sprintf("phrase.%s.%s", "de", format)
	testFilenameForLocale(t, name, filename, f)
	testDirectoryForLocale(t, name, "./", f)
	if got, _ := f.extractLocaleFromPath(&phrase.Client{}, ""); got != "" {
		t.Errorf("%s format cannot extract locale from path", name)
	}
}

func TestFormats_json(t *testing.T) {
	testSimpleFormat(t, "json", "json", false)
}

func TestFormats_csv(t *testing.T) {
	testSimpleFormat(t, "csv", "csv", false)
}

func TestFormats_ini(t *testing.T) {
	testSimpleFormat(t, "ini", "ini", false)
}

func TestFormats_yml(t *testing.T) {
	testSimpleFormat(t, "yml", "yml", true)
}

func TestFormats_tmx(t *testing.T) {
	testSimpleFormat(t, "tmx", "tmx", false)
}

func TestFormats_yml_symfony(t *testing.T) {
	testSimpleFormat(t, "yml_symfony", "yml", false)
}

func TestFormats_yml_symfony2(t *testing.T) {
	testSimpleFormat(t, "yml_symfony2", "yml", false)
}

func TestFormats_php_array(t *testing.T) {
	testSimpleFormat(t, "php_array", "php", false)
}

func TestFormats_angular_translate(t *testing.T) {
	testSimpleFormat(t, "angular_translate", "json", false)
}

func TestFormats_nested_json(t *testing.T) {
	testSimpleFormat(t, "nested_json", "json", false)
}

func TestFormats_simple_json(t *testing.T) {
	testSimpleFormat(t, "simple_json", "json", false)
}

func TestFormats_laravel(t *testing.T) {
	testSimpleFormat(t, "laravel", "php", false)
}

func TestFormats_properties(t *testing.T) {
	testSimpleFormat(t, "properties", "properties", true)
}

func TestFormats_properties_xml(t *testing.T) {
	testSimpleFormat(t, "properties_xml", "xml", false)
}

func TestFormats_mozilla_properties(t *testing.T) {
	testSimpleFormat(t, "mozilla_properties", "properties", true)
}

func TestFormats_qph(t *testing.T) {
	testSimpleFormat(t, "qph", "qph", true)
}

func TestFormats_plist(t *testing.T) {
	testSimpleFormat(t, "plist", "plist", true)
}

func TestFormats_ts(t *testing.T) {
	testSimpleFormat(t, "ts", "ts", true)
}

func TestFormats_resx(t *testing.T) {
	testSimpleFormat(t, "resx", "resx", false)
}

func TestFormats_resx_windowsphone(t *testing.T) {
	testSimpleFormat(t, "resx_windowsphone", "resx", false)
}

func TestFormats_resw(t *testing.T) {
	testSimpleFormat(t, "windows8_resource", "resw", false)
}

func TestFormats_xliff(t *testing.T) {
	f := formats["xlf"]
	props := f.properties()
	name := "xliff"
	extensions := []string{"xlf", "xliff"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, true)
	filename := fmt.Sprintf("phrase.%s.%s", "de", "xlf")
	testFilenameForLocale(t, name, filename, f)
	testDirectoryForLocale(t, name, "./", f)
}

func TestFormats_go_i18n(t *testing.T) {
	f := formats["go_i18n"]
	props := f.properties()
	name := "go_i18n"
	extensions := []string{"json"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, false)
	testFilenameForLocale(t, name, "de.all.json", f)
	testDirectoryForLocale(t, name, "./", f)
	testTargetDirectory(t, name, props, "locales/")
}

func TestFormats_gettext(t *testing.T) {
	f := formats["gettext"]
	props := f.properties()
	name := "gettext"
	extensions := []string{"po"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, true)
	testFilenameForLocale(t, name, "testing.po", f)
	testDirectoryForLocale(t, name, "./fr/", f)
	testTargetDirectory(t, name, props, "locales/")
}

func TestFormats_gettext_template(t *testing.T) {
	f := formats["gettext_template"]
	props := f.properties()
	name := "gettext_template"
	extensions := []string{"pot"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, false)
	testFilenameForLocale(t, name, "phrase.pot", f)
	testDirectoryForLocale(t, name, "./", f)
}

func TestFormats_play_properties(t *testing.T) {
	f := formats["play_properties"]
	props := f.properties()
	name := "play_properties"
	testLocaleAware(t, name, props, false)
	testFilenameForLocale(t, name, "messages.de", f)
	testDirectoryForLocale(t, name, "./", f)
	if !props.localeExtension {
		t.Error("play_properties renders locale as extension")
	}
}

func TestFormats_node_json(t *testing.T) {
	f := formats["node_json"]
	props := f.properties()
	name := "node_json"
	extensions := []string{"js"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, false)
	testFilenameForLocale(t, name, "de.js", f)
	testDirectoryForLocale(t, name, "./", f)
	testTargetDirectory(t, name, props, "locales/")
}

func TestFormats_strings(t *testing.T) {
	f := formats["strings"]
	props := f.properties()
	name := "strings"
	extensions := []string{"strings"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, true)
	testFilenameForLocale(t, name, "Localizable.strings", f)
	testDirectoryForLocale(t, name, "fr.lproj", f)
	directory := "foo-BAR.lproj"
	if got := f.directoryForLocale(&Config{}, &phrase.Locale{Name: "foo-bar"}); got != directory {
		t.Errorf("%s format directory expects %s, got %s", name, directory, got)
	}
	directory = "zh-Cn.lproj"
	if got := f.directoryForLocale(&Config{}, &phrase.Locale{Name: "ZH-cn"}); got != directory {
		t.Errorf("%s format directory expects %s, got %s", name, directory, got)
	}
	locale := "fr_FR"
	if got, _ := f.extractLocaleFromPath(nil, "/foo/fr_FR.lproj/Localizable.strings"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	locale = ""
	if got, _ := f.extractLocaleFromPath(nil, "/foo/bar/Localizable.strings"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
}

func TestFormats_stringsdict(t *testing.T) {
	f := formats["stringsdict"]
	props := f.properties()
	name := "stringsdict"
	extensions := []string{"stringsdict"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, true)
	testFilenameForLocale(t, name, "Localizable.stringsdict", f)
	testDirectoryForLocale(t, name, "fr.lproj", f)
	directory := "foo-BAR.lproj"
	if got := f.directoryForLocale(&Config{}, &phrase.Locale{Name: "foo-bar"}); got != directory {
		t.Errorf("%s format directory expects %s, got %s", name, directory, got)
	}
	directory = "zh-Cn.lproj"
	if got := f.directoryForLocale(&Config{}, &phrase.Locale{Name: "ZH-cn"}); got != directory {
		t.Errorf("%s format directory expects %s, got %s", name, directory, got)
	}
	locale := "fr_FR"
	if got, _ := f.extractLocaleFromPath(nil, "/foo/fr_FR.lproj/Localizable.strings"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	locale = ""
	if got, _ := f.extractLocaleFromPath(nil, "/foo/bar/Localizable.strings"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
}

func TestFormats_xml(t *testing.T) {
	f := formats["xml"]
	props := f.properties()
	name := "xml"
	extensions := []string{"xml"}
	testExtensions(t, name, props, extensions)
	testLocaleAware(t, name, props, true)
	testFilenameForLocale(t, name, "strings.xml", f)
}

func TestFormats_xml_extractLocaleFromPath(t *testing.T) {
	f := formats["xml"]
	name := "xml"
	locale := "fr"
	if got, _ := f.extractLocaleFromPath(nil, "/foo/values-fr/strings.xml"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	locale = "de-DE"
	if got, _ := f.extractLocaleFromPath(nil, "/foo/values-de-DE/strings.xml"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	locale = "pt-BR"
	if got, _ := f.extractLocaleFromPath(nil, "/foo/values-pt-rBR/strings.xml"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	locale = ""
	if got, _ := f.extractLocaleFromPath(nil, "/foo/bar/strings.xml"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
	setupAPI()
	defer shutdownAPI()
	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"name":"default","is_default":true}]`)
	})
	locale = "default"
	if got, _ := f.extractLocaleFromPath(client, "/foo/values/strings.xml"); got != locale {
		t.Errorf("%s format extractLocaleFromPath expects %s, got %s", name, locale, got)
	}
}

func TestFormats_xml_directoryForLocale(t *testing.T) {
	f := formats["xml"]
	name := "xml"
	locale, directory := &phrase.Locale{Default: true}, "values"
	if got := f.directoryForLocale(nil, locale); got != directory {
		t.Errorf("%s format directoryForLocale expects %s, got %s", name, directory, got)
	}
	locale, directory = &phrase.Locale{Name: "fooish", Default: false}, "values-fooish"
	if got := f.directoryForLocale(nil, locale); got != directory {
		t.Errorf("%s format directoryForLocale expects %s, got %s", name, directory, got)
	}
	locale, directory = &phrase.Locale{Name: "foo-ish", Default: false}, "values-foo-rISH"
	if got := f.directoryForLocale(nil, locale); got != directory {
		t.Errorf("%s format directoryForLocale expects %s, got %s", name, directory, got)
	}
}
