package documents

func (x *DocUpdated) HasChanges() bool {
	if x == nil {
		return false
	}

	if x.GetTitleCdiff() != nil ||
		x.GetContentCdiff() != nil {
		return true
	}

	if x.GetTitleDiff() != "" ||
		x.GetContentDiff() != "" {
		return true
	}

	if !x.GetFilesChange().IsEmpty() {
		return true
	}

	return false
}

func (x *DocFilesChange) IsEmpty() bool {
	if x == nil {
		return true
	}

	return x.GetAdded() == 0 && x.GetDeleted() == 0
}
