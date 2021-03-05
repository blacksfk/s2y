package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func PerformReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	res, e := client.Do(req)

	if e != nil {
		return nil, e
	}

	body, e := io.ReadAll(res.Body)

	if e != nil {
		return nil, e
	}

	if res.StatusCode >= 400 {
		// something went wrong with the request
		ae := &Error{}
		e = json.Unmarshal(body, ae)

		if e != nil {
			// busted JSON
			return nil, e
		}

		return nil, ae
	}

	// request was successful, return the body
	return body, nil
}
