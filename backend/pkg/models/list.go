package models

type TaskList struct {
	Id       int    `json:"id,omitempty"`
	BoardId  int    `json:"boardId" db:"board_id"`
	Title    string `json:"title" valid:"length(1|32)"`
	Position int    `json:"position" valid:"type(int)"`
}

type UpdateTaskList struct {
	Title    *string `json:"title"`
	Position *int    `json:"position" valid:"type(*int)"`
}
