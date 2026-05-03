CREATE TABLE categoria (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
nome VARCHAR(255) NOT NULL
);
CREATE TABLE produtos (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
nome VARCHAR(255) NOT NULL,
descricao TEXT,
categoria_id INT NOT NULL,
preco_venda DECIMAL(10, 2) NOT NULL,
estoque INT NOT NULL comment 'Quantidade baseada em caixa não em unidades',
preco_compra DECIMAL(10, 2) NOT NULL,
codigo_barra VARCHAR(255) NOT NULL,
FOREIGN KEY (categoria_id) REFERENCES categoria(id)
);
CREATE TABLE usuarios (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
username VARCHAR(255) NOT NULL,
senha_hash VARCHAR(255) NOT NULL,
cargo VARCHAR(255) NOT NULL,
nome_completo VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL,
numero_telefone VARCHAR(14) NOT NULL comment 'Formato: +55 (XX) XXXXX-XXXX'
);
CREATE TABLE vendas (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
data_venda DATETIME NOT NULL,
total DECIMAL(10, 2) NOT NULL,
usuario_id INT NOT NULL,
FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);
CREATE TABLE item_venda (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
venda_id INT NOT NULL,
produto_id INT NOT NULL,
quantidade INT NOT NULL,
preco_unitario DECIMAL(10, 2) NOT NULL,
FOREIGN KEY (venda_id) REFERENCES vendas(id),
FOREIGN KEY (produto_id) REFERENCES produtos(id)
);
CREATE TABLE relatorio_vendas (
id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
data_inicio DATE NOT NULL,
data_fim DATE NOT NULL, total_vendas DECIMAL(10, 2) NOT NULL,
total_itens_vendidos INT NOT NULL
);