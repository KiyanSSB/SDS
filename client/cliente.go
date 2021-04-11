package main

import (
	"fmt"
	"os"

	arg "github.com/alexflint/go-arg"
)

/**************************************************************************************************************
***			En este fichero se encuentra todo lo relacionado con las funcionalidades del cliente			***
***************************************************************************************************************/

func Opciones(resp Resp) {
	fmt.Println("")
	fmt.Println("-" + resp.Msg + "-")

	if !resp.Ok {
		fmt.Println("Salir")
		return
	} else {
		if resp.Msg == "Usuario registrado" {
			fmt.Println("")
			fmt.Println("Inicia sesión")
			signin(client, "signin")
		} else if resp.Msg == "Credenciales válidas" {
			fmt.Println("---- MENÚ PRINCIPAL ----")
			fmt.Println("1. Crear un tema")
			fmt.Println("2. Ver los temas disponibles")
			fmt.Println("3. Añadir un comentario")
			fmt.Println("4. Cerrar el programa")
			fmt.Println("------------------------")
			fmt.Print("¿Qué opción quieres realizar? ")
			number := StringAInt(leerTerminal())

			switch number {
			case 1:
				guardar_tema("crear_tema", resp)
				return
			case 2:
				return
			case 3:
				return
			case 4:
				return
			default:
				Opciones(resp)
				return
			}
		}
	}

}

/**********************************
***		MAIN DEL CLIENTE		***
***********************************/

//Seleccion de inicio sesion o registrarse
func main() {

	var args struct {
		Operation string `arg:"positional, required" help:"(signup|signin)"`
	}

	fmt.Println("***********************************************************************")
	fmt.Println("*** Bienvenido al sistema de Foros de la asignatura de SDS en 20/21 ***")
	fmt.Println("***********************************************************************")

	parser := arg.MustParse(&args)

	switch args.Operation {
	case "signup":
		signup(client, "signup")
	case "signin":
		signin(client, "signin")
	case "help":
		parser.WriteHelp(os.Stdin)
	default:
		parser.Fail(args.Operation)
	}
}
