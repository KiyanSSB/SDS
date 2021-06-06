package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/**************************************************************************************************************************
***			En este fichero se encuentra todo lo relacionado con la ejecución de las funcionalidades del servidor		***
***************************************************************************************************************************/

/**************************************************
***		ESTRUCTURAS DE DATOS DEL SERVIDOR		***
***************************************************/

type user struct {
	Name string            // nombre de usuario
	Hash []byte            // hash de la contraseña
	Salt []byte            // sal para la contraseña
	Data map[string]string // datos adicionales del usuario
}

type Comentario struct {
	Usuario    string
	Comentario string
}

type tema struct {
	Titulo      string
	Descripcion string
	Comentarios map[string]Comentario
}

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

type respf struct {
	Ok   bool // true -> correcto, false -> error
	File string
	Msg  string // mensaje adicional
}

type registry struct {
	Key   []byte
	Users map[string]user
}

type registryTema struct {
	Key   []byte
	Temas map[string]tema
}

var gUsers map[string]user //Mapa que contiene un mapa donde un string coincide con un usuario
var gTemas map[string]tema //Mapa que contiene un mapa donde un string coincide con un tema
var codee []byte

/******************************************
***		FUNCIONES BASE DEL SERVIDOR		***
*******************************************/

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func leerTerminal() string {
	text := bufio.NewReader(os.Stdin)
	read, _ := text.ReadString('\n')
	tipo := strings.TrimRight(read, "\r\n")
	return tipo
}

// función para escribir una respuesta del servidor
func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

func responsefix(w io.Writer, file string, ok bool, msg string) {
	r := respf{Ok: ok, File: file, Msg: msg}
	rJSON, err := json.Marshal(&r)
	chk(err)
	w.Write(rJSON)
}

//Función que maneja las diferentes opciones del servidor cuando recibe una petición del cliente
func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              //Es necesario parsear el usuario
	w.Header().Set("Content-Type", "text/plain") //Cabecera estandar

	switch req.Form.Get("cmd") { //Comprobamos el comando desde el cliente
	case "signup":
		fmt.Println("El cliente ha seleccionado REGISTRO")
		signup(w, req)
	case "signin":
		fmt.Println("El cliente ha seleccionado LOGIN")
		signin(w, req)
	case "crear_tema":
		fmt.Println("Enviar datos")
		crear_tema(w, req)
	case "publicos":
		fmt.Println("El cliente ha seleccionado ver los temas PUBLICOS")
		enviar_temas_publicos(w, req)
	case "comentario":
		fmt.Println("El cliente ha seleccionado AÑADIR UN COMENTARIO")
		crear_comentario(w, req)
	default:
		response(w, false, "Comando inválido")
	}
}

func abrirArchivo() {
	file, err := os.Open("registro.json") // abrimos el primer fichero (entrada)

	gUsers = make(map[string]user) //Inicializamos el mapa de los usuarios

	if err != nil {
		file, err = os.Create("registro.json") // abrimos el segundo fichero (salida)
		// inicializamos mapa de usuarios
		if err != nil {
			panic(err)
		}
	} else {
		defer file.Close() //Por último cerramos el fichero

		byteValue, _ := ioutil.ReadAll(file) //Guardamos el contenido del fichero entero en la variable
		var code []byte = nil                //Creamos una variable de bytes
		Regi := registry{Key: code, Users: gUsers}

		/*Utilizamoe el contenido del fichero y lo desencriptamos con el valor codee, que es la primera parte de la contraseña del servidor
		Luego, guardamos en la variable Regi el contenido del json que almacena el servidor*/
		json.Unmarshal(decrypt(byteValue, codee), &Regi)
		//fmt.Println(Regi)

		/*Comprobamos si la contraseña introducida por la persona que inicializa el servidor es la correcta si coincide con la misma que se ha utilizado para
		Codificar los datos del json, si coinciden el codee y la Key del registro, podemos continuar */
		verdad := bytes.Equal(Regi.Key, codee)
		if verdad == false {
			fmt.Println("La contraseña no es la correcta")
			panic(err)
		} else {
			fmt.Println("Esto es el valor del Regi.key: ", Regi.Key)
			fmt.Println("Esto es el valor del codee:", codee)
			fmt.Println("La contraseña es correcta, puedes continuar")
		}
	}

	/* Cargamos el apartado de los temas del servidor */
	fichero, error := os.Open("temas.json")
	gTemas = make(map[string]tema)

	if error != nil {
		fichero, error = os.Create("temas.json")
		if error != nil {
			panic(error)
		}
	} else {
		defer fichero.Close()
		contenidoFichero, _ := ioutil.ReadAll(fichero)
		var codigo []byte = []byte("1")
		Regi2 := registryTema{Key: codigo, Temas: gTemas}
		json.Unmarshal(decrypt(contenidoFichero, codee), &Regi2)
		fmt.Println(Regi2)
	}
}

func almacenarArchivo() {
	var code []byte = nil
	Regi := registry{Key: code, Users: gUsers}
	Regi.Key = codee
	Regi.Users = gUsers
	os.Remove("registro.json")
	_, err := os.Create("registro.json")
	chk(err)
	jsonFD, err := json.Marshal(&Regi) //Recordar cambiar los nombres
	//jsonFD := encrypt(jsonF, codee)
	err = ioutil.WriteFile("registro.json", jsonFD, 0644)
	chk(err)
}

//Almacenamos los temas en un fichero borrando el archivo anterior y actualizandolo
func almacenarTema() {
	var code []byte = nil
	Tem := registryTema{Key: code, Temas: gTemas}
	Tem.Key = codee
	Tem.Temas = gTemas
	os.Remove("temas.json")
	_, err := os.Create("temas.json")
	chk(err)
	jsonF, err := json.Marshal(&Tem)

	//Encriptamos el json de los temas con el codigo de la contraseña del server
	//jsonFD := encrypt(jsonF, codee)
	var jsonFD = jsonF

	err = ioutil.WriteFile("temas.json", jsonFD, 0644)
	chk(err)
}

/**********************************
***		MAIN DEL PROGRAMA		***
***********************************/

func main() {
	fmt.Println("Bienvenido al sistema de foros de SDS 20/21")
	fmt.Println("Te encuentras ejecutando el SERVIDOR")
	fmt.Println("------------------------------------")
	fmt.Print("Dime la contraseña del servidor: ")
	key := leerTerminal()
	data := sha512.Sum512([]byte(key))
	fmt.Println("El valor de data es el siguiente: ", data)
	codee = data[:32] //El codigo es los primeros 32
	abrirArchivo()
	http.HandleFunc("/", handler)
	chk(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}
