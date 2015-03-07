# go-phrase #

go-phrase is a Go client library for accessing the [Phrase API](http://docs.phraseapp.com/api/v1/).

**Documentation:** [![GoDoc](https://godoc.org/github.com/weynsee/go-phrase?status.svg)](https://godoc.org/github.com/weynsee/go-phrase)  
**Build Status:** [![Build Status](https://travis-ci.org/weynsee/go-phrase.svg?branch=master)](https://travis-ci.org/weynsee/go-phrase)  
**Test Coverage:** [![Test Coverage](https://coveralls.io/repos/weynsee/go-phrase/badge.svg?branch=master)](https://coveralls.io/r/weynsee/go-phrase?branch=master)

go-phrase requires Go version 1.1 or greater.

## Usage ##

```go
import "github.com/weynsee/go-phrase/phrase"
```

Construct a new Phrase client, then use the various services on the client to
access different parts of the Phrase API.  For example, to list all
the locales for your project token:

```go
client := phrase.New(token)
locales, err := client.Locales.ListAll()
```
## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.
