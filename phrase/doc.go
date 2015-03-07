/*
Package phrase provides a client for using the Phrase API.

Construct a new Phrase client, then use the various services on the client to
access different parts of the Phrase API. For example:

	client := phrase.New(token)

	// list all the locales in the current project
	locales, err := client.Locales.ListAll()

The services of a client correspond to
the structure of the Phrase API documentation at
http://docs.phraseapp.com/api/v1/.

Authentication

The phrase client sends the authentication token (obtained from your project
overview page in Phrase)  for all API requests,
but some requests require that you perform a user login before it can be
performed. These requests are marked as signed requests in the documentation.
To do a user login, use the sessions service:

	userAuthToken, err := client.Sessions.Create(email, password)

The userAuthToken can now be used to create a new Phrase client that
can be used to do signed API requests:

	newClient, err := phrase.NewClient(userAuthToken, token, nil)

For more information on authentication, see http://docs.phraseapp.com/api/v1/authentication/
*/
package phrase
