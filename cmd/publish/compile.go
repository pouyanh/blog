package main

import (
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type Compiler interface {
	Compile(tpl *template.Template, destination string, sources ...string) error
}

type LocalCompiler struct {
}

func (c LocalCompiler) Compile(tpl *template.Template, destination string, sources ...string) error {
	err := os.MkdirAll(destination, 0755)
	if nil != err {
		return err
	}

	fh, err := os.Create(filepath.Join(destination, "index.html"))
	if nil != err {
		return err
	}
	defer func() { _ = fh.Close() }()

	err = tpl.ExecuteTemplate(fh, "root.html.tpl", newArgs())
	if nil != err {
		return err
	}

	return nil
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
