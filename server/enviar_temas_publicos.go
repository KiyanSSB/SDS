package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func enviar_temas_publicos(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("temas.json")

	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero en la variable en bytes

	fmt.Println(string(byteValue))

	json.Marshal(byteValue)

	response(w, true, string(byteValue))
}
