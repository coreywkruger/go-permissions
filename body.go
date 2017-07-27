package main

import (
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/http"
)

// Body struct is a shortcut for handling request bodies
type Body struct {
	R    *http.Request
	W    http.ResponseWriter
	body *gabs.Container
}

// NewBody is a factory for Body structs
func NewBody(w http.ResponseWriter, r *http.Request) *Body {
	b := Body{R: r, W: w}
	bodyBytes, err := ioutil.ReadAll(b.R.Body)
	if err != nil {
		b.W.WriteHeader(422)
		b.W.Write([]byte("Could not process request"))
		return nil
	}
	defer b.R.Body.Close()

	b.body, err = gabs.ParseJSON(bodyBytes)
	if err != nil {
		b.W.WriteHeader(422)
		b.W.Write([]byte("Could not process request"))
		return nil
	}
	return &b
}

// GetField gets a field from a json request body
func (b *Body) GetField(field string) string {
	value, _ := b.body.Path(field).Data().(string)
	return value
}
