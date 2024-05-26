package api

import (
	"ethereum-parser/internal/domain/parser"
	"log"
)

func Setup() {
	p, err := parser.NewParser()
	if err != nil {
		log.Fatalf("set up dep error:[%s]", err.Error())
	}
	parser.SetParser(p)
}
