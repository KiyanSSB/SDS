package main

import (
	"encoding/json"
	"net/http"
)

func crear_comentario(w http.ResponseWriter, req *http.Request) {
	if req != nil {

		desencript := decode64(req.Form.Get("json"))
		patata := decrypt([]byte(desencript), []byte(req.Form.Get("name")))

		jsonString := (string(patata)) //Convertimos el valor a string porque está en los valores raros

		if err := json.Unmarshal([]byte(jsonString), &gTemas); err != nil {
			panic(err)
		}

		almacenarTema()
		response(w, true, "Comentario añadido a la base de datos")
	}
}
