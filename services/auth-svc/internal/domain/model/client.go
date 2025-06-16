package model

import "time"

type Client struct {
	ID        string
	Secret    string
	Roles     []string
	Status    int16
	CreatedAt time.Time
}
