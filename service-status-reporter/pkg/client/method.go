package client

type Method interface {
	GetValue() string
}

type GET struct{}

func (g GET) GetValue() string {
	return "GET"
}

type POST struct{}

func (p POST) GetValue() string {
	return "POST"
}

type PUT struct{}

func (p PUT) GetValue() string {
	return "PUT"
}

type DELETE struct{}

func (d DELETE) GetValue() string {
	return "DELETE"
}

type PATCH struct{}

func (p PATCH) GetValue() string {
	return "PATCH"
}

type HEAD struct{}

func (h HEAD) GetValue() string {
	return "HEAD"
}

type OPTIONS struct{}

func (o OPTIONS) GetValue() string {
	return "OPTIONS"
}

type TRACE struct{}

func (t TRACE) GetValue() string {
	return "TRACE"
}

type CONNECT struct{}

func (c CONNECT) GetValue() string {
	return "CONNECT"
}
