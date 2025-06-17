package database // Adicione esta linha

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SyncToElasticsearch(db *pgxpool.Pool) {
	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_HOST")},
	})

	rows, _ := db.Query(context.Background(), "SELECT id, primeiro_nome, sobrenome, nome_completo FROM pessoas_fisicas")
	defer rows.Close()

	for rows.Next() {
		var id int
		var primeiroNome, sobrenome, nomeCompleto string
		rows.Scan(&id, &primeiroNome, &sobrenome, &nomeCompleto)

		doc := fmt.Sprintf(`{ "primeiro_nome": "%s", "sobrenome": "%s", "nome_completo": "%s" }`,
			primeiroNome, sobrenome, nomeCompleto)
		es.Index(
			"pessoas_fisicas",
			strings.NewReader(doc),
			es.Index.WithDocumentID(strconv.Itoa(id)),
		)
	}
}
