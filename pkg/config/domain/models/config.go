package models

import "errors"

type Config struct {
	Servers []Server `yaml:"servers"`
}

func (c *Config) Validate() error {
	if len(c.Servers) == 0 {
		return errors.New("No server found")
	}
	for _, v := range c.Servers {
		err := v.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
