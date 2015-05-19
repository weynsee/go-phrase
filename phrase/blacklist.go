package phrase

// BlacklistService provides access to the blacklist related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/blacklisted_keys/
type BlacklistService struct {
	client *Client
}

type blacklistKey struct {
	Name string
}

// List all blacklisted keys.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/blacklisted_keys/
func (s *BlacklistService) Keys() ([]string, error) {
	req, err := s.client.NewRequest("GET", "blacklisted_keys", nil)
	if err != nil {
		return nil, err
	}

	blacklists := new([]blacklistKey)
	_, err = s.client.Do(req, blacklists)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(*blacklists))
	for i, key := range *blacklists {
		keys[i] = key.Name
	}

	return keys, err
}
