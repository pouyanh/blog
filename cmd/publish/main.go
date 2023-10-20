package main

import (
	"html/template"
	"log"
	"path/filepath"

	"github.com/pouyanh/blog"
	"github.com/pouyanh/blog/compilers"
)

func main() {
	cfg := blog.Settings{
		OutputDirectory: "public",
	}

	templates, err := filepath.Glob(filepath.Join("templates", "*.html.tpl"))
	if err != nil {
		log.Fatalf("reading templates error: %s", err)
	}
	tpl := template.Must(template.ParseFiles(templates...))
	log.Printf("template list loaded: %+v\n", templates)

	sources, err := filepath.Glob(filepath.Join("content", "articles", "*", "*.md"))
	if err != nil {
		log.Fatalf("reading contents error: %s", err)
	}
	log.Printf("source list to process: %+v\n", sources)
	lc := compilers.Local{}
	err = lc.Compile(tpl, cfg.OutputDirectory, sources...)
	if err != nil {
		log.Fatalf("compile error: %s", err)
	}
}
