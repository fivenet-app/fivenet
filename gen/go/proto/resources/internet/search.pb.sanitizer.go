// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/internet/search.proto

package internet

func (m *SearchResult) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Domain
	if m.Domain != nil {
		if v, ok := any(m.GetDomain()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
