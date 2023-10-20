package uploaders

type FtpUploader struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (u FtpUploader) Upload(destination, source string) error {
	panic("not implemented")
}
