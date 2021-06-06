package main

/**********************************************************************************************************************
 ***			En este fichero se encuentra todo lo relacionado con el registro en el apartado del cliente				***
 ***********************************************************************************************************************/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var caracteresInvalidos = map[int]string{0: "!", 1: "\"", 2: "#", 3: "$", 4: "%", 5: "&", 6: "(", 7: ")",
	8: "*", 9: "+", 10: ",", 11: "-", 12: ".", 13: "/", 14: ":", 15: ";", 16: "<", 17: "=",
	18: ">", 19: "?", 20: "@", 21: "[", 22: "\\", 23: "]", 24: "_", 25: "{", 26: "|", 27: "}",
	28: "á", 29: "Á", 30: "é", 31: "É", 32: "í", 33: "Í", 34: "ó", 35: "Ó", 36: "ú", 37: "Ú",
	38: "à", 39: "À", 40: "è", 41: "È", 42: "ì", 43: "Ì", 44: "ò", 45: "Ò", 46: "ù", 47: "Ù",
	48: "ä", 49: "Ä", 50: "ë", 51: "Ë", 52: "ï", 53: "Ï", 54: "ö", 55: "Ö", 56: "ü", 57: "Ü",
	58: "'", 59: "^", 60: "¬", 61: "·"}

func signup(client *http.Client, cmd string) {
	fmt.Println("Registrar un usuario")
	fmt.Println("--------------------")

	nombreCorrecto := false
	fmt.Print("Nombre de usuario: ")
	user := leerTerminal()

	for !nombreCorrecto {
		nombreCorrecto = true
		for i := 0; i < 62; i++ {
			if strings.Contains(user, caracteresInvalidos[i]) {
				nombreCorrecto = false
				i = 38
				fmt.Println("\nEl nombre de usuario contiene caracteres inválidos")
				fmt.Println("Por favor, repita el nombre de usuario: ")
				user = leerTerminal()
			}
		}
	}

	fmt.Print("Contraseña: ")
	pass := leerTerminal()

	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := keyClient[:32]
	keyData := keyClient[32:64] //La otra para los datos

	//Generamos un par de claves (privada, pública) para el servidor
	pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
	chk(err)
	pkClient.Precompute()

	pkJSON, err := json.Marshal(&pkClient)
	chk(err)

	keyPub := pkClient.Public()
	pubJSON, err := json.Marshal(&keyPub)
	chk(err)

	data := url.Values{}
	data.Set("cmd", cmd)
	data.Set("user", user)
	data.Set("pass", encode64(keyLogin))

	//Comprimimos y codificamos la clave pública
	data.Set("pubkey", encode64(compress(pubJSON)))

	//Comprimimos ciframos y codificamos la clave privada
	data.Set("prikey", encode64(encrypt(compress(pkJSON), keyData)))

	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	resp := Resp{}
	byteValue, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(byteValue), &resp)
	Opciones(resp)
}
