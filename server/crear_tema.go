package main

import (
	"net/http"
)

func crear_tema(w http.ResponseWriter, req *http.Request) {
	if req != nil {
		almacenarTema()
		response(w, true, "Añadido a la base de datos")
	}
}
