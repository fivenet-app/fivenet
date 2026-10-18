package centrumunits

// Merge replaces each independently versioned unit component in place. The
// store keeps pointer identity for local-cache consumers, while each KV value
// remains an authoritative complete snapshot of that component.
func (x *UnitSpec) Merge(in *UnitSpec) *UnitSpec {
	if x == nil || in == nil || x.GetId() != in.GetId() {
		return x
	}

	x.CreatedAt = in.CreatedAt
	x.UpdatedAt = in.UpdatedAt
	x.Job = in.Job
	x.JobLabel = in.JobLabel
	x.SortOrder = in.SortOrder
	x.Name = in.Name
	x.Initials = in.Initials
	x.Color = in.Color
	x.Icon = in.Icon
	x.Description = in.Description
	x.Attributes = in.Attributes
	x.HomePostal = in.HomePostal
	x.Access = in.Access
	return x
}

func (x *UnitAssignments) Merge(in *UnitAssignments) *UnitAssignments {
	if x == nil || in == nil || x.GetUnitId() != in.GetUnitId() {
		return x
	}

	x.Job = in.Job
	x.Users = in.Users
	return x
}

func (x *UnitStatus) Merge(in *UnitStatus) *UnitStatus {
	if x == nil || in == nil {
		return x
	}

	x.Id = in.Id
	x.CreatedAt = in.CreatedAt
	x.UnitId = in.UnitId
	x.Unit = in.Unit
	x.Status = in.Status
	x.Reason = in.Reason
	x.Code = in.Code
	x.UserId = in.UserId
	x.User = in.User
	x.X = in.X
	x.Y = in.Y
	x.Postal = in.Postal
	x.CreatorId = in.CreatorId
	x.Creator = in.Creator
	x.CreatorJob = in.CreatorJob
	return x
}
