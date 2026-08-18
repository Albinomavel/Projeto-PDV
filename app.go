package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte("sua_chave_secreta_aqui")

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
	http.HandleFunc("/add-category", middleware("operador", AddCategory))
	http.HandleFunc("/add-user", AddUser)
	http.HandleFunc("/login", Login)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}

}

// BLOCO DE STRUCTS PARA DECODIFICAR O JSON
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
type usuarios struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	HashPassword string `json:"HashPassword"`
	Cargo        string `json:"cargo"`
	CompleteName string `json:"completeName"`
	Email        string `json:"email"`
	Telefone     string `json:"telefone"`
	Ativo        bool   `json:"ativo"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
type ErrorResponse struct {
	Error  string `json:"error"`
	Codigo int    `json:"codigo"`
}
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Cargo    string `json:"cargo"`
	jwt.RegisteredClaims
}
type RefreshClaims struct {
	UsuarioID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// FIM DO BLOCO DE STRUCTS

// BLOCO DE FUNÇÕES PARA INSERIR DADOS NO BANCO DE DADOS
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
func AddUser(w http.ResponseWriter, r *http.Request) {
	var user usuarios
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao decodificar o usuário", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user.HashPassword) == "" {
		http.Error(w, "A senha é obrigatória", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao gerar oa senha", http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO usuarios (username, senha_hash, cargo, nome_completo, email, numero_telefone, Ativo) VALUES (?, ?, ?, ?, ?, ?, ?)"
	res, err := db.Exec(query, user.Username, hash, user.Cargo, user.CompleteName, user.Email, user.Telefone, user.Ativo)
	if err != nil {
		http.Error(w, "Erro ao inserir o usuário no banco de dados", http.StatusInternalServerError)
		return
	}
	user.HashPassword = ""
	resultID, _ := res.LastInsertId()
	user.ID = int(resultID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	fmt.Println("Usuário inserido com sucesso!")
}

// FIM DO BLOCO DE FUNÇÕES PARA INSERIR DADOS NO BANCO DE DADOS
func getSecret() ([]byte, error) {
	secret := jwtKey
	if secret == nil {
		return nil, ErrSecretAysento
	}
	return []byte(secret), nil
}
func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição de login", http.StatusBadRequest)
		return
	}
	var user usuarios
	query := "SELECT id, username, senha_hash, cargo, ativo FROM usuarios WHERE username = ?"

	err = db.QueryRow(query, loginReq.Username).Scan(&user.ID, &user.Username, &user.HashPassword, &user.Cargo, &user.Ativo)
	if err != nil {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}
	resp, err := generateToken(user.ID, user.Username, user.Cargo)
	if err != nil {
		http.Error(w, "Erro ao gerar o token de acesso", http.StatusInternalServerError)
		return
	}
	tokenResponse := TokenResponse{
		AccessToken: resp,
		TokenType:   "Bearer",
		ExpiresIn:   int64(TokenExpirationTime.Seconds()),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse)

	fmt.Println("Login bem-sucedido!")
}

// BLOCO DE SEGURANÇA
func generateToken(userID int, username, cargo string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}
	agora := time.Now()

	claims := Claims{
		UserID:   userID,
		Username: username,
		Cargo:    cargo,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "app.go",
			IssuedAt:  jwt.NewNumericDate(agora),
			ExpiresAt: jwt.NewNumericDate(agora.Add(TokenExpirationTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func validateToken(tokenString string) (*Claims, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}
	Claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return Claims, nil
}
func middleware(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len("Bearer "):]
		claims, err := validateToken(tokenString)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				http.Error(w, "Token expirado", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		if claims.Cargo != role {
			http.Error(w, "Acesso negado: permissão insuficiente", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

const (
	TokenExpirationTime        = time.Minute * 15
	RefreshTokenExpirationTime = time.Hour * 24 * 7
)

var (
	ErrTokenExpired  = errors.New("token expirado")
	ErrTokenInvalid  = errors.New("token inválido")
	ErrSecretAysento = errors.New("chave secreta ausente")
)

// FIM DO BLOCO DE SEGURANÇA
