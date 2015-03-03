package phrase

import (
	"fmt"
	"net/url"
)

// SessionsService provides access to the authentication related functions
// in the Phrase API.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/authentication/
type SessionsService struct {
	client *Client
}

type session struct {
	Success bool   `json:"success"`
	Token   string `json:"auth_token"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role_name"`
}

type checkLoginResponse struct {
	LoggedIn bool  `json:"logged_in"`
	User     *User `json:"user"`
}

// Sign in a user identified by email and password.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/authentication/#create
func (s *SessionsService) Create(email, password string) (string, error) {
	params := url.Values{}
	params.Set("email", email)
	params.Set("password", password)
	req, err := s.client.NewRequest("POST", "sessions", params)
	if err != nil {
		return "", err
	}

	sess := new(session)
	_, err = s.client.Do(req, sess)
	if err != nil {
		return "", err
	}

	return sess.Token, err
}

// Log the current user out.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/authentication/#destroy
func (s *SessionsService) Destroy() error {
	req, err := s.client.NewRequest("DELETE", "sessions", nil)
	if err != nil {
		return err
	}

	sess := new(session)
	_, err = s.client.Do(req, sess)

	return err
}

// Check the validity of an auth_token and return information of the current user.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/authentication/#check_login
func (s *SessionsService) CheckLogin() (*User, error) {
	req, err := s.client.NewRequest("GET", "auth/check_login", nil)
	if err != nil {
		return nil, err
	}

	loggedIn := new(checkLoginResponse)
	_, err = s.client.Do(req, loggedIn)
	if err != nil {
		return nil, err
	}

	return loggedIn.User, err
}

func (u User) String() string {
	return fmt.Sprintf("User ID: %d Username: %s",
		u.ID, u.Username)
}
