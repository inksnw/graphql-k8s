package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"graphql-k8s/lib"
	"os"
)

const Depth = 2

func main() {
	initLog()
	resources := lib.ResourceTypes()
	actionMap, err := lib.GenerateGraphQLSchema(resources, Depth)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	r := gin.Default()

	r.POST("/graphql", lib.GraphqlHandler(actionMap))

	err = r.Run(":8080")
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
