package spotify

import "fmt"

type ErrorReason string

const (
	noActiveDevice ErrorReason = "NO_ACTIVE_DEVICE"
)

type errorResponse struct {
	Error SpotifyError `json:"error"`
}

type SpotifyError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Reason  ErrorReason `json:"reason"`
}

func (e SpotifyError) Error() string {
	return fmt.Sprintf("spotify error: %d %s %s", e.Status, e.Reason, e.Message)
}
