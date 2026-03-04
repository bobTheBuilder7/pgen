package main

import (
	"context"
	"strings"

	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/valkdb/postgresparser"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "schema"

var tablesCol = syncmap.Map[string, []postgresparser.DDLColumn]{}

func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}

func main() {
	// text := "select * from users where id = @id_asd;"

	// re := regexp.MustCompile(`@([a-zA-Z0-9_]+)[,); ]`)
	// reMatches := re.FindAllString(text, -1)

	// var matches []string

	// for _, m := range reMatches {
	// 	matches = append(matches, m[:len(m)-1])
	// }

	// fmt.Println(matches)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx)
	if err != nil {
		panic(err.Error())
	}
}
