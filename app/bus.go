package app

type CommandIssuer interface {
	Send(cmd interface{}) error
}