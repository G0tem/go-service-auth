package dto

type SmtpMessageType uint8

const (
	MailUnknown SmtpMessageType = iota
	MailPlain
	MailHtml
)

func ParseSmtpMessageType(smtpMessageType string) SmtpMessageType {
	switch smtpMessageType {
	case "@mail.Plain":
		return MailPlain
	case "@mail.Html":
		return MailHtml
	}
	return MailUnknown
}

func (smtpMessageType SmtpMessageType) String() string {
	switch smtpMessageType {
	case MailPlain:
		return "@mail.Plain"
	case MailHtml:
		return "@mail.Html"
	}
	return ""
}
