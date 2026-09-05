# 🛒 SGE (Sistema de Gerenciamento de Estoque) - v2

Segunda versão da construção de um sistema de gerenciamento de estoque (SGE) desktop. 

A primeira tentativa utilizava Node.js no back-end, mas a arquitetura foi reformulada para utilizar **Go**, garantindo melhor desempenho e menor consumo de recursos para uma aplicação desktop.

---

## 🛠️ Tecnologias Utilizadas

* **Back-End:** [Go](https://go.dev/)
* **Front-End:** C# / .NET (sujeito a alterações ao longo do desenvolvimento)

---

## 📄 Documentação da API (End-points)

A API do Back-End roda por padrão no endereço: `http://localhost:8080`

| Método | Rota | Descrição |
| :--- | :--- | :--- |
| **POST** | `/add-product` | Cadastra um novo produto |
| **POST** | `/add-category` | Cadastra uma nova categoria |
| **POST** | `/add-user` | Cadastra um novo usuário |
| **PUT** | `/edit-product` | Atualiza os dados de um produto |
| **PUT** | `/edit-category` | Atualiza os dados de uma categoria |
| **PUT** | `/edit-user` | Atualiza os dados de um usuário |
| **PUT** | `/edit-password` | Altera a senha do usuário |
| **DELETE** | `/delete-product` | Remove um produto |
| **DELETE** | `/delete-category` | Remove uma categoria |
| **DELETE** | `/delete-user` | Remove um usuário |

---

## 🚀 Como Executar o Projeto

### Pré-requisitos
* Go (versão 1.22 ou superior)
* SDK do .NET (versão 8.0 ou superior)

### Passo a passo
1. Clone o repositório:
   git clone [https://github.com/Albinomavel/Projeto-PDV](https://github.com/Albinomavel/Projeto-PDV)```

2. Preste atenção nas tags, é uma boa prática conferi-las para garantir o bom funcionamento, nesse caso eu utilizei um protocolo que conversava com o meu próprio computador, caso mudar preste atenção

---

## 📌 Licenças e uso

- Esta é uma versão padronizada que pode ser alterada e personalizada conforme suas necessidades.
- Em caso de alteração ou reutilização do código, o autor dessa versão não se responsabiliza por modificações de terceiros ou inclusão de arquivos maliciosos.