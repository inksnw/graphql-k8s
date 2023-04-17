package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/handler"
	"github.com/phuslu/log"
	"graphql-k8s/lib"
	"io"
	"os"
)

func main() {
	initLog()
	resources := lib.ResourceTypes()
	shc, err := lib.GenerateGraphQLSchema(resources, 2)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	r := gin.Default()
	r.Use(ModifyRequest())
	r.POST("/graphql", GraphqlHandler(shc))

	err = r.Run(":8080")
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}

func GraphqlHandler(shc *graphql.Schema) gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema: shc,
		Pretty: true,
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func ModifyRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body map[string]any
		data, err := c.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = json.Unmarshal(data, &body)
		if err != nil {
			return
		}
		ql := body["query"]

		astDoc, err := parser.Parse(parser.ParseParams{Source: ql})
		if err != nil {
			panic(err)
		}
		lib.AddPrefix(astDoc, nil)
		modifiedQuery := printer.Print(astDoc)
		body["query"] = modifiedQuery

		marshal, err := json.Marshal(body)
		if err != nil {
			return
		}
		log.Info().Msgf("原查询语句:\n  %s", ql)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(marshal))
		c.Next()

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
