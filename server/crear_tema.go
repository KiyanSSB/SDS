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

		jsonString := (string(patata)) //Convertimos el valor a string porque está en los valores raros
		fmt.Println(jsonString)

		//Guardamos en un tema T
		var t tema
		if err := json.Unmarshal([]byte(jsonString), &t); err != nil {
			panic(err)
		}

		fmt.Println("Este es el tema", t.Titulo)
		fmt.Println("Esta es la descripción: ", t.Descripcion)

		gTemas[t.Titulo] = t

		/*var jsonchulo gTemas[]

		json.Unmarshal([]byte(string(patata)), &jsonchulo)

		fmt.Println(jsonchulo)

		fmt.Println(jsonchulo.Titulo)*/

		//des64 := (decode64((string(desencript))))

		//fmt.Println([]byte(des64))

		//temilla := tema{}

		almacenarTema()
		response(w, true, "Añadido a la base de datos")
	}
}
