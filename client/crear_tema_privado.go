package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

func crear_tema_privado(resp Resp) {
	fmt.Println("Crear un tema PRIVADO")
	fmt.Println("-------------")

	//COmprobar que el tema no está repetido
	if gTemas == nil {
		gTemas = make(map[string]tema)
	}

	t := tema{}
	fmt.Print("Título: ")
	t.Titulo = leerTerminal()
	fmt.Print("Descripción: ")
	t.Descripcion = leerTerminal()

	//Comprobamos que el tema no está repetido
	_, ok := gTemas[t.Titulo]
	if ok {
		fmt.Println("Este tema ya existe")
	} else {
		gTemas[t.Titulo] = t
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(t)

	println(reqBodyBytes.Bytes())

	//K del tema aleatoria para poder cifrar los temas
	haseo := sha512.New()
	key := haseo.Sum(nil)
	keyTema := sha512.Sum512([]byte(key))
	keyTema64 := keyTema[:32]

	//Apartado de selección de gente que lo puede leer
	fmt.Println("Escribe el nombre del usuario con el que lo quieres compartir")
	nombreUsuario := leerTerminal()

	//Enviamos la petición al servidor para que nos devuelva la clave publica del usuario
	data := url.Values{}
	data.Set("cmd", "damepubkey")
	data.Set("user", nombreUsuario)
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	respuestatonta := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &respuestatonta)
	fmt.Println("Clave pública del usuario elegido: ", respuestatonta.Msg)
	fmt.Println("Codificando contraseña aleatoria de tema con clave pública del usuario")
	//->>>
	temaAenvia := encode64(encrypt(compress(reqBodyBytes.Bytes()), keyTema64))
	println(temaAenvia)

	println("estás detrás de temaEnvia")

	temaJSON, err := json.Marshal(t)

	//Enviamos la contraseña random encriptada con la clave pública al servidor para que la almacena
	data2 := url.Values{}
	data2.Set("cmd", "guardarKey")
	data2.Set("user", nombreUsuario)
	data2.Set("tema", t.Titulo)
	data2.Set("encriptado", encode64(encrypt(compress([]byte(respuestatonta.Msg)), keyTema64)))
	data2.Set("file", encode64(encrypt(compress([]byte(temaJSON)), keyTema64))) //Tema encoded

	//Variable para chckear que funciona
	ficheroRandom := encode64(encrypt(compress([]byte(temaJSON)), keyTema64))

	println("estás detrás de fichero Random")

	r2, err := client.PostForm("https://localhost:10443", data2)
	chk(err)

	respuestaSiCreado := Resp{}
	byteValue, _ = ioutil.ReadAll(r2.Body)
	json.Unmarshal([]byte(byteValue), &respuestaSiCreado)
	print("ESta es la clave privada del usuario : ", respuestaSiCreado.Msg)

	//Desencriptar la llave privada del usuario
	privateKet := encode64(decrypt(decompress(decode64(respuestaSiCreado.Msg)), u.KeyData))
	var pk rsa.PrivateKey

	println("hola que tal ")

	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(&pk)

	json.Unmarshal(decode64(privateKet), &pk)

	_ = ioutil.WriteFile("datos_"+nombreUsuario+"_"+t.Titulo+".json.enc", decode64(ficheroRandom), 0644)

	fichetiro, err := ioutil.ReadFile("datos_" + nombreUsuario + "_" + t.Titulo + ".json.enc")
	chk(err)
	ficheritoEnconded := encode64(fichetiro)

	//Igual hace falta un unmarshal

	temaDesencriptado := decrypt(decode64(ficheritoEnconded), privateKeyBytes)

	var temaFinal tema

	json.Unmarshal(temaDesencriptado, &temaFinal)

	fmt.Println(temaDesencriptado)

	Opciones(resp)

}
