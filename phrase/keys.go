package phrase

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"strconv"
)

// KeysService provides access to the translation key related functions
// in the Phrase API.
//
// http://docs.phraseapp.com/api/v1/translation_keys/
type KeysService struct {
	client *Client
}

// Key represents a translation key.
type Key struct {
	ID int `json:"id" url:",omitempty"`

	// Name of the key. This is the only required field.
	Name string `json:"name" url:"translation_key[name]"`

	// Plural name of the key, e.g. for use in Gettext format.
	NamePlural  string `json:"name_plural" url:"translation_key[name_plural],omitempty"`
	Description string `json:"description" url:"translation_key[description]"`
	Pluralized  bool   `json:"pluralized" url:"translation_key[pluralized],int,omitempty"`

	// Data type of the key. The values can be string, number, boolean, or array
	DataType string   `json:"data_type" url:"translation_key[data_type],omitempty"`
	Tags     []string `json:"tag_list" url:"translation_key[tag_names],comma"`

	// Mark key as unformatted.
	Unformatted bool `json:"unformatted" url:"translation_key[unformatted],int,omitempty"`

	// Mark xml:space="preserve" as true, i.e. for Android XML format.
	MaxCharacters int `json:"max_characters_allowed" url:"translation_key[max_characters_allowed]"`

	// Max. number of characters for this key, default is 0 (unlimited).
	XMLSpacePreserve bool `json:"xml_space_preserve" url:"translation_key[xml_space_preserve],int,omitempty"`
}

type translateResponse struct {
	Success   bool            `json:"success"`
	Translate json.RawMessage `json:"translate"`
}

// KeyTranslation represents a response from the Translate API call.
// The response can be either a String or a Map, which are stored
// in the appropriate fields.
type KeyTranslation struct {
	String string
	Map    map[string]string
}

// UploadRequest represents a request to the Upload API call.
type UploadRequest struct {
	Filename           string   `url:"filename"`
	FileContent        string   `url:"file_content"`
	Tags               []string `url:"tags[],omitempty"`
	LocaleCode         string   `url:"locale_code"`
	FileFormat         string   `url:"file_format,omitempty"`
	UpdateTranslations bool     `url:"update_translations,int,omitempty"`
	SkipUnverification bool     `url:"skip_unverification,int,omitempty"`
	SkipUploadTags     bool     `url:"skip_upload_tags,int,omitempty"`
}

// List all keys for your current project.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#index
func (s *KeysService) ListAll() ([]Key, error) {
	return s.Get(nil)
}

// Return only keys that match the given names. This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#index
func (s *KeysService) Get(keyNames []string) ([]Key, error) {
	params := url.Values{}
	for _, x := range keyNames {
		params.Add("key_names[]", x)
	}
	req, err := s.client.NewRequest("GET", "translation_keys", params)
	if err != nil {
		return nil, err
	}

	keys := new([]Key)
	_, err = s.client.Do(req, keys)
	if err != nil {
		return nil, err
	}

	return *keys, err
}

// Create a new key for the current project. This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#create
func (s *KeysService) Create(k *Key) (*Key, error) {
	return s.submitKey("POST", "translation_keys", k)
}

// Update an existing key in the current project. This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#update
func (s *KeysService) Update(k *Key) (*Key, error) {
	u := fmt.Sprintf("translation_keys/%d", k.ID)
	return s.submitKey("PATCH", u, k)
}

// Delete key identified by id. Be careful: This will delete all associated translations as well!
// This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#destroy
func (s *KeysService) Destroy(id int) error {
	u := fmt.Sprintf("translation_keys/%d", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	del := new(successResponse)
	_, err = s.client.Do(req, del)
	return err
}

// Delete multiple keys identified by their ids. Be careful: This will delete all associated translations as well! The number of keys to delete is limited to 50 per request.
// This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#destroy_multiple
func (s *KeysService) DestroyMultiple(ids []int) error {
	params := url.Values{}
	for _, x := range ids {
		params.Add("ids[]", strconv.Itoa(x))
	}
	req, err := s.client.NewRequest("DELETE", "translation_keys/destroy_multiple", params)
	if err != nil {
		return err
	}

	del := new(successResponse)
	_, err = s.client.Do(req, del)

	return err
}

// Returns all untranslated keys for the given locale. Note: This will return 10,000 keys at most.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#untranslated
func (s *KeysService) ListUntranslated(l string) ([]Key, error) {
	params := url.Values{}
	params.Set("locale_name", l)
	req, err := s.client.NewRequest("GET", "translation_keys/untranslated", params)
	if err != nil {
		return nil, err
	}

	keys := new([]Key)
	_, err = s.client.Do(req, keys)
	if err != nil {
		return nil, err
	}

	return *keys, err
}

// Add tags to the given keys. Existing tags for the given keys will not be removed.
// This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#tag
func (s *KeysService) Tag(ids []int, tags []string) error {
	params := url.Values{}
	for _, id := range ids {
		params.Add("ids[]", strconv.Itoa(id))
	}
	for _, tag := range tags {
		params.Add("tags[]", tag)
	}
	req, err := s.client.NewRequest("POST", "translation_keys/tag", params)
	if err != nil {
		return err
	}

	resp := new(successResponse)
	_, err = s.client.Do(req, resp)

	return err
}

// Returns what I18n.translate(key) would return in Ruby in JSON format. This is to be able to provide data if the translation has not been used in a pure key-value-fashion, but to store an array or even a hash of its child items.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#translate
func (s *KeysService) Translate(key string) (*KeyTranslation, error) {
	params := url.Values{}
	params.Set("key", key)
	req, err := s.client.NewRequest("GET", "translation_keys/translate", params)
	if err != nil {
		return nil, err
	}

	resp := new(translateResponse)
	_, err = s.client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.Success {
		translation, err := decodeTranslation(resp.Translate)
		return translation, err
	} else {
		return nil, err
	}
}

// Upload a localization file to the current project. This will add new keys and their translations.
// This is a signed request.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/translation_keys/#upload
func (s *KeysService) Upload(u *UploadRequest) error {
	params, err := query.Values(u)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest("POST", "translation_keys/upload", params)
	if err != nil {
		return err
	}

	resp := new(successResponse)
	_, err = s.client.Do(req, resp)

	return err
}

func decodeTranslation(data json.RawMessage) (*KeyTranslation, error) {
	translation := new(KeyTranslation)
	var s string
	var m map[string]string
	var err error
	if err = json.Unmarshal(data, &s); err == nil {
		translation.String = s
		return translation, err
	}
	if err = json.Unmarshal(data, &m); err == nil {
		translation.Map = m
		return translation, err
	}
	return translation, err
}

func (s *KeysService) submitKey(method, url string, k *Key) (*Key, error) {
	params, err := query.Values(k)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}

	key := new(Key)
	_, err = s.client.Do(req, key)
	if err != nil {
		return nil, err
	}

	return key, err
}

func (k Key) String() string {
	return fmt.Sprintf("Key ID: %d Name: %s Description: %s Type: %s Tags: %v",
		k.ID, k.Name, k.Description, k.DataType, k.Tags)
}
