package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func crear_vinculo_privado(w http.ResponseWriter, req *http.Request) {
	if req != nil {

		req.ParseForm()
		w.Header().Set("Content-Type", "text/plain") //Cabecera estandard

		usuario := req.Form.Get("user")
		tema := req.Form.Get("tema")
		clave := req.Form.Get("encriptado")

		filezilla := req.Form.Get("file") //Recuperamos el file

		_ = ioutil.WriteFile("datos_"+usuario+"_"+tema+".json.enc", decode64(filezilla), 0644)

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

		nameFicherito := "datos_" + usuario + "_" + tema + ".json.enc"

		lectura, err := ioutil.ReadFile(nameFicherito)
		chk(err)
		encode64(lectura)

		clavePrivadaUsuario := gUsers[usuario].Data["prikey"]

		response(w, true, clavePrivadaUsuario) //Devolvemos la clave privada del usuario que va a decriptar el tema

		//println(gUsers[usuario].Data["private"])
		//println(gUsers[usuario].Data["prikey"])
		//println(desencriptado)
	}
}
