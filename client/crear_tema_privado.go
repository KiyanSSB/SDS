package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

func crear_tema_privado(resp Resp) {
	fmt.Println("Crear un tema PRIVADO")
	fmt.Println("-------------")

	//Inicializamos el mapa de temas
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
		return
	} else {
		gTemas[t.Titulo] = t
	}

	//1º Elegir el usuario con el que lo queremos compartir
	fmt.Println("Escribe el nombre del usuario con el que lo quieres compartir")
	nombreUsuario := leerTerminal()

	//2º Recibir la clave pública del usuario seleccionado desde el servidor
	data := url.Values{}
	data.Set("cmd", "damepubkey")
	data.Set("user", nombreUsuario)
	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	respuestatonta := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &respuestatonta)
	fmt.Println("Clave pública del usuario elegido: ", respuestatonta.Msg)

	var crypto crypto.PublicKey
	file, err := os.Open("cp.json")

	byteValue, _ = ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &crypto) // -> de bytes a clave publica

	var cryptoFD []byte
	cryptoFD, err = json.Marshal(&crypto) // -> de clave publica a bytes [] de json

	println(crypto)
	println(string(cryptoFD))

	fmt.Println("\n \n \n ")

	//3º cifrar  el tema con la clave random
	haseo := sha512.New()
	key := haseo.Sum(nil)
	keyTema := sha512.Sum512([]byte(key))
	keyTema32 := keyTema[:32] //Utilizamos solo la primera parte de la clave

	temaJSON, err := json.Marshal(t)                                           //Cremos el tema en formato json
	temaEncriptado := encode64(encrypt(compress([]byte(temaJSON)), keyTema32)) //Encriptamos el tema en formato json

	_ = ioutil.WriteFile("datos_"+nombreUsuario+"_"+t.Titulo+".json.enc", decode64(temaEncriptado), 0644) //Creamos el fichero con los datos en bytes

	//4º Cifrar la clave random con la clave pública del usuario que queremos que pueda ver el tema
	fmt.Println("Codificando contraseña aleatoria de tema con clave pública del usuario ...")
	var clavePubUnmarshal rsa.PublicKey

	var clavePublica []byte
	json.Unmarshal([]byte(clavePublica), &clavePubUnmarshal)

	encClave, err := rsa.EncryptPKCS1v15(rand.Reader, &clavePubUnmarshal, keyTema32)

	println("Esta es la clave random encodeada: ", encClave)

	//-----> Hasta este momento tenemos : El fichero con el tema encriptado, la clave random del usuario + tema creada
	//-----> Necesitamos Clave PRIVADA del usuario para leer el tema

	//5º Pedir al servidor la clave privada cifrada del usuario

	//Recibir la clave privada del usuario
	data2 := url.Values{}
	data2.Set("cmd", "guardarKey")
	data2.Set("user", nombreUsuario)
	data2.Set("tema", t.Titulo)
	data2.Set("encriptado", string(encClave))                                   //Key random codificada con la clave publica del usuario
	data2.Set("file", encode64(encrypt(compress([]byte(temaJSON)), keyTema32))) //Tema creado encodeado

	r2, err := client.PostForm("https://localhost:10443", data2)
	chk(err)

	respuestaSiCreado := Resp{}                                                                //Recibimos la clave privada del usuario CODIFICADA			//
	byteValue, _ = ioutil.ReadAll(r2.Body)                                                     //
	json.Unmarshal([]byte(byteValue), &respuestaSiCreado)                                      //
	println("Esta es la clave privada del usuario SIN DESENCRIPTAR : ", respuestaSiCreado.Msg) //
	print("\n")

	//6º  Desencriptar la clave privada del usuario que viene desde el servidor
	var pkClient rsa.PrivateKey

	privateKet := decompress(decrypt(decode64(respuestaSiCreado.Msg), u.KeyData)) //Está en bytes

	json.Unmarshal(privateKet, &pkClient)

	//7º Desencriptar la clave random con la clave privada del usuario
	claveRandomDesencriptada, err := rsa.DecryptPKCS1v15(rand.Reader, &pkClient, encClave)

	//8º Desencriptar el tema con la clave random
	temaDesencriptado := decompress(decrypt(decode64(temaEncriptado), claveRandomDesencriptada))

	//9º Printear el tema
	println(temaDesencriptado)

	/*
		//Encodea el tema de alguna manera rara
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

		//Creamos la clave random del tema encriptada con la clave pública del usuario que lo quiere leer
		var clavePubUnmarshal rsa.PublicKey
		json.Unmarshal([]byte(respuestatonta.Msg), &clavePubUnmarshal)
		encClave, err := rsa.EncryptPKCS1v15(rand.Reader, &clavePubUnmarshal, keyTema64)

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
		data2.Set("encriptado", string(encClave))                                   //Key random codificada con la clave publica del usuario
		data2.Set("file", encode64(encrypt(compress([]byte(temaJSON)), keyTema64))) //Tema creado encodeado

		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		//Simulación de lectura

		//Variable para chckear que funciona
		ficheroRandom := encode64(encrypt(compress([]byte(temaJSON)), keyTema64))

		println("estás detrás de fichero Random")

		r2, err := client.PostForm("https://localhost:10443", data2)
		chk(err)

		//Recibir la clave privada del usuario

		//Esto está bien
		respuestaSiCreado := Resp{}                                                                //Recibimos la clave privada del usuario CODIFICADA			//
		byteValue, _ = ioutil.ReadAll(r2.Body)                                                     //
		json.Unmarshal([]byte(byteValue), &respuestaSiCreado)                                      //
		println("Esta es la clave privada del usuario SIN DESENCRIPTAR : ", respuestaSiCreado.Msg) //
		print("\n")

		//1º Desencriptar la clave privada del usuario que recibimos del servidor
		privateKet := decompress(decrypt(decode64(respuestaSiCreado.Msg), u.KeyData)) //Está en bytes

		println(string(privateKet))

		privateKetMarshall, err := json.Marshal(privateKet)

		var clavePrivadaDesencriptada rsa.PrivateKey

		json.Unmarshal(privateKetMarshall, clavePrivadaDesencriptada) //Guardamos en formato rsa.privateKey //Esto de aquí está mal

		//2º Desencriptar la clave Random con la clave privada
		claveRandomDesencriptada, err := rsa.DecryptPKCS1v15(rand.Reader, &clavePrivadaDesencriptada, encClave)

		println(claveRandomDesencriptada)

		//3º Desencriptar el tema encriptado con la clave random
		temaDesencriptado := decompress(decrypt(decode64(ficheroRandom), claveRandomDesencriptada))

		//4º Printear el tema
		println(temaDesencriptado)

		//////////////

		var pk rsa.PrivateKey //Esto es una private KEY  pk = 0

		println(string(privateKet))

		//pk, err = x509.ParsePKCS1PrivateKey(privateKet)

		//Pasar los bytes de privateKet a RSA Key

		println("hola que tal")

		//var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(&pk) //Esto es la clave privada en bytes

		json.Unmarshal(privateKet, &pk)

		_ = ioutil.WriteFile("datos_"+nombreUsuario+"_"+t.Titulo+".json.enc", decode64(ficheroRandom), 0644)

		//fichetiro, err := ioutil.ReadFile("datos_" + nombreUsuario + "_" + t.Titulo + ".json.enc")
		chk(err)
		//ficheritoEnconded := encode64(fichetiro)

		//Igual hace falta un unmarshal

		//temaDesencriptado := decrypt(decode64(ficheritoEnconded), privateKet)

		var temaFinal tema

		json.Unmarshal(temaDesencriptado, &temaFinal)

		fmt.Println(temaDesencriptado)*/

	Opciones(resp)

}
