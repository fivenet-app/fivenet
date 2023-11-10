package dispatch

import (
	"database/sql/driver"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (x *Attributes) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf Attributes.
func (x *Attributes) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}

func (x *Dispatch) Merge(in *Dispatch) {
	if x.Id != in.Id {
		return
	}

	if in.CreatedAt != nil {
		if x.CreatedAt == nil {
			x.CreatedAt = in.CreatedAt
		} else {
			proto.Merge(x.CreatedAt, in.CreatedAt)
		}
	}

	if in.UpdatedAt != nil {
		if x.UpdatedAt == nil {
			x.UpdatedAt = in.UpdatedAt
		} else {
			proto.Merge(x.UpdatedAt, in.UpdatedAt)
		}
	}

	if x.Job != in.Job {
		x.Job = in.Job
	}

	if in.Status != nil {
		if x.Status == nil {
			x.Status = in.Status
		} else {
			proto.Merge(x.Status, in.Status)
		}
	}

	if x.Message != in.Message {
		x.Message = in.Message
	}

	if in.Description != nil && (x.Description == nil || x.Description != in.Description) {
		x.Description = in.Description
	}

	if in.Attributes != nil {
		if x.Attributes == nil {
			x.Attributes = in.Attributes
		} else {
			x.Attributes.List = in.Attributes.List
		}
	}

	if x.X != in.X {
		x.X = in.X
	}

	if x.Y != in.Y {
		x.Y = in.Y
	}

	if in.Postal != nil && (x.Postal == nil || x.Postal != in.Postal) {
		x.Postal = in.Postal
	}

	if x.Anon != in.Anon {
		x.Anon = in.Anon
	}

	if in.CreatorId != nil && (x.CreatorId == nil || x.CreatorId != in.CreatorId) {
		x.CreatorId = in.CreatorId
	}

	if in.Creator != nil {
		if x.Creator == nil {
			x.Creator = in.Creator
		} else {
			proto.Merge(x.Creator, in.Creator)
		}
	}

	for _, unit := range x.Units {
		fmt.Println("MERGE - Dispatch Units", x.Id, unit.Unit.Name, "length", len(x.Units))
	}
	x.Units = in.Units
}

func (x *Dispatch) Point() orb.Point {
	return orb.Point{x.X, x.Y}
}

func (x *DispatchStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{*x.X, *x.Y}
}
