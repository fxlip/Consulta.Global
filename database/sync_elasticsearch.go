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

	rows, _ := db.Query(context.Background(), "SELECT id, nome FROM usuarios")
	defer rows.Close()

	for rows.Next() {
		var id int
		var nome string
		rows.Scan(&id, &nome)

		doc := fmt.Sprintf(`{ "nome": "%s" }`, nome)
		es.Index(
			"usuarios",
			strings.NewReader(doc),
			es.Index.WithDocumentID(strconv.Itoa(id)),
		)
	}
}
