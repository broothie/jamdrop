package job

import (
	"context"
	"jamdrop/logger"
	"time"

	"jamdrop/model"
)

func (j *Job) EjectSessionTokens(ctx context.Context) error {
	j.Logger.Info("job.EjectSessionTokens")

	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	j.Logger.Info("deleting session tokens", logger.Fields{"updated_at before": thirtyDaysAgo})

	docs, err := j.DB.
		Collection(model.CollectionSessionTokens).
		Where("updated_at", "<", thirtyDaysAgo).
		Documents(ctx).
		GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		j.Logger.Info("no expired session tokens found")
		return nil
	}

	batch := j.DB.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	j.Logger.Info("deleting session tokens", logger.Fields{"n": len(docs)})
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}

	return nil
}
