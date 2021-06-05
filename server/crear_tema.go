package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func crear_tema(w http.ResponseWriter, req *http.Request) {
	if req != nil {

		//t := tema{}
		//var data = req.Form.Get("json")
		//var dataReformed = decode64((string(decrypt([]byte(data), []byte(req.Form.Get("name"))))))
		//fmt.Println(dataReformed)
		//print(req.Form.Get("json"))
		//desencript := decrypt([]byte(req.Form.Get("json")), []byte(req.Form.Get("name")))

		desencript := decode64(req.Form.Get("json"))
		patata := decrypt([]byte(desencript), []byte(req.Form.Get("name")))

		fmt.Println(string(patata)) /***PARA VER EL JSON RECIBIDO***/

		jsonString := (string(patata))

		fmt.Println(jsonString)
		////////////////////////////////////////////////

		t := make(map[string]tema)
		if err := json.Unmarshal([]byte(jsonString), &t); err != nil {
			panic(err)
		}
		println(t[req.Form.Get("name")].Titulo)

		gTemas[req.Form.Get("name")] = t[req.Form.Get("name")]

		/*var jsonchulo gTemas[]

		json.Unmarshal([]byte(string(patata)), &jsonchulo)

		fmt.Println(jsonchulo)

		fmt.Println(jsonchulo.Titulo)*/

		//des64 := (decode64((string(desencript))))

		//fmt.Println([]byte(des64))

		almacenarTema()
		response(w, true, "Añadido a la base de datos")
	}
}
