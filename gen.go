package main

import (
	"log"
	"os"
	"path"
	"text/template"
	"time"
)

//go:generate go run gen.go

const (
	dirnameOut    = "public"
	basenameIndex = "index"
	filenameDst   = basenameIndex + ".html"
	filenameTpl   = filenameDst + ".tpl"
)

var (
	tplIndex = template.Must(template.ParseFiles(filenameTpl))
)

func main() {
	err := os.Mkdir(dirnameOut, 0755)
	if nil != err {
		log.Fatal(err)
	}

	fh, err := os.Create(path.Join(dirnameOut, filenameDst))
	if nil != err {
		log.Fatal(err)
	}
	defer func() { _ = fh.Close() }()

	err = tplIndex.Execute(fh, newArgs())
	if nil != err {
		log.Fatal(err)
	}
}

type args struct {
	Header header
}

func newArgs() args {
	return args{
		Header: newHeader(),
	}
}

type header struct {
	Date time.Time

	OS   string
	Arch string

	PackageName string
	Filename    string
	LineNumber  string
}

func newHeader() header {
	return header{
		Date: time.Now(),

		OS:   os.Getenv("GOOS"),
		Arch: os.Getenv("GOARCH"),

		PackageName: os.Getenv("GOPACKAGE"),
		Filename:    os.Getenv("GOFILE"),
		LineNumber:  os.Getenv("GOLINE"),
	}
}
