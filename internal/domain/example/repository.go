// name of the package
package example

import "context"

/**************************************************
* This file constains domain functionality to be  *
* implemented.                                    *
***************************************************/

// IdentityProvider provides new Identifier types
type IdentityProvider interface {
	NewID() Identifier
	ParseID(string) (Identifier, error)
}

type LineRepository interface {
	Write(context.Context, Line) error
	Read(context.Context, Identifier) (*Line, error)
}
