# Consulta Global

Desenvolvido em [Go](https://golang.org/) para trabalhar com [PostgreSQL](https://postgresql.org/), utilizando [Redis](https://redis.io/) e [Elasticsearch](https://elastic.co/), e executado com [Docker](https://docker.com/), esse projeto tem o objetivo de fornecer uma aplicação rápida e moderna para consulta de dados.

## Features

- API REST para consultas.
- Integração com bancos de dados PostgreSQL, Redis e Elasticsearch.
- Cache de resultados.
- Handlers organizados por responsabilidade.
- Configuração via arquivo `.env`.

## Estrutura

```
/
├── database/           # Integração com bancos de dados
├── handlers/           # Handlers de cache e usuário
├── main.go             # Ponto de entrada da aplicação
├── Dockerfile          # Build e deploy com Docker
├── go.mod              # Dependências Go
└── .env                # Configurações de ambiente
```

---

> Projeto para fins de estudo e demonstração.  