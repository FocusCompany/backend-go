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
	socket, _ := zmq4.NewSocket(zmq4.PUSH)

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
			WindowName:          []byte(windows[rand.Int() % len(windows)]),
			ProcessName:         []byte(processes[rand.Int() % len(processes)]),
		}
		any, _:= ptypes.MarshalAny(&payload)

		event := Focus.Event{
			Timestamp:            &timestamp.Timestamp{Seconds:previousTime.Unix()},
			PayloadType:          "ContextChanged",
			Payload:              any,
		}
		previousTime = previousTime.Add(time.Duration(rand.Int() % 40) * time.Minute)

		envelope := Focus.Envelope{
			DeviceID:             "someID",
			Jwt:                  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMjY4YTQ4YzctZjliNi00MzI0LWI5MzEtYmFkYmM3Y2M1YzM5IiwiZXhwIjoxNTQzMTYzNzI5LCJpYXQiOjE1NDMxNjE5Mjl9.vLj7VwwDIAM5_8amfyqyGMRvYTvu3iejLseTYroL8A8xAvrufs6AuJ_00DfKJZ91kL4W0q4TqMyAPbX5zRFDgS3o-fwdkejvzTyOCci3fRmow3449ROQu9acBagW76SOolwEa731WdJU-iC1HZJCA6Z74o1-NdLlCrAkeFEHSbw",
			Events: []*Focus.Event{&event},
		}

		bytes, _ := proto.Marshal(&envelope)

		socket.Send(string(bytes), 0)
	}
	time.Sleep(time.Second)
}