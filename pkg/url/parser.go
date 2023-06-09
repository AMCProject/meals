package url

import "github.com/labstack/echo/v4"

type PathMap map[string]struct {
	Target *string
	Err    error
}

func ParseURLPath(c echo.Context, paramsMap PathMap) error {
	for name, mapValue := range paramsMap {
		value := c.Param(name)
		if value == "" {
			return mapValue.Err
		}
		*mapValue.Target = value
	}
	return nil
}
