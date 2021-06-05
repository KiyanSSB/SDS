package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
)

func leer_temas_publicos(resp Resp) {

	if gTemas == nil {
		gTemas = make(map[string]tema)
	}

	data := url.Values{}

	data.Set("cmd", "publicos")

	//Hacemos FETCH de todos los temas que sean públicos
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	respuestilla := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &respuestilla)

}
