package phrase

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"strconv"
)

type KeysService struct {
	client *Client
}

type Key struct {
	ID               int      `json:"id" url:",omitempty"`
	Name             string   `json:"name" url:"translation_key[name]"`
	NamePlural       string   `json:"name_plural" url:"translation_key[name_plural],omitempty"`
	Description      string   `json:"description" url:"translation_key[description]"`
	Pluralized       bool     `json:"pluralized" url:"translation_key[pluralized],int,omitempty"`
	DataType         string   `json:"data_type" url:"translation_key[data_type],omitempty"`
	Tags             []string `json:"tag_list" url:"translation_key[tag_names],comma"`
	Unformatted      bool     `json:"unformatted" url:"translation_key[unformatted],int,omitempty"`
	MaxCharacters    int      `json:"max_characters_allowed" url:"translation_key[max_characters_allowed]"`
	XMLSpacePreserve bool     `json:"xml_space_preserve" url:"translation_key[xml_space_preserve],int,omitempty"`
}

type translateResponse struct {
	Success   bool            `json:"success"`
	Translate json.RawMessage `json:"translate"`
}

type KeyTranslation struct {
	String string
	Map    map[string]string
}

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

func (s *KeysService) ListAll() ([]Key, error) {
	return s.Get(nil)
}

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

func (s *KeysService) Create(k *Key) (*Key, error) {
	return s.submitKey("POST", "translation_keys", k)
}

func (s *KeysService) Update(k *Key) (*Key, error) {
	u := fmt.Sprintf("translation_keys/%d", k.ID)
	return s.submitKey("PATCH", u, k)
}

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
