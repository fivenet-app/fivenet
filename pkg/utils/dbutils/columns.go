package dbutils

import (
	jet "github.com/go-jet/jet/v2/mysql"
)

const DisableColumnName = "-"

type CustomColumns struct {
	User    UserColumns    `yaml:"user"`
	Vehicle VehicleColumns `yaml:"vehicle"`
}

type UserColumns struct {
	Visum    string `default:"visum" yaml:"visum"`
	Playtime string `default:"playtime" yaml:"playtime"`
}

func (c *UserColumns) GetVisum(alias string) jet.Projection {
	if c.Visum == DisableColumnName {
		return nil
	}
	return jet.RawInt(alias + "." + c.Visum).AS(alias + ".visum")
}

func (c *UserColumns) GetPlaytime(alias string) jet.Projection {
	if c.Playtime == DisableColumnName {
		return nil
	}
	return jet.RawInt(alias + "." + c.Playtime).AS(alias + ".playtime")
}

type VehicleColumns struct {
	Model string `default:"model" yaml:"model"`
}

func (c *VehicleColumns) GetModel(alias string) jet.Projection {
	if c.Model == DisableColumnName {
		return nil
	}
	return jet.RawInt(alias + "." + c.Model).AS(alias + ".model")
}
