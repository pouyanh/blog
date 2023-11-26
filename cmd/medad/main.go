package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/pouyanh/blog"
	"github.com/pouyanh/blog/compilers"
	"github.com/pouyanh/blog/uploaders"
)

func main() {
	var cmd medad
	if err := cmd.run(context.Background()); err != nil {
		log.Fatalf("error: %s", err)
	}
}

type medad struct {
	stn blog.Settings
}

func (c *medad) run(ctx context.Context) error {
	root := cobra.Command{
		Use:              "medad",
		TraverseChildren: true,
	}
	root.Flags().StringVarP(&c.stn.DistDirectory, "dist", "d", "public", "Dist directory")

	compile := &cobra.Command{
		Use:   "compile",
		Short: "Compiles articles to html",
		RunE:  c.compile,
	}
	compile.Flags().StringVar(&c.stn.ArticlesGlob, "articles",
		filepath.Join("content", "articles", "*", "*.md"), "Article files glob (pattern)")
	compile.Flags().StringVarP(&c.stn.TemplatesGlob, "templates", "t", filepath.Join("templates", "*.gohtml"),
		"Template files glob (pattern)")

	upload := &cobra.Command{
		Use:   "upload",
		Short: "Uploads compiled html files to server",
		RunE:  c.upload,
	}
	upload.Flags().StringVarP(&c.stn.FtpHost, "host", "H", "example.com",
		"FTP Host to upload")
	upload.Flags().StringVarP(&c.stn.FtpPort, "port", "P", "21",
		"FTP Port to upload")
	upload.Flags().StringVarP(&c.stn.FtpUsername, "username", "u", "anonymous",
		"FTP Username to upload")
	upload.Flags().StringVarP(&c.stn.FtpPassword, "password", "p", "",
		"FTP Password to upload")
	upload.Flags().StringVar(&c.stn.RemoteDirectory, "rdir", "blog",
		"Remote directory to upload dist to")

	update := &cobra.Command{
		Use:   "update",
		Short: "Combination of compile & upload",
		RunE:  c.update,
	}
	update.Flags().AddFlagSet(compile.Flags())
	update.Flags().AddFlagSet(upload.Flags())

	root.AddCommand(compile)
	root.AddCommand(upload)
	root.AddCommand(update)

	return root.ExecuteContext(ctx)
}

func (c *medad) compile(cmd *cobra.Command, args []string) error {
	templates, err := filepath.Glob(c.stn.TemplatesGlob)
	if err != nil {
		return fmt.Errorf("reading templates error: %w", err)
	}
	tpl := template.Must(template.ParseFiles(templates...))
	log.Printf("template list loaded: %+v\n", templates)

	sources, err := filepath.Glob(c.stn.ArticlesGlob)
	if err != nil {
		return fmt.Errorf("reading contents error: %w", err)
	}
	log.Printf("source list to process: %+v\n", sources)
	lc := compilers.LocalCompiler()
	err = lc.Compile(tpl, c.stn.DistDirectory, sources...)
	if err != nil {
		return fmt.Errorf("compile error: %w", err)
	}

	return nil
}

func (c *medad) upload(cmd *cobra.Command, args []string) error {
	fu := uploaders.FtpUploader{
		Username: c.stn.FtpUsername,
		Password: c.stn.FtpPassword,
		Host:     c.stn.FtpHost,
		Port:     c.stn.FtpPort,

		Timeout: 5 * time.Second,
	}

	err := fu.Upload(c.stn.RemoteDirectory, c.stn.DistDirectory)
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}

	return nil
}

func (c *medad) update(cmd *cobra.Command, args []string) error {
	if err := c.compile(cmd, args); err != nil {
		return err
	}

	if err := c.upload(cmd, args); err != nil {
		return err
	}

	return nil
}
