package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

func crear_tema_privado(resp Resp) {
	fmt.Println("Crear un tema PRIVADO")
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

	//K del tema aleatoria para poder cifrar los temas
	haseo := sha512.New()
	key := haseo.Sum(nil)
	keyTema := sha512.Sum512([]byte(key))
	keyTema64 := keyTema[:32]

	//Apartado de selección de gente que lo puede leer
	fmt.Println("Escribe el nombre del usuario con el que lo quieres compartir")
	nombreUsuario := leerTerminal()

	//Enviamos la petición al servidor para que nos devuelva la clave privada del usuario
	data := url.Values{}
	data.Set("cmd", "damepubkey")
	data.Set("user", nombreUsuario)
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	respuestatonta := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &respuestatonta)
	fmt.Println(respuestatonta.Msg)
	fmt.Println("Codificando contraseña aleatoria de tema con clave pública del usuario")

	//Enviamos la contraseña random encriptada con la clave pública al servidor para que la almacena
	data2 := url.Values{}
	data2.Set("cmd", "guardarKey")
	data2.Set("user", nombreUsuario)
	data2.Set("tema", t.Titulo)

	data2.Set("encriptado", encode64(encrypt(compress([]byte(respuestatonta.Msg)), keyTema64)))

	r2, err := client.PostForm("https://localhost:10443", data2)
	chk(err)

	respuestaSiCreado := Resp{}
	byteValue, _ = ioutil.ReadAll(r2.Body)
	json.Unmarshal([]byte(byteValue), &respuestaSiCreado)
	print(respuestaSiCreado.Msg)

	Opciones(resp)

}
