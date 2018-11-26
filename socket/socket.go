package socket

import (
	"errors"
	"github.com/pebbe/zmq4"
)

// InitSocket returns a socket listening on tcp://*:5555
func InitSocket() (*zmq4.Socket, error) {
	socket, err := zmq4.NewSocket(zmq4.DEALER)

	if err != nil {
		return nil, errors.New("failed to create socket" + err.Error())
	}

	err = socket.SetCurveServer(1)
	if err == nil {
		err = socket.SetCurveSecretkey("JTKVSB%%)wK0E.X)V>+}o?pNmC{O&4W4b!Ni{Lh6")
	} else {
		return nil, errors.New("failed to set socket curve" + err.Error())
	}

	if err != nil {
		return nil, errors.New("failed to set socket curve" + err.Error())
	}

	err = socket.Bind("tcp://*:5555")
	if err != nil {
		return nil, errors.New("failed to bind socket" + err.Error())
	}

	return socket, nil
}
