package commands

type CommandError struct {
	Err error
	BuilderRESP
}

func (c CommandError) Error() string {
	return c.EncodeAsSimpleString(c.Err.Error(), ERROR).String()
}
