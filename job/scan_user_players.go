package job

import (
	"context"
	"sync"
	"time"

	"jamdrop/model"

	"cloud.google.com/go/firestore"
)

func (j *Job) ScanUserPlayers(ctx context.Context) error {
	docs, err := j.DB.Collection(model.CollectionUsers).Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	now := time.Now()
	batch := j.DB.Batch()
	var wg sync.WaitGroup
	docChan := make(chan *firestore.DocumentSnapshot)
	for i := 0; i < j.App.Config.ScanWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for doc := range docChan {
				user := new(model.User)
				if err := doc.DataTo(user); err != nil {
					j.Logger.Printf("failed to read user data; user_id: %s; %v\n", doc.Ref.ID, err)
					continue
				}

				isPlaying, err := j.Spotify.GetCurrentlyPlaying(user)
				if err != nil {
					j.Logger.Printf("failed to get currently playing: user_id: %s, access_token: %s; %v\n", user.ID, user.AccessToken, err)
					time.Sleep(10 * time.Millisecond)
					continue
				}

				if isPlaying {
					batch.Update(doc.Ref, []firestore.Update{{Path: "last_playing", Value: now}})
				}
			}
		}()
	}

	for _, doc := range docs {
		docChan <- doc
	}

	close(docChan)
	wg.Wait()

	if _, err := batch.Commit(ctx); err != nil {
		j.Logger.Println("error committing batch", err)
		return err
	}

	return nil
}
