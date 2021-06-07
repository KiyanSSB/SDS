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
	if !resp.Ok {
		fmt.Println("Salir")
		return
	} else {
		if resp.Msg == "Usuario registrado" {
			fmt.Println("")
			fmt.Println("Inicia sesión")
			signin(client, "signin")
		} else if resp.Msg == "Credenciales válidas" || resp.Msg == "Añadido a la base de datos" {
			fmt.Println("")
			fmt.Println("---- MENÚ PRINCIPAL ----")
			fmt.Println("1. Crear un tema PUBLICO")
			fmt.Println("2. Crear un tema PRIVADO")
			fmt.Println("3. Ver todos los temas PUBLICOS")
			fmt.Println("4. Ver todos los temas PRIVADOS")
			fmt.Println("5. Añadir un comentario")
			fmt.Println("6. Cerrar el programa")
			fmt.Println("------------------------")
			fmt.Print("¿Qué opción quieres realizar? ")
			number := StringAInt(leerTerminal())
			fmt.Println("")

			switch number {
			case 1:
				guardar_tema("crear_tema", resp)
				return
			case 2:
				crear_tema_privado(resp)
			case 3:
				leer_temas_publicos(resp)
				return
			case 4:
				leer_temas_privados(resp)
			case 5:
				comentario("comentario", resp)
				return
			case 6:
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
