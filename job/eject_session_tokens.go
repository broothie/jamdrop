package job

import (
	"context"
	"jamdrop/logger"
	"jamdrop/requestid"
	"time"

	"jamdrop/model"
)

func (j *Job) EjectSessionTokens(ctx context.Context) error {
	j.Logger.Debug("job.EjectSessionTokens", requestid.LogContext(ctx))

	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	j.Logger.Info("deleting session tokens", logger.Fields{"updated_at before": thirtyDaysAgo}, requestid.LogContext(ctx))

	docs, err := j.DB.
		Collection(model.CollectionSessionTokens).
		Where("updated_at", "<", thirtyDaysAgo).
		Documents(ctx).
		GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		j.Logger.Info("no expired session tokens found", requestid.LogContext(ctx))
		return nil
	}

	batch := j.DB.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	j.Logger.Info("deleting session tokens", logger.Fields{"n": len(docs)}, requestid.LogContext(ctx))
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}

	return nil
}
