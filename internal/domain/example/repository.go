package example

import "context"

/**************************************************
* This file constains domain functionality to be  *
* implemented.                                    *
***************************************************/

type LineRepository interface {
	Write(context.Context, Line) error
	Read(context.Context, Identifier) (*Line, error)
}
