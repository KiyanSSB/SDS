package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type User struct {
	Nombre    string
	Apellidos string
	Email     string
	Password  string
}

type Comentario struct {
	Usuario    User
	Comentario string
}

type Tema struct {
	Titulo      string
	Descripcion string
	Creador     User
	Comentarios []Comentario
}

func escribir_temas() {
	tema := Tema{
		Titulo:      "Programación 1",
		Descripcion: "La asignatura Programación 1 es el primer contacto del estudiante con la programación de ordenadores.",
		Creador: User{
			Nombre:    "Gema",
			Apellidos: "Lozano Jimenez",
			Email:     "glj4@alu.ua.es",
			Password:  "12345",
		},
		Comentarios: []Comentario{
			Comentario{
				Usuario: User{
					Nombre:    "Gema",
					Apellidos: "Lozano Jimenez",
					Email:     "glj4@alu.ua.es",
					Password:  "12345",
				},
				Comentario: "Mu weno",
			},
			Comentario{
				Usuario: User{
					Nombre:    "Antero",
					Apellidos: "Guarinos Caballero",
					Email:     "agc8@alu.ua.es",
					Password:  "12345",
				},
				Comentario: "Spesial",
			},
		},
	}

	file, _ := json.MarshalIndent(tema, "", " ")
	_ = ioutil.WriteFile("temas.json", file, 0644)
}

func leer_temas() {
	manejadorDeArchivo, err := ioutil.ReadFile("temas.json")
	if err != nil {
		log.Fatal(err)
	}
	c := Tema{}
	err = json.Unmarshal(manejadorDeArchivo, &c)

	if c.Titulo == "Programación 1" {
		fmt.Println("Título correcto")
	}
}

func main() {
	escribir_temas()
	leer_temas()
}
