package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/PROJETO_PDV")
	if err != nil {
		panic(err)
	}
	route := "http://localhost:8080"

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println(`a rota é ` + route)

	http.HandleFunc("/add-product", AddProduct)
	http.HandleFunc("/add-category", AddCategory)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}

}

type produto struct {
	ID           int     `json:"id"`
	Nome         string  `json:"nome"`
	Descricao    string  `json:"descricao"`
	CategoriaID  int     `json:"categoria_id"`
	PrecoVenda   float64 `json:"preco_venda"`
	Estoque      int     `json:"estoque"`
	PrecoCompra  float64 `json:"preco_compra"`
	CodigoBarras int     `json:"codigo_barra"`
}
type categoria struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product produto
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Erro ao decodificar o produto", http.StatusBadRequest)
		return
	}
	query := "INSERT INTO produtos (nome, descricao, categoria_id, preco_venda, estoque, preco_compra, codigo_barra) VALUES (?, ?, ?, ?, ?, ?, ?)"
	res, err := db.Exec(query, product.Nome, product.Descricao, product.CategoriaID, product.PrecoVenda, product.Estoque, product.PrecoCompra, product.CodigoBarras)
	if err != nil {
		http.Error(w, "Erro ao inserir o produto no banco de dados", http.StatusInternalServerError)
		return
	}
	resultID, _ := res.LastInsertId()
	product.ID = int(resultID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)

	fmt.Println("Produto inserido com sucesso!")
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	var category categoria
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Erro ao decodificar a categoria", http.StatusBadRequest)
		return
	}
	query := "INSERT INTO categoria (nome) VALUES (?)"
	res, err := db.Exec(query, category.Nome)
	if err != nil {
		http.Error(w, "Erro ao inserir a categoria no banco de dados", http.StatusInternalServerError)
		return
	}
	resultID, _ := res.LastInsertId()
	category.ID = int(resultID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

	fmt.Println("Categoria inserida com sucesso!")
}
