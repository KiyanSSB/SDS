package main

import (
	"fmt"
	"net/http"
)

func obtener_pubkey(w http.ResponseWriter, req *http.Request) {
	u := user{}
	u.Name = req.Form.Get("user")

	publicKey := gUsers[u.Name].Data["pubkey"]

	fmt.Println(publicKey)

	response(w, true, publicKey)
}
