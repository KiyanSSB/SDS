package main

import "fmt"

func escribir_comentarios(resp Resp) {
	fmt.Println("------------------------------------")
	fmt.Println("Que tema quieres comentar")
	fmt.Println("------------------------------------")

	for k := range gTemas {
		fmt.Println("Titulo: " + gTemas[k].Titulo)

		fmt.Println("Descripcion: " + gTemas[k].Descripcion)
		fmt.Println("Comentarios")
		for i := range gTemas[k].Comentarios {
			fmt.Println("    nombre de ususario :" + gTemas[k].Comentarios[i].NombreUsuario)
			fmt.Println("    comentario :" + gTemas[k].Comentarios[i].Comentario)
		}

	}
	nombreTema := leerTerminal()

	for k := range gTemas {
		if gTemas[k].Titulo == nombreTema {
			fmt.Println("introduce un comentario")
			comentarioTerminal := leerTerminal()

			c := comentario{}

			c.NombreUsuario = u.Name

			c.Comentario = comentarioTerminal

			gComentarios[c.NombreUsuario] = c

			gTemas[k].Comentarios = gComentarios

		}

	}

}
