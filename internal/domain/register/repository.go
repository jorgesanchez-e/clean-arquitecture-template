package register

/**************************************************
* This file constains domain functionality to be  *
* implemented.                                    *
***************************************************/

type LineRepository interface {
	Write(Line) error
	Read(id Identifier) *Line
}
