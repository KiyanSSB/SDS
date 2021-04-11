package main

import (
	"bytes"
	"net/http"

	"golang.org/x/crypto/scrypt"
)

func signin(w http.ResponseWriter, req *http.Request) {
	u, ok := gUsers[req.Form.Get("user")] // ¿existe ya el usuario?
	if !ok {
		response(w, false, "Usuario inexistente")
		return
	}

	password := decode64(req.Form.Get("pass"))               // obtenemos la contraseña
	hash, _ := scrypt.Key(password, u.Salt, 16384, 8, 1, 32) // scrypt(contraseña)

	if bytes.Compare(u.Hash, hash) != 0 {
		response(w, false, "Credenciales inválidas")
		return
	}
	response(w, true, "Credenciales válidas")
}
