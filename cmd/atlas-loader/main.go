package main

import (
	"fmt"
	"io"
	"os"

	"go-cms/internal/model"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&model.User{},
		&model.Blog{},
		&model.Portfolio{},
		&model.Session{},
		&model.SysConfig{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
