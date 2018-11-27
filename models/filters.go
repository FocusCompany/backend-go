package models

import (
	"github.com/satori/go.uuid"
)

type Filters struct {
	tableName struct{}  `sql:"focus.filters"`
	UserId    uuid.UUID `sql:"user_id"`
	Name      string    `sql:"name"`
}

