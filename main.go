package main

import (
	"log"
	"net/http"
	"encoding/json"
	"context"
	"strings"
	"mime"
	"path/filepath"

	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/ws"
	
	"github.com/mikeifomin/graphql-subscriptions-go-example/schema"
	"github.com/mikeifomin/graphql-subscriptions-go-example/assets"
)

func main() {
	schemaCompiled := graphql.MustParseSchema(schema.Schema, &schema.Resolver{})

	http.Handle("/ws", &ws.Handler{
		Schema: schemaCompiled,
		OnConnect: func(connectionParam json.RawMessage, request *http.Request) (context.Context, error){
      return context.Background(),nil
		},
	})

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path,"/")
		if filename == "" {
			filename = "index.html"
		}
		b, err := assets.Asset(filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m := mime.TypeByExtension(filepath.Ext(filename))
		w.Header().Set("Content-Type",m)
		w.Write(b)
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
