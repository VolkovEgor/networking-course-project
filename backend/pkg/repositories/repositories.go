package repositories

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetAll() ([]*models.User, error)
	GetById(id int) (*models.User, error)
	Create(user *models.User) (int, error)
	Get(nickname, password string) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	SignOut(token string) (int, error)
	FindToken(token string) error
	Update(id int, profile *models.UpdateUser) error
}

type Project interface {
	Create(project *models.Project) (int, error)
	GetAll(userId int) ([]*models.Project, error)
	GetById(projectId int) (*models.Project, error)
	Delete(projectId int) error
	Update(projectId int, project *models.UpdateProject) error
	GetPermissions(userId, projectId int) (*models.Permission, error)
	GetMembers(projectId int) ([]*models.Member, error)
}

type Board interface {
	Create(userId int, board *models.Board) (int, error)
	GetAll(userId, projectId int) ([]*models.Board, error)
	GetById(boardId int) (*models.Board, error)
	Delete(boardId int) error
	Update(boardId int, board *models.UpdateBoard) error
	GetPermissions(userId, boardId int) (*models.Permission, error)
	GetBoardsCountByOwnerId(projectId, ownerId int) (int, error)
	GetMembers(projectId int) ([]*models.Member, error)
}

type TaskList interface {
	Create(list *models.TaskList) (int, error)
	GetAll(listId int) ([]*models.TaskList, error)
	GetById(listId int) (*models.TaskList, error)
	Delete(listId int) error
	Update(listId int, list *models.UpdateTaskList) error
	// GetPermissions(userId, boardId int) (*models.Permission, error)
}

type Task interface {
	Create(task *models.Task) (int, error)
	GetAll(taskId int) ([]*models.Task, error)
	GetById(taskId int) (*models.Task, error)
	Delete(taskId int) error
	Update(taskId int, task *models.UpdateTask) error
}

type Label interface {
	Create(label *models.Label) (int, error)
	CreateInTask(taskId, labelId int) (int, error)
	GetAllInTask(taskId int) ([]*models.Label, error)
	GetAll(boardId int) ([]*models.Label, error)
	GetById(labelId int) (*models.Label, error)
	DeleteInTask(taskId, labelId int) error
	Delete(labelId int) error
	Update(labelId int, label *models.UpdateLabel) error
}

type ObjectPerms interface {
	Create(objectId, objectType int, memberNickname string, permissions *models.Permission) (int, error)
	GetById(objectId, memberId, objectType int) (*models.Permission, error)
	GetByNickname(objectId, objectType int, memberId string) (*models.Permission, error)
	Delete(objectId, oldOwnerId, newOwnerId, objectType int) error
	Update(objectId, oldOwnerId, newOwnerId, objectType int, permissions *models.UpdatePermission) error
}

type Repository struct {
	User
	Project
	Board
	TaskList
	Task
	Label
	ObjectPerms
}

// func NewRepository(db *mongo.Database) *Repository {
// 	return &Repository{
// 		User:     repoMongo.NewUserMongo(db),
// 		Project:  repoMongo.NewProjectMongo(db),
// 		Board:    repoMongo.NewBoardMongo(db),
// 		TaskList: repoMongo.NewTaskListMongo(db),
// 		Task:     repoMongo.NewTaskMongo(db),
// 	}
// }

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        postgres.NewUserPg(db),
		Project:     postgres.NewProjectPg(db),
		Board:       postgres.NewBoardPg(db),
		TaskList:    postgres.NewTaskListPg(db),
		Task:        postgres.NewTaskPg(db),
		Label:       postgres.NewLabelPg(db),
		ObjectPerms: postgres.NewObjectPermsPg(db),
	}
}
