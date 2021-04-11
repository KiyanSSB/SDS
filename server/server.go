package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
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

type comentario struct {
	Usuario    user
	Comentario string
}

type tema struct {
	Titulo      string
	Descripcion string
	//Creador     user
	//Comentarios []comentario
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

var gUsers map[string]user
var gTemas map[string]tema
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

// función para cifrar (con AES en este caso), adjunta el IV al principio
func encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    // reservamos espacio para el IV al principio
	rand.Read(out[:16])                 // generamos el IV
	blk, err := aes.NewCipher(key)      // cifrador en bloque (AES), usa key
	chk(err)                            // comprobamos el error
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// función para descifrar (con AES en este caso)
func decrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16)     // la salida no va a tener el IV
	blk, err := aes.NewCipher(key)       // cifrador en bloque (AES), usa key
	chk(err)                             // comprobamos el error
	ctr := cipher.NewCTR(blk, data[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out, data[16:])     // desciframos (doble cifrado) los datos
	return
}

// función para comprimir
func compress(data []byte) []byte {
	var b bytes.Buffer      // b contendrá los datos comprimidos (tamaño variable)
	w := zlib.NewWriter(&b) // escritor que comprime sobre b
	w.Write(data)           // escribimos los datos
	w.Close()               // cerramos el escritor (buffering)
	return b.Bytes()        // devolvemos los datos comprimidos
}

// función para descomprimir
func decompress(data []byte) []byte {
	var b bytes.Buffer // b contendrá los datos descomprimidos

	r, err := zlib.NewReader(bytes.NewReader(data)) // lector descomprime al leer

	chk(err)         // comprobamos el error
	io.Copy(&b, r)   // copiamos del descompresor (r) al buffer (b)
	r.Close()        // cerramos el lector (buffering)
	return b.Bytes() // devolvemos los datos descomprimidos
}

// función para codificar de []bytes a string (Base64)
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// función para decodificar de string a []bytes (Base64)
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
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

//Función que maneja las diferentes opciones del servidor
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
	default:
		response(w, false, "Comando inválido")
	}
}

func abrirArchivo() {
	file, err := os.Open("registro.json") // abrimos el primer fichero (entrada)
	gUsers = make(map[string]user)

	if err != nil {
		file, err = os.Create("registro.json") // abrimos el segundo fichero (salida)
		// inicializamos mapa de usuarios
		if err != nil {
			panic(err)
		}
	} else {
		defer file.Close()
		byteValue, _ := ioutil.ReadAll(file)
		var code []byte = nil
		Regi := registry{Key: code, Users: gUsers}
		json.Unmarshal(decrypt(byteValue, codee), &Regi)
		verdad := bytes.Equal(Regi.Key, codee)
		if verdad == false {
			fmt.Println("La contraseña no es la correcta")
			panic(err)
		} else {
			Regi.Key = codee
			fmt.Println("La contraseña es correcta, puedes continuar")
		}
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
	jsonF, err := json.Marshal(&Regi)
	jsonFD := encrypt(jsonF, codee)
	err = ioutil.WriteFile("registro.json", jsonFD, 0644)
	chk(err)
}

func almacenarTema() {
	var code []byte = nil
	Tem := registryTema{Key: code, Temas: gTemas}
	Tem.Key = codee
	Tem.Temas = gTemas
	os.Remove("temas.json")
	_, err := os.Create("temas.json")
	chk(err)
	jsonF, err := json.Marshal(&Tem)
	jsonFD := encrypt(jsonF, codee)
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
	codee = data[:32]
	abrirArchivo()
	http.HandleFunc("/", handler)
	chk(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}
