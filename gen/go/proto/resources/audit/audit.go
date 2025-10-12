package audit

func (x *AuditEntryMeta) Set(key string, value string) {
	if x.Meta == nil {
		x.Meta = make(map[string]string)
	}

	x.Meta[key] = value
}
