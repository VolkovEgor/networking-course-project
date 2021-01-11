package models

type Task struct {
	Id          int        `json:"_id,omitempty"`
	ListId      int        `json:"listId"`
	Title       string     `json:"title" valid:"length(1|32)"`
	Description string     `json:"description,omitempty"`
	Datetimes   *Datetimes `json:"datetimes,omitempty"`
	Position    int        `json:"position" valid:"type(int)"`
}

type UpdateTask struct {
	ListId      *int             `json:"listId" valid:"type(*int)"`
	Title       *string          `json:"title" valid:"length(1|32)"`
	Description *string          `json:"description,omitempty"`
	Datetimes   *UpdateDatetimes `json:"datetimes,omitempty"`
	Position    *int             `json:"position" valid:"type(*int)"`
}
