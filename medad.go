package blog

import "html/template"

type Compiler interface {
	Compile(tpl *template.Template, destination string, sources ...string) error
}

type Uploader interface {
	Upload(destination, source string) error
}

type Settings struct {
	ArticlesGlob  string
	TemplatesGlob string
	DistDirectory string

	FtpHost         string
	FtpPort         string
	FtpUsername     string
	FtpPassword     string
	RemoteDirectory string
}
