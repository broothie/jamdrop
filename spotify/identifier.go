package spotify

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func IDFromIdentifier(identifier string) (string, error) {
	id := identifier
	if strings.Contains(identifier, "/") {
		u, err := url.Parse(identifier)
		if err != nil {
			return "", errors.Wrap(err, "invalid user link")
		}

		segments := strings.Split(u.Path, "/")
		id = segments[len(segments)-1]
	} else if strings.Contains(identifier, ":") {
		segments := strings.Split(identifier, ":")
		id = segments[len(segments)-1]
	}

	return id, nil
}

func SongURI(id string) string {
	return TrackURI(id)
}

func TrackURI(id string) string {
	return URI("track", id)
}

func UserURI(id string) string {
	return URI("user", id)
}

func URI(entity, id string) string {
	return fmt.Sprintf("spotify:%s:%s", entity, id)
}
