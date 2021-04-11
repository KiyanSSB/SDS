package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type User struct {
	Name    string
	KeyData []byte
}

var u = User{}

func signin(client *http.Client, cmd string) {
	fmt.Print("Nombre de usuario: ")
	user := leerTerminal()

	fmt.Print("Contraseña: ")
	pass := leerTerminal()

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	u.Name = user
	u.KeyData = keyData

	data := url.Values{}                 // estructura para contener los valores
	data.Set("cmd", cmd)                 // comando (string)
	data.Set("user", user)               // usuario (string)
	data.Set("pass", encode64(keyLogin)) // "contraseña" a base64

	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	//Respuesta
	resp := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp)
	Opciones(resp)
}
