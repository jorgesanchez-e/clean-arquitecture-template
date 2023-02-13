package register

import (
	"time"

	"github.com/google/uuid"
)

/**************************************************
* This file constains entities (aka data) used to *
* implement our business domain.                  *
***************************************************/

type Line struct {
	id      uuid.UUID
	created time.Time
	data    string
}

func NewLine(data string) Line {
	return Line{
		id:      uuid.New(),
		created: time.Now(),
		data:    data,
	}
}
