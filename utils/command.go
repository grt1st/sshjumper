package utils

import "strings"

type Command struct {
	Name   string            // command
	Params map[string]string // args
	Peaces []string
}

func NewCommand(input string) (c Command) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	for k, v := range parts {
		if k == 0 {
			c.Name = parts[0]
			continue
		}
		if strings.HasPrefix(v, "--") {
			// eg: --ignore
			if c.Params == nil {
				c.Params = make(map[string]string, 0)
			}
			c.Params[strings.TrimPrefix(v, "--")] = "true"
		} else if strings.HasPrefix(v, "-") {
			// eg: -l、-p 25
			if c.Params == nil {
				c.Params = make(map[string]string, 0)
			}
			value := "true"
			if k+1 < len(parts) && !strings.HasPrefix(parts[k+1], "-") {
				value = parts[k+1]
			}
			c.Params[strings.TrimPrefix(v, "-")] = value
		} else {
			if v != "" {
				c.Peaces = append(c.Peaces, v)
			}
		}
	}
	return
}
