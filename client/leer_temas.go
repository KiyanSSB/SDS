package main

import (
	"fmt"
)

func leer_temas(resp Resp) {
	if gTemas == nil {
		gTemas = make(map[string]tema)
	}

	for k := range gTemas {
		fmt.Println("Titulo: " + gTemas[k].Titulo)
		fmt.Println("Descripcion: " + gTemas[k].Descripcion)
		fmt.Println("")
	}
}
