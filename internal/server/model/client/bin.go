package client

type Bin struct {
	Data []byte
}

func (b Bin) String() string {
	return "binary"
}

func (b Bin) Kind() int {
	return 1
}