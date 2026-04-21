package api

type Mailer interface {
	SendMagicLink(recipient, link string) error
}
