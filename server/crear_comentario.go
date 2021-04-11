package main

import (
	"net/http"
)

func crear_comentario(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		almacenarTema()
		response(w, true, "Comentario añadido a la base de datos")
	}
}
