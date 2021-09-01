package contagions

import "time"

type Contagion struct {
	UserGeneratedCode string    `json:"userGeneratedCode"`
	SpaceId           string    `json:"spaceId,omitempty"`
	Timestamp         time.Time `json:"timestamp,omitempty"`
}
