package models

type Group struct {
	tableName struct{} `sql:"focus.groups"`
	ID int32 `sql:"id"`
	Name string `sql:"name"`
}

type GroupEvent struct {
	tableName struct{} `sql:"focus.event_group"`
	ID int32 `sql:"id"`
	EventId int32 `sql:"event_id"`
	GroupId int32 `sql:"group_id"`
}