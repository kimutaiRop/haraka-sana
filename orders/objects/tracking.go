package objects

type CreateTrackingStep struct {
	Name           string `json:"name"`
	StreamLocation bool   `json:"stream_location"`
}
