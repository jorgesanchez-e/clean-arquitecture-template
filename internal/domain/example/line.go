package example

import (
	"time"
)

/**************************************************
* This file constains entities (aka data) used to *
* implement our business domain.                  *
***************************************************/

type Identifier interface {
	String() string
}

type Line struct {
	ID      Identifier
	Created time.Time
	Data    string
}
