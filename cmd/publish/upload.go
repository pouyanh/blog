package main

type Uploader interface {
	Upload(destination, source string) error
}

type FtpUploader struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (u FtpUploader) Upload(destination, source string) error {
	panic("not implemented")
}
