package netjson

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Posts) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}

func (u *Users) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}

func (f *Followers) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}
