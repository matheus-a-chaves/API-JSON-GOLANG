package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Livro struct {
	Id     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

var Livros []Livro = []Livro{
	Livro{
		Id:     1,
		Titulo: "O guarani",
		Autor:  "José de Alencar",
	},
	Livro{
		Id:     2,
		Titulo: "Cazuza",
		Autor:  "Virato Correia",
	},
	Livro{
		Id:     3,
		Titulo: "Dom Casmurro",
		Autor:  "Machado de Assis",
	},
}

func rotaPrincipal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem vindo")
}

func excluirLivro(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(partes[2])
	if err != nil { //tratando o erro
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err) //Mostra o erro no terminal
		return
	}
	indiceDoLivro := -1
	for indice, livro := range Livros {
		if livro.Id == id {
			indiceDoLivro = indice
			break
		}
	}
	if indiceDoLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Atualizando a lista (retirando o excluido)
	ladoEsquerdo := Livros[0:indiceDoLivro]
	ladoDireito := Livros[indiceDoLivro+1 : len(Livros)]
	Livros = append(ladoEsquerdo, ladoDireito...)

	w.WriteHeader(http.StatusNoContent)
}

func modificarLivro(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(partes[2])

	corpo, erroCorpo := ioutil.ReadAll(r.Body)
	if erroCorpo != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	var livroModificado Livro
	erroJson := json.Unmarshal(corpo, &livroModificado)

	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err != nil { //tratando o erro
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err) //Mostra o erro no terminal
		return
	}
	indiceDoLivro := -1
	for indice, livro := range Livros {
		if livro.Id == id {
			indiceDoLivro = indice
			break
		}
	}
	if indiceDoLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Mudando dados dos livros
	Livros[indiceDoLivro] = livroModificado
	json.NewEncoder(w).Encode(livroModificado)
	w.WriteHeader(http.StatusNoContent)

}

func rotearLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Formatando JsonView na Web
	// /livros
	// /livros/ -> /livros/id
	partes := strings.Split(r.URL.Path, "/")
	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		if r.Method == "GET" {
			listarLivros(w, r)
		} else if r.Method == "POST" {
			cadastrarLivros(w, r)
		}
	} else if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		if r.Method == "GET" {
			buscarLivros(w, r)
		} else if r.Method == "DELETE" {
			excluirLivro(w, r)
		} else if r.Method == "PUT" {
			modificarLivro(w, r)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func listarLivros(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(Livros)

}
func cadastrarLivros(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated) // Mudando status de resposta

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		//Lidar com erro
	} else {
		var novoLivro Livro
		json.Unmarshal(body, &novoLivro)
		novoLivro.Id = len(Livros) + 1     //Definindo um Id
		Livros = append(Livros, novoLivro) // Adicionando novo livro
		encoder := json.NewEncoder(w)
		encoder.Encode(novoLivro)
	}
}

func buscarLivros(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	//Quebrando a URL
	partes := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(partes[2])
	if err != nil { //tratando o erro
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err) //Mostra o erro no terminal
	}
	for _, livro := range Livros {
		if livro.Id == id {
			json.NewEncoder(w).Encode(livro)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func configurarRotas() {
	http.HandleFunc("/", rotaPrincipal)
	http.HandleFunc("/livros", rotearLivros)

	//e.g. GET /livros/123 -> seria o id pelo GET
	http.HandleFunc("/livros/", rotearLivros)
}

func configurarServidor() {
	configurarRotas()
	fmt.Println("O servidor está rodando")
	log.Fatal(http.ListenAndServe(":1337", nil))
}

func main() {
	configurarServidor()
}
