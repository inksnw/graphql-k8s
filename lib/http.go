package lib

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/phuslu/log"
	"io"
	"strings"
)

func GraphqlHandler(actionMap map[string]*graphql.Schema) gin.HandlerFunc {

	return func(c *gin.Context) {
		var body map[string]any
		data, err := c.GetRawData()
		if err != nil {
			log.Fatal().Msgf("%s", err)
		}
		err = json.Unmarshal(data, &body)
		if err != nil {
			log.Fatal().Msgf("%s", err)
		}
		ql := body["query"]
		list := strings.Split(ql.(string), "{")
		kind := strings.TrimSpace(list[1])
		shc, ok := actionMap[strings.ToLower(kind)]
		if !ok {
			log.Error().Msgf("不支持的类型 %s", kind)
		}

		marshal, err := json.Marshal(body)
		if err != nil {
			log.Fatal().Msgf("%s", err)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(marshal))
		h := handler.New(&handler.Config{
			Schema: shc,
			Pretty: true,
		})

		h.ServeHTTP(c.Writer, c.Request)
	}
}
