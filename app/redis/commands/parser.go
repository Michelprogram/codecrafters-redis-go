package commands

import (
	"bytes"
)

type CMD struct {
	ICommand
	Name       string
	Parameters [][]byte
}

type Parser struct {
	Commands map[string]ICommand
	Data     []byte
}

func NewParser() Parser {
	return Parser{
		Data: nil,
		Commands: map[string]ICommand{
			"ping":     Ping{},
			"echo":     Echo{},
			"set":      Set{},
			"get":      Get{},
			"info":     Info{},
			"replconf": ReplConf{},
			"psync":    Psync{},
			"redis001": RDB{},
			"type":     Type{},
			"xadd":     Xadd{},
			"xrange":   Xrange{},
			"xread":    XRead{},
			"config":   Config{},
			"keys":     Keys{},
		},
	}
}

func (p *Parser) arrays() *CMD {

	args := bytes.Split(p.Data, CRLF)

	arg := string(bytes.ToLower(args[2]))

	cmd, ok := p.Commands[arg]

	if ok {
		return &CMD{
			ICommand:   cmd,
			Name:       arg,
			Parameters: args[4:],
		}
	}

	return nil
}

func (p *Parser) bulkstring() *CMD {

	//firstCRLF := bytes.Index(p.Data, CRLF)

	//content := p.Data[firstCRLF+2:]

	//log.Println(strings.Split(string(content), string(CRLF)))

	return &CMD{
		ICommand:   p.Commands["redis001"],
		Name:       "redis001",
		Parameters: nil,
	}
}

func (p *Parser) ParseArgs(data []byte) *CMD {

	p.Data = data

	start := p.Data[0]

	switch start {
	case '$':
		return p.bulkstring()
	case '*':
		return p.arrays()
	default:
		return nil
	}

}
