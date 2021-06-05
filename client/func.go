package main

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func leerTerminal() string {
	text := bufio.NewReader(os.Stdin)
	read, _ := text.ReadString('\n')
	tipo := strings.TrimRight(read, "\r\n")
	return tipo
}

func StringAInt(s string) int64 {
	number, _ := strconv.ParseInt(s, 10, 0)
	return number
}

// respuesta del servidor
type Resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

// respuesta del servidor
type Respfix struct {
	Ok   bool // true -> correcto, false -> error
	File string
	Msg  string // mensaje adicional
}

type Comentario struct {
	Usuario    string
	Comentario string
}

type registryTema struct {
	Key   []byte
	Temas map[string]tema
}

var gComentarios map[string]Comentario

type tema struct {
	Titulo      string
	Descripcion string
	Comentarios map[string]Comentario
}

var gTemas map[string]tema

type data struct {
	Tema map[string]tema
}

var gData map[string]data

//Cliente global
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client = &http.Client{Transport: tr}
