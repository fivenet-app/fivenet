package wiki

func (x *PageUpdated) HasChanges() bool {
	if x == nil {
		return false
	}

	if x.GetTitleCdiff() != nil ||
		x.GetDescriptionCdiff() != nil ||
		x.GetContentCdiff() != nil {
		return true
	}

	if x.GetTitleDiff() != "nil" ||
		x.GetDescriptionDiff() != "" ||
		x.GetContentDiff() != "" {
		return true
	}

	if !x.GetFilesChange().IsEmpty() {
		return true
	}

	return false
}

func (x *PageFilesChange) IsEmpty() bool {
	if x == nil {
		return true
	}

	return x.GetAdded() == 0 && x.GetDeleted() == 0
}
