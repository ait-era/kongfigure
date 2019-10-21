package kongfigure

import "fmt"

type KongfigureHttpError struct {
	CustomMessage string
	ResponseUrl   string
	ResponseError string
}

func (e *KongfigureHttpError) Error() string {
	return fmt.Sprintf("%v | <%s> | Response data: <%v>", e.CustomMessage, e.ResponseUrl, e.ResponseError)
}

