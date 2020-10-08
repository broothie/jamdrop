package job

import "jamdrop/app"

type Job struct {
	*app.App
}

func New(app *app.App) *Job {
	return &Job{App: app}
}
