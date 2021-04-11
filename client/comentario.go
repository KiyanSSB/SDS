package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

func comentario(cmd string, resp Resp) {
	fmt.Print("¿A qué tema quieres añadir el comentario? ")
	tituloTema := leerTerminal()

	vacio := false
	t := tema{}
	c := Comentario{}
	for k := range gTemas {
		if tituloTema == gTemas[k].Titulo {
			vacio = true
			if gComentarios == nil {
				gComentarios = make(map[string]Comentario)
			}
			fmt.Print("Escribe tu comentario: ")
			c.Comentario = leerTerminal()
			gComentarios[c.Comentario] = c
			t.Titulo = gTemas[k].Titulo
			t.Descripcion = gTemas[k].Descripcion
			t.Comentarios = gComentarios
			gTemas[k] = t
		}
	}

	if !vacio {
		fmt.Println("No existe el tema")
	}

	jsonData, err := json.Marshal(&gTemas)
	chk(err)
	fmt.Println(gTemas)
	jsonData = []byte(encode64(encrypt(jsonData, u.KeyData)))

	data := url.Values{}
	data.Set("cmd", "crear_tema")
	data.Set("json", string(jsonData))
	data.Set("name", string(u.KeyData))
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	resp2 := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp2)
	Opciones(resp2)
}
