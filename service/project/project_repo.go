package project

type Repository interface {
	GetAllProjectByUser(userId int) (projects []Project, err error)
	GetProjectById(projectId int) (project Project, err error)
	CreateProject(input Project) (project Project, err error)
	UpdateProject(input Project) (project Project, err error)
	DeleteProject(projectId int) (err error)
}
