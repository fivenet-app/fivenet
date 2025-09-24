package dbutils

import (
	"github.com/go-jet/jet/v2/mysql"
)

type Columns mysql.ProjectionList

func (c Columns) Get() mysql.ProjectionList {
	out := mysql.ProjectionList{}

	for i := range c {
		if c[i] != nil {
			out = append(out, c[i])
		}
	}

	return out
}

const DisableColumnName = "-"

type CustomColumns struct {
	User    UserColumns    `yaml:"user"`
	Vehicle VehicleColumns `yaml:"vehicle"`
}

type UserColumns struct {
	Visum    string `default:"visum"    yaml:"visum"`
	Playtime string `default:"playtime" yaml:"playtime"`
}

func (c *UserColumns) GetVisum(alias string) mysql.Projection {
	if c.Visum == DisableColumnName {
		return nil
	}
	return mysql.RawInt(alias + "." + c.Visum).AS(alias + ".visum")
}

func (c *UserColumns) GetPlaytime(alias string) mysql.Projection {
	if c.Playtime == DisableColumnName {
		return nil
	}
	return mysql.RawInt(alias + "." + c.Playtime).AS(alias + ".playtime")
}

type VehicleColumns struct {
	Model string `default:"model" yaml:"model"`
}

func (c *VehicleColumns) GetModel(alias string) mysql.Projection {
	if c.Model == DisableColumnName {
		return nil
	}
	return mysql.RawInt(alias + "." + c.Model).AS(alias + ".model")
}
