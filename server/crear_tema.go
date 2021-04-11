package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

/******************************************************************************************************************************
***					En este fichero se encuentra todo lo relacionado con la creación de los temas en el servidor			***
*******************************************************************************************************************************/

func crear_tema(w http.ResponseWriter, req *http.Request) {
	os.Remove("temas.json")
	if req != nil {
		_, err := os.Create("temas.json")
		chk(err)
		jsonF := decode64(req.Form.Get("json"))
		err = ioutil.WriteFile("temas.json", jsonF, 0644)
		chk(err)
		response(w, true, "Añadido a la base de datos")
	}
}
