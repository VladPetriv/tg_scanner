package errors

import "fmt"

type CreateError struct {
	Name       string
	ErrorValue error
}

func (c *CreateError) Error() string {
	return fmt.Sprintf("create %s error: %s", c.Name, c.ErrorValue)
}

type GetError struct {
	Name       string
	ErrorValue error
}

func (g *GetError) Error() string {
	return fmt.Sprintf("get %s error: %s", g.Name, g.ErrorValue)
}
