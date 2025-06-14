// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/filestore/filestore.proto

package filestore

func (m *DeleteFileByPathRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *DeleteFileByPathResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *ListFilesRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := any(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ListFilesResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Files
	for idx, item := range m.Files {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Pagination
	if m.Pagination != nil {
		if v, ok := any(m.GetPagination()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
