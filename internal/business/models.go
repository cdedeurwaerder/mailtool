package business

type User struct {
	Firstname string
	Lastname  string
	Email     string
}

type Email struct {
	ID           string
	From         string
	To           string
	EnvelopeFrom string
	RcptTo       string
	Title        string
	Body         string
	Attachments  []byte
}

type SuspiciousEmail struct {
	ID           string
	From         string
	To           string
	EnvelopeFrom string
	RcptTo       string
	Title        string
	Body         string
	Flag         int
	Reason       string
}
