package socket

import (
	"errors"
	"github.com/pebbe/zmq4"
)

var (
	frontend *zmq4.Socket
)

// InitSocket returns a socket listening on tcp://*:5555
func InitSocket() (*zmq4.Socket, error) {
	var err error
	frontend, err = zmq4.NewSocket(zmq4.ROUTER)

	if err != nil {
		return nil, errors.New("failed to create frontend" + err.Error())
	}

	err = frontend.SetCurveServer(1)
	if err == nil {
		err = frontend.SetCurveSecretkey("JTKVSB%%)wK0E.X)V>+}o?pNmC{O&4W4b!Ni{Lh6")
	} else {
		return nil, errors.New("failed to set frontend curve" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to set frontend curve" + err.Error())
	}

	err = frontend.Bind("tcp://*:5555")
	if err != nil {
		return nil, errors.New("failed to bind frontend" + err.Error())
	}

	return frontend, nil
}
