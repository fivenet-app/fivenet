// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/completor/completor.proto

package completor

func (m *CompleteCitizenLabelsRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *CompleteCitizenLabelsResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Labels
	for idx, item := range m.Labels {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *CompleteCitizensRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *CompleteCitizensResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Users
	for idx, item := range m.Users {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *CompleteDocumentCategoriesRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *CompleteDocumentCategoriesResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Categories
	for idx, item := range m.Categories {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *CompleteJobsRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *CompleteJobsResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Jobs
	for idx, item := range m.Jobs {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *ListLawBooksRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *ListLawBooksResponse) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Books
	for idx, item := range m.Books {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	return nil
}
