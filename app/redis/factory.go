package redis

func NewServerFactory(port uint, role string) IServer {

	if role == "" {
		return newMain(port, role)
	}

	return newNode(port, role)

}
