package main

import (
	"fmt"
	"github.com/graphql-go/handler"
	"github.com/phuslu/log"
	"graphql-k8s/lib"
	"net/http"
	"os"
)

func main() {
	initLog()
	resources := lib.ResourceTypes()

	var test []lib.ResourceType
	for _, i := range resources {
		if i.ResourceName == "io.k8s.api.core.v1.Pod" {
			test = append(test, i)
		}
	}

	schema, err := lib.GenerateGraphQLSchema(test, 3)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("Listening on :8080/graphql...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}

func initLog() {
	if !log.IsTerminal(os.Stderr.Fd()) {
		return
	}
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}
