package register

import "github.com/google/uuid"

/**************************************************
* This file constains domain functionality to be  *
* implemented.                                    *
***************************************************/

type LineRepository interface {
	Write(Line) error
	Read(id uuid.UUID) *Line
}
