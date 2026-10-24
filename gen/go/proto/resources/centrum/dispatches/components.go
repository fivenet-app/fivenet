package centrumdispatches

// Merge replaces a dispatch specification snapshot. Component stores contain
// complete snapshots, so field-wise protobuf merge semantics are undesirable.
func (x *DispatchSpec) Merge(in *DispatchSpec) *DispatchSpec {
	if x == nil || in == nil || x.GetId() != in.GetId() {
		return x
	}

	x.CreatedAt = in.CreatedAt
	x.UpdatedAt = in.UpdatedAt
	x.Jobs = in.Jobs
	x.Message = in.Message
	x.Description = in.Description
	x.Attributes = in.Attributes
	x.X = in.X
	x.Y = in.Y
	x.Postal = in.Postal
	x.Anon = in.Anon
	x.CreatorId = in.CreatorId
	x.Creator = in.Creator
	x.References = in.References
	return x
}

// Merge replaces the complete assignment snapshot.
func (x *DispatchAssignments) Merge(in *DispatchAssignments) *DispatchAssignments {
	if x == nil || in == nil || x.GetDispatchId() != in.GetDispatchId() {
		return x
	}

	x.Job = in.Job
	x.Units = in.Units
	return x
}

// Merge replaces the complete status snapshot.
func (x *DispatchStatus) Merge(in *DispatchStatus) *DispatchStatus {
	if x == nil || in == nil || x.GetDispatchId() != in.GetDispatchId() {
		return x
	}

	x.Id = in.Id
	x.CreatedAt = in.CreatedAt
	x.DispatchId = in.DispatchId
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
	x.CreatorJob = in.CreatorJob
	return x
}
