package conn

import (
	"log"
)

var (
	connErrF = "连接外部程序出错%s"
)

type Conner interface {
	Conn() error
	Close() error
}

type ExternalProcedure struct {
	Conners []Conner
}

func NewExternalProcedure(conners ...Conner) *ExternalProcedure {
	e := &ExternalProcedure{}
	for _, conner := range conners {
		err := conner.Conn()
		if err != nil {
			log.Printf(connErrF, err)
		}
		e.Conners = append(e.Conners, conner)
	}
	return e
}

func (e *ExternalProcedure) Close() {
	for _, conner := range e.Conners {
		_ = conner.Close()
	}
}
