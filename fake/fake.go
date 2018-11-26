package main

import (
	"fmt"
	"github.com/FocusCompany/backend-go/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pebbe/zmq4"
	"math/rand"
	"time"
)

func main() {
	socket, _ := zmq4.NewSocket(zmq4.DEALER)

	err := socket.SetCurveServerkey("rq:rM>}U?@Lns47E1%kR.o@n%FcmmsL/@{H8]yf7")
	err = socket.SetCurvePublickey("Yne@$w-vo<fVvi]a<NY6T1ed:M$fCG*[IaLV{hID")
	err = socket.SetCurveSecretkey("D:)Q[IlAW!ahhC2ac:9*A}h:p?([4%wOTJ%JR%cs")
	if err != nil {
		fmt.Println(err)
	}

	socket.Connect("tcp://127.0.0.1:5555")
	defer socket.Close()

	rand.Seed(time.Now().Unix())
	windows := []string{
		"Facebook",
		"Youtube",
		"Gmail",
		"( ͡° ͜ʖ ͡°)",
	}

	processes := []string{
		"chrome.exe",
		"outlook",
	}

	previousTime := time.Now()
	for i := 0; i < 10; i++ {

		payload := Focus.ContextEventPayload{
			WindowName:  []byte(windows[rand.Int()%len(windows)]),
			ProcessName: []byte(processes[rand.Int()%len(processes)]),
		}
		any, _ := ptypes.MarshalAny(&payload)

		event := Focus.Event{
			Timestamp:   &timestamp.Timestamp{Seconds: previousTime.Unix()},
			PayloadType: "ContextChanged",
			Payload:     any,
		}
		previousTime = previousTime.Add(time.Duration(rand.Int() % 40) * time.Minute)

		envelope := Focus.Envelope{
			DeviceID: "23",
			Jwt:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMjY4YTQ4YzctZjliNi00MzI0LWI5MzEtYmFkYmM3Y2M1YzM5IiwiZXhwIjoxNTQzMjc3MTEwLCJpYXQiOjE1NDMyNzUzMTB9.GjMgbcdrIg-_aT85Rwl_k6m2OH2M4pf8rls6utGe-9SlfSBcsPO1J85gNsISOIAzIylZ9T3dLx1HXXdMBT2k5wsVErPCiQaYOGDccvzBNE1pz6ABqjT9HECExyBFKaadhYWc6xsH5o0c0OlrVogud45tHcMNbEIXpspYITPe3Bk",
			Events:   []*Focus.Event{&event},
		}

		bytes, _ := proto.Marshal(&envelope)

		socket.Send(string(bytes), 0)
	}
	time.Sleep(time.Second)
}
