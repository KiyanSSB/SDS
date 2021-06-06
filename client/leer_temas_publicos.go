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

	respuestilla := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal([]byte(byteValue), &respuestilla)
	//fmt.Println(respuestilla.Msg)

	tem := registryTema{Key: nil, Temas: nil}
	json.Unmarshal([]byte(respuestilla.Msg), &tem)

	//Printea el contenido que tenemos en tem.Temas
	/*for k := range tem.Temas {
		fmt.Println(k)
		fmt.Println("Titulo: " + tem.Temas[k].Titulo)
		fmt.Println("Descripcion: " + tem.Temas[k].Descripcion)
		fmt.Println("Comentarios: ")
		for w := range gComentarios {
			fmt.Println(tem.Temas[k].Comentarios[w].Comentario)
		}
		fmt.Println("")
	}*/

	gTemas = tem.Temas

	for k := range gTemas {
		fmt.Println("Titulo: " + gTemas[k].Titulo)
		fmt.Println("Descripcion: " + gTemas[k].Descripcion)
		fmt.Println("Comentarios: ")
		for w := range gComentarios {
			fmt.Println(gTemas[k].Comentarios[w].Comentario)
		}
		fmt.Println("")
	}
	Opciones(resp)
}
