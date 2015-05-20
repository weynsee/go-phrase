# go-phrase #

go-phrase is a Go client library for accessing the [PhraseApp API](http://docs.phraseapp.com/api/v1/). It also includes a command line interface that implements all the commands in the [official PhraseApp command-line client](https://github.com/phrase/phrase).

**Documentation:** [![GoDoc](https://godoc.org/github.com/weynsee/go-phrase?status.svg)](https://godoc.org/github.com/weynsee/go-phrase)  
**Build Status:** [![Build Status](https://travis-ci.org/weynsee/go-phrase.svg?branch=master)](https://travis-ci.org/weynsee/go-phrase)  
**Test Coverage:** [![Test Coverage](https://coveralls.io/repos/weynsee/go-phrase/badge.svg?branch=master)](https://coveralls.io/r/weynsee/go-phrase?branch=master)

go-phrase requires Go version 1.1 or greater.

## CLI ##

The CLI has 4 commands:

```
    init    Initializes a phrase project
    pull    Download the translation files in the current project
    push    Upload the translation files in the current project to PhraseApp
    tags    List all the tags in the current project
```

Options and arguments for the commands are the same those used in the [official command-line client](https://github.com/phrase/phrase).

## API ##

```go
import "github.com/weynsee/go-phrase/phrase"
```

Construct a new API client, then use the various services on the client to
access different parts of the PhraseApp API.  For example, to list all
the locales for your project token:

```go
client := phrase.New(token)
locales, err := client.Locales.ListAll()
```
## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.
