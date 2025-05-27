package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// GetSurnameSuggestions retorna sugestões de sobrenomes com base em um primeiro nome
func GetSurnameSuggestions(db *pgxpool.Pool, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		firstName := c.Query("firstName")
		if firstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nome não fornecido"})
			return
		}

		ctx := context.Background()

		// Tenta obter do cache primeiro
		cacheKey := fmt.Sprintf("surname-suggestions:%s", firstName)
		cachedSuggestions, err := rdb.Get(ctx, cacheKey).Bytes()
		if err == nil {
			c.Data(http.StatusOK, "application/json", cachedSuggestions)
			return
		}

		// Se não estiver em cache, busca no banco de dados
		query := `
            SELECT DISTINCT
                SUBSTRING(nome FROM POSITION(' ' IN nome) + 1) AS surname
            FROM 
                usuarios
            WHERE 
                nome ILIKE $1 || ' %' 
                AND POSITION(' ' IN nome) > 0
            ORDER BY 
                surname
            LIMIT 10
        `

		rows, err := db.Query(ctx, query, firstName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar sugestões"})
			return
		}
		defer rows.Close()

		var surnames []string
		for rows.Next() {
			var surname string
			if err := rows.Scan(&surname); err != nil {
				continue // Ignora erros individuais
			}
			surnames = append(surnames, surname)
		}

		// Armazena em cache por 1 hora
		jsonData, _ := json.Marshal(surnames)
		rdb.Set(ctx, cacheKey, jsonData, time.Hour)

		c.JSON(http.StatusOK, surnames)
	}
}

func GetUserDetails(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Busca dados do usuário
		var usuario struct {
			ID        int
			Nome      string
			CPF_CNPJ  string
			Operadora string
		}
		err := db.QueryRow(context.Background(),
			`SELECT id, nome, cpf_cnpj, operadora 
             FROM usuarios 
             WHERE id = $1`, id,
		).Scan(&usuario.ID, &usuario.Nome, &usuario.CPF_CNPJ, &usuario.Operadora)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}

		// Busca telefones
		type Telefone struct {
			DDD    string `json:"ddd"`
			Numero string `json:"numero"`
			Tipo   string `json:"tipo"`
		}

		var telefones []Telefone

		rowsTel, err := db.Query(context.Background(),
			"SELECT ddd, numero, tipo FROM telefones WHERE usuario_id = $1", id)
		if err == nil {
			defer rowsTel.Close()
			for rowsTel.Next() {
				var tel Telefone // Usando o tipo definido
				if err := rowsTel.Scan(&tel.DDD, &tel.Numero, &tel.Tipo); err == nil {
					telefones = append(telefones, tel)
				}
			}
		}

		// Busca endereços
		type Endereco struct {
			Logradouro     string `json:"logradouro"`
			NumeroEndereco string `json:"numero"`
			Cidade         string `json:"cidade"`
			UF             string `json:"uf"`
		}

		var enderecos []Endereco

		rowsEnd, err := db.Query(context.Background(),
			"SELECT logradouro, numero_endereco, cidade, uf FROM enderecos WHERE usuario_id = $1", id)
		if err == nil {
			defer rowsEnd.Close()
			for rowsEnd.Next() {
				var end Endereco // Usando o tipo definido
				if err := rowsEnd.Scan(&end.Logradouro, &end.NumeroEndereco, &end.Cidade, &end.UF); err == nil {
					enderecos = append(enderecos, end)
				}
			}
		}

		// Busca contatos adicionais
		type Contato struct {
			Tipo       string `json:"tipo"`
			Valor      string `json:"valor"`
			Observacao string `json:"observacao,omitempty"`
		}

		var contatos []Contato

		rowsCont, err := db.Query(context.Background(),
			"SELECT tipo, valor, observacao FROM contatos_adicionais WHERE usuario_id = $1", id)
		if err == nil {
			defer rowsCont.Close()
			for rowsCont.Next() {
				var cont Contato // Usando o tipo definido
				if err := rowsCont.Scan(&cont.Tipo, &cont.Valor, &cont.Observacao); err == nil {
					contatos = append(contatos, cont)
				}
			}
		}

		// Adiciona ao objeto de resposta
		response := struct {
			Usuario   interface{} `json:"usuario"`
			Telefones interface{} `json:"telefones"`
			Enderecos interface{} `json:"enderecos"`
			Contatos  interface{} `json:"contatos_adicionais"`
		}{
			Usuario:   usuario,
			Telefones: telefones,
			Enderecos: enderecos,
			Contatos:  contatos,
		}

		c.JSON(http.StatusOK, response)

	}
}

func SearchUsers(db *pgxpool.Pool, rdb *redis.Client, es *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		term := c.Query("term")
		ctx := context.Background()

		// 1. Verifica cache no Redis
		cachedResult, err := rdb.Get(ctx, "search:"+term).Bytes()
		if err == nil {
			c.Data(http.StatusOK, "application/json", cachedResult)
			return
		}

		// 2. Busca no Elasticsearch
		query := fmt.Sprintf(`
		{
			"query": {
				"match": {
					"nome": {
						"query": "%s",
						"fuzziness": "AUTO"
					}
				}
			}
		}`, term)

		res, err := es.Search(
			es.Search.WithIndex("usuarios"),
			es.Search.WithBody(strings.NewReader(query)),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha na busca"})
			return
		}

		// 3. Processa resultados
		var result map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decode falhou"})
			return
		}

		hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
		results := make([]map[string]interface{}, 0)
		for _, hit := range hits {
			source := hit.(map[string]interface{})["_source"].(map[string]interface{})
			results = append(results, map[string]interface{}{
				"id":   hit.(map[string]interface{})["_id"],
				"nome": source["nome"],
			})
		}

		// 4. Armazena no Redis com TTL dinâmico
		jsonData, _ := json.Marshal(results)
		ttl := 5 * time.Minute
		if len(results) > 0 {
			ttl = 30 * time.Minute
		}
		rdb.Set(ctx, "search:"+term, jsonData, ttl)

		c.JSON(http.StatusOK, results)
	}
}
