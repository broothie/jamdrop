package job

import (
	"context"
	"time"

	"jamdrop/app"
	"jamdrop/model"
)

type Job struct {
	*app.App
}

func New(app *app.App) *Job {
	return &Job{App: app}
}

func (j *Job) EjectSessionTokens(ctx context.Context) error {
	j.Logger.Println("job.EjectSessionTokens")

	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	j.Logger.Printf("deleting session tokens last updated_at before %v\n", thirtyDaysAgo)

	docs, err := j.DB.
		Collection(model.CollectionSessionTokens).
		Where("updated_at", "<", thirtyDaysAgo).
		Documents(ctx).
		GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		j.Logger.Println("no expired session tokens found")
		return nil
	}

	batch := j.DB.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	j.Logger.Printf("deleting %d session tokens\n", len(docs))
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}

	return nil
}
