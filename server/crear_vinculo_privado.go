package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func crear_vinculo_privado(w http.ResponseWriter, req *http.Request) {
	if req != nil {

		usuario := req.Form.Get("user")
		tema := req.Form.Get("tema")
		clave := req.Form.Get("encriptado")

		println(usuario)
		println(tema)
		println(clave)

		file, err := os.Open(usuario + "_" + tema + "_priv.json")

		if err != nil {
			file, err = os.Create(usuario + "_" + tema + "_priv.json")
			if err != nil {
				panic(err)
			}
		}

		defer file.Close()
		jsonFD, err := json.Marshal(&clave)
		err = ioutil.WriteFile(usuario+"_"+tema+"_priv.json", jsonFD, 0644)

		//println(gUsers[usuario].Data["private"])
		//println(gUsers[usuario].Data["prikey"])
		//println(desencriptado)
	}
}
