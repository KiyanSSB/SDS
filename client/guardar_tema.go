package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

func guardar_tema(cmd string, resp Resp) {
	fmt.Println("Crear un tema")
	fmt.Println("-------------")

	if gTemas == nil {
		gTemas = make(map[string]tema)
	}

	t := tema{}
	fmt.Print("Título: ")
	t.Titulo = leerTerminal()
	fmt.Print("Descripción: ")
	t.Descripcion = leerTerminal()

	_, ok := gTemas[t.Titulo]
	if ok {
		fmt.Println("Este tema ya existe")
	} else {
		gTemas[t.Titulo] = t
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
