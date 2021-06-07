package main

import (
	"encoding/json"
	"fmt"
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

	resp2 := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal([]byte(byteValue), &resp2)

	tem := registryTema{Key: nil, Temas: nil}
	json.Unmarshal([]byte(resp2.Msg), &tem)

	gTemas = tem.Temas

	if gTemas != nil {
		for k := range gTemas {
			fmt.Println("Titulo: " + gTemas[k].Titulo)
			fmt.Println("Descripcion: " + gTemas[k].Descripcion)
			fmt.Println("Comentarios: ")
			for w := range gTemas[k].Comentarios {
				fmt.Println(gTemas[k].Comentarios[w].Comentario)
			}
		}
	} else {
		fmt.Println("No hay temas públicos disponibles")
	}

	Opciones(resp)
}
