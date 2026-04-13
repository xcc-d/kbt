package service

type HelloService struct {
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) SayHello(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}
