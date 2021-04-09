package main

import (
	"fmt"
	"os"
)

/**************************************************************************************************************
***			En este fichero se encuentra todo lo relacionado con las funcionalidades del cliente			***
***************************************************************************************************************/

/**************************************************
***		ESTRUCTURAS DE DATOS DEL CLIENTE		***
***************************************************/

/******************************************
***		FUNCIONES BASE DEL CLIENTE		***
*******************************************/

/**********************************
***		MAIN DEL CLIENTE		***
***********************************/

//Seleccion de inicio sesion o registrarse
func main() {

	var args struct {
		Operation string `arg:"positional, required" help:"(signup|signin)"`
	}

	fmt.Println("*********************************************************************************")
	fmt.Println("***	Bienvenido al sistema de Foros de la asignatura de SDS en 20/21		******")
	fmt.Println("*********************************************************************************")

	parser := arg.MustParse(&args)

	switch args.Operation {
	case "signup":
		signup(client, "register")
	case "signin":
		signin(client, "login")
	case "help":
		parser.WriteHelp(os.Stdin)
	default:
		parser.Fail(args.Operation)
	}
}
