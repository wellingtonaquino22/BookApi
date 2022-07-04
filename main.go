package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http" // pacote de http
	"strconv"
	"strings"
)

type book struct {
	Id     int    `json:"id"`
	Titulo string `json:"title"`
	Autor  string `json:"actor"`
}

var books []book = []book{ //Lista de livros (Fatias)
	{
		Id:     1,
		Titulo: "Title1",
		Autor:  "Actor2",
	},
	{
		Id:     2,
		Titulo: "Title2",
		Autor:  "Actor2",
	},
	{
		Id:     3,
		Titulo: "Title3",
		Autor:  "Actor3",
	},
	{
		Id:     4,
		Titulo: "Title4",
		Autor:  "Actor4",
	},
	{
		Id:     5,
		Titulo: "Title5",
		Autor:  "Actor5",
	},
	{
		Id:     6,
		Titulo: "Title6",
		Autor:  "Actor6",
	},
}

func mainRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem Vindo")

}
func listBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	encoder := json.NewEncoder(w) // codificar em json passando w como argumento, Vai gerar um encoding
	encoder.Encode(books)         //passa o que quer que seja codificado no formato json, automaticamente envia resposta para o cliente.

}
func registerBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	//readall = retorna o corpo que foi decodificado para fatia de byte e um erro.
	body, erro := ioutil.ReadAll(r.Body) // readall leva como argumento um leitor read e responde com uma fatia de byte e segundo erro!, Ler o corpo da requisição (Request)
	if erro != nil {
		//lidar com erro
	}
	var newBook book
	json.Unmarshal(body, &newBook) //decodificar, segundo argumento é o que vc quer aplicar (converter)
	newBook.Id = len(books) + 1    // pega o comprimento e adiciona mais um no id (Organizar Id)
	books = append(books, newBook) // adiciona o novolivro no final da lista
	encoder := json.NewEncoder(w)
	encoder.Encode(newBook)
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	// delete /livros/123
	parts := strings.Split(r.URL.Path, "/") // extrae o ID
	id, erro := strconv.Atoi(parts[2])      // Converte de string para Int

	if erro != nil { // se der algum erro na conversão manda badrequest
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indiceBook := -1 // livro ainda não encontrado
	for i, book := range books {
		if book.Id == id { // se achar um livro cujo id foi o fornecido
			indiceBook = i // atribui o indice
			break
		}
	}
	// se não encontrou o livro, status not found
	if indiceBook < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leftside := books[0:indiceBook] // 0 1 2 3 4 5 6 7 8
	rightside := books[indiceBook+1 : len(books)]
	books = append(leftside, rightside...) // atualiza a lista
	w.WriteHeader(http.StatusNoContent)    // sem conteúdo na resposta de exclusão
}
func editBook(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id, erro := strconv.Atoi(parts[2])

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, bodyErro := ioutil.ReadAll(r.Body)

	if bodyErro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var modifiedBook book
	erroJson := json.Unmarshal(body, &modifiedBook)

	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indiceBook := -1

	for i, book := range books {
		if book.Id == id {
			indiceBook = i
			break
		}
	}
	if indiceBook < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	books[indiceBook] = modifiedBook

	json.NewEncoder(w).Encode(modifiedBook)
}

func routeBooks(w http.ResponseWriter, r *http.Request) {
	//livros
	// /livros/ ->  /livros/qualquercoisa
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 2 || len(parts) == 3 && parts[2] == "" {

		if r.Method == "GET" {
			listBooks(w, r)
		} else if r.Method == "POST" {
			registerBook(w, r)
		}
	} else if len(parts) == 3 || len(parts) == 4 && parts[3] == "" {
		if r.Method == "GET" {
			searchBooks(w, r)
		} else if r.Method == "DELETE" {
			deleteBook(w, r)
		} else if r.Method == "PUT" {
			editBook(w, r)
		}
		searchBooks(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
func searchBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path) //rota que esta batendo (requerida)
	//split = quebra a string em varias partes
	part := strings.Split(r.URL.Path, "/") //quebra a url.path (nesse caso caminho / books / alguma coisa)

	id, _ := strconv.Atoi(part[2]) //extraiu id e converteu uma string para um indice

	for _, book := range books { //loop para encontrar o livro
		if book.Id == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func configRoutes() { //Rotas
	http.HandleFunc("/", mainRoute) // 1 argumento local da rota (no caso a principal) , quando chegar entrar na rota, executa a função   )
	http.HandleFunc("/books", routeBooks)

	//e.g get /books/123
	http.HandleFunc("/books/", routeBooks) //Barra no final = qualquer coisa depois (Quantas quiser)
}

func serverConfig() { // w = Resposta , r = pedido
	configRoutes()
	fmt.Println("Server Rodando")
	log.Fatal(http.ListenAndServe(":3000", nil)) // Rodar o servidor, metodo 2 argumentos (endereço (Número da porta), servidor mux, o nil retorna o default servemux)

}

func main() {
	serverConfig()

}

/*
Obs: C:\Users\Wellington\go\bin\CompileDaemon        |Roda o servidor automaticamente
	C:\Users\Wellington\go\src\github.com\githubnemo\CompileDaemon     |Com o executável
*/
