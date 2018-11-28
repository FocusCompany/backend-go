package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type Event struct {
	tableName   struct{}  `sql:"focus.events"`
	ID          int32     `sql:"id" json:"id"`
	UserId      uuid.UUID `sql:"user_id,notnull" json:"userId"`
	DeviceId    int `sql:"device_id,notnull" json:"deviceId"`
	WindowsName string    `sql:"window_name" json:"windowName"`
	ProcessName string    `sql:"process_name" json:"processName"`
	Afk         bool      `sql:"afk" json:"afk"`
	Time        time.Time `sql:"time" json:"time"`
}
