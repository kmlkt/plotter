package pltt

import (
	"fmt"
	"net/http"
)

func keysPost(w http.ResponseWriter, r *http.Request) error {
	var readKey, writeKey string
	parsedArgs, _ := fmt.Fscanf(r.Body, "%s %s", &readKey, &writeKey)
	var err error
	if parsedArgs == 2 {
		_, err = createKeyPair(readKey, writeKey)
	} else {
		readKey, writeKey, err = generateKeyPair()
	}
	if err != nil {
		return err
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "%s %s", readKey, writeKey)
	return nil
}
