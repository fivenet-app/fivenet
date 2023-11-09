package dispatch

import (
	"fmt"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func (x *Unit) Merge(in *Unit) {
	if x.Id != in.Id {
		return
	}

	if in.CreatedAt != nil {
		proto.Merge(x.CreatedAt, in.CreatedAt)
	}

	if in.UpdatedAt != nil {
		proto.Merge(x.UpdatedAt, in.UpdatedAt)
	}

	if x.Job != in.Job {
		x.Job = in.Job
	}

	if x.Name != in.Name {
		x.Name = in.Name
	}

	if x.Initials != in.Initials {
		x.Initials = in.Initials
	}

	if x.Color != in.Color {
		x.Color = in.Color
	}

	if in.Description != nil && x.Description != in.Description {
		x.Description = in.Description
	}

	if in.Status != nil {
		proto.Merge(x.Status, in.Status)
	}

	for _, user := range x.Users {
		fmt.Println("MERGE - Unit Users", x.Id, user.UnitId, user.UserId, "length", len(x.Users))
	}
	x.Users = in.Users
}

func (x *UnitStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{*x.X, *x.Y}
}
