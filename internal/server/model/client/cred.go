package client

import "fmt"

type Credential struct {
	Login string
	Password string
}

func (c Credential) String() string {
	return fmt.Sprintf("login %s password %s", c.Login, c.Password)
}

func (c Credential) Kind() int {
	return 2
}