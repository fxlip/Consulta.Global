-- 01-create-schema.sql
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    cpf_cnpj VARCHAR(14) NOT NULL,
    tipo_pessoa CHAR(2) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    operadora VARCHAR(50) NOT NULL,
    UNIQUE (cpf_cnpj, operadora)
);

CREATE TABLE telefones (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
    ddd CHAR(2) NOT NULL,
    numero VARCHAR(9) NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    UNIQUE (ddd, numero)
);

CREATE TABLE enderecos (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
    logradouro VARCHAR(255) NOT NULL,
    numero_endereco VARCHAR(10),
    complemento VARCHAR(100),
    bairro VARCHAR(100),
    cidade VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL,
    cep CHAR(8) NOT NULL,
    tipo_endereco VARCHAR(50),
    UNIQUE (usuario_id, logradouro, numero_endereco)
);

CREATE TABLE contatos_adicionais (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
    tipo VARCHAR(50) NOT NULL,
    valor VARCHAR(255) NOT NULL,
    observacao TEXT,
    UNIQUE (usuario_id, tipo, valor)
);
