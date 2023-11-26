package compilers

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/janstoon/toolbox/tricks"
	"github.com/russross/blackfriday/v2"
	"github.com/yuin/goldmark"
	"golang.org/x/sync/errgroup"

	"github.com/pouyanh/blog"
)

type localCompiler struct {
	RootTemplate string
	RootFilename string

	ArticleTemplate string
	ArticlesDir     string
}

func LocalCompiler(opts ...tricks.Option[localCompiler]) blog.Compiler {
	c := localCompiler{
		RootTemplate: "root.gohtml",
		RootFilename: "index.html",

		ArticleTemplate: "article.gohtml",
		ArticlesDir:     "articles",
	}
	tricks.ApplyOptions(&c, opts...)

	return c
}

func (c localCompiler) Compile(tpl *template.Template, destination string, sources ...string) error {
	err := os.MkdirAll(destination, 0755)
	if nil != err {
		return err
	}

	articles, err := c.compileArticles(tpl, destination, sources...)
	if nil != err {
		return err
	}

	err = c.compileIndex(tpl, destination, articles...)
	if nil != err {
		return err
	}

	return nil
}

func (c localCompiler) compileArticles(
	tpl *template.Template, destination string, sources ...string,
) ([]article, error) {
	dstArticles := filepath.Join(destination, c.ArticlesDir)
	err := os.MkdirAll(dstArticles, 0755)
	if nil != err {
		return nil, err
	}

	aa := make([]article, len(sources))
	wg := errgroup.Group{}
	for k, filename := range sources {
		func(k int, filename string) {
			wg.Go(func() error {
				a, err := c.compileArticle(tpl, dstArticles, filename)
				if err != nil {
					return err
				}

				aa[k] = tricks.PtrVal(a)

				return nil
			})
		}(k, filename)
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return aa, nil
}

func (c localCompiler) compileArticle(tpl *template.Template, destination, source string) (*article, error) {
	bb, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}

	lang := extensionless(filepath.Base(source))
	an := filepath.Base(filepath.Dir(source))
	bn := fmt.Sprintf("%s.%s.html", an, lang)

	doc := blackfriday.New().Parse(bb)
	title := getTitle(doc, an)

	fn := filepath.Join(destination, bn)
	fh, err := os.Create(fn)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fh.Close() }()

	md := goldmark.New()
	err = md.Convert(bb, fh)
	if err != nil {
		return nil, err
	}

	return &article{
		Title: title,
		Lang:  lang,
		Link:  filepath.Join(c.ArticlesDir, bn),
	}, nil
}

func (c localCompiler) compileIndex(tpl *template.Template, destination string, articles ...article) error {
	fh, err := os.Create(filepath.Join(destination, c.RootFilename))
	if nil != err {
		return err
	}
	defer func() { _ = fh.Close() }()

	err = tpl.ExecuteTemplate(fh, c.RootTemplate, newArgs(articles...))
	if nil != err {
		return err
	}

	return nil
}

func extensionless(basename string) string {
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func getTitle(doc *blackfriday.Node, fallback string) string {
	heading := doc.FirstChild
	for heading != nil && heading.Level != 1 {
		heading = heading.Next
	}

	if heading == nil || heading.FirstChild == nil {
		return fallback
	}

	return string(heading.FirstChild.Literal)
}

type args struct {
	Header    header
	Defaults  defaults
	Languages map[string]lang
	Articles  []article
}

func newArgs(articles ...article) args {
	scope := args{
		Header: newHeader(),
		Defaults: defaults{
			Lang: "en",
		},
		Languages: make(map[string]lang),
		Articles:  articles,
	}

	for _, a := range articles {
		scope.Languages[a.Lang] = langs[a.Lang]
	}

	return scope
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

type defaults struct {
	Lang string
}

type lang struct {
	UpperCode string // Contains ISO 639-1 two-letter code in uppercase
	Title     string
	Rtl       bool
}

var langs = map[string]lang{
	"en": {
		UpperCode: "EN",
		Title:     "English",
		Rtl:       false,
	},
	"fa": {
		UpperCode: "FA",
		Title:     "فارسی",
		Rtl:       true,
	},
}

type article struct {
	Title string
	Lang  string
	Link  string
}
