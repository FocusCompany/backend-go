package models

import "github.com/satori/go.uuid"

type Dnd struct {
	tableName struct{} `sql:"focus.dnd"`
	UserId uuid.UUID `sql:"user_id"`
	Activations int `sql:"activations"`
}
