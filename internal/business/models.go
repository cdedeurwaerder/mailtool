package business

import "github.com/google/uuid"

type EmailFlag int

const (
	EmailFraudSuspicious EmailFlag = iota + 1
	EmailSpearPhishingSuspicious
)

type User struct {
	ID        uuid.UUID
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
	UserID       string
	EmailID      string
	From         string
	To           string
	EnvelopeFrom string
	RcptTo       string
	Title        string
	Body         string
	Flag         EmailFlag
	Reason       string
}
