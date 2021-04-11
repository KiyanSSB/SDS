package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"golang.org/x/crypto/scrypt"
)

/******************************************************************************************************************************
***					En este fichero se encuentra todo lo relacionado con el registro en el apartado del servidor			***
*******************************************************************************************************************************/

func signup(w http.ResponseWriter, req *http.Request) {
	u := user{}
	u.Name = req.Form.Get("user")              // nombre
	u.Salt = make([]byte, 16)                  // sal (16 bytes == 128 bits)
	rand.Read(u.Salt)                          // la sal es aleatoria
	u.Data = make(map[string]string)           // reservamos mapa de datos de usuario
	password := decode64(req.Form.Get("pass")) // contraseña (keyLogin)

	// "hasheamos" la contraseña con scrypt
	u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)

	_, ok := gUsers[u.Name] // ¿existe ya el usuario?
	if ok {
		response(w, false, "Usuario ya registrado")
	} else {
		gUsers[u.Name] = u
		almacenarArchivo()
		fmt.Println("Almacenado")
		response(w, true, "Usuario registrado")
	}
}
