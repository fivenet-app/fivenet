package content

func (dr *ContentDiff) HasChanges() bool {
	if dr.GetStats() == nil {
		return false
	}
	return dr.GetStats().GetInsertedRunes() > 0 || dr.GetStats().GetDeletedRunes() > 0
}
