package job

type Repository interface {
	GetAllJobs(input string) (jobs []Job, err error)
	GetAllJobByUser(userId int) (jobs []Job, err error)
	GetJobById(jobId int) (job Job, err error)
	CreateJob(input Job) (job Job, err error)
	UpdateJob(input Job) (job Job, err error)
	DeleteJob(jobId int) (err error)
}
