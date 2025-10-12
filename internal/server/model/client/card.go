package client

import "fmt"

type Card struct {
	Holder string 
	Cvc string
	Expiry string 
	Number string
}

func (c Card) String() string {
	return fmt.Sprintf("number %s holder %s expiry %s cvc %s", c.Number, c.Holder, c.Expiry, c.Cvc)
}

func (c Card) Kind() int {
	return 3
}