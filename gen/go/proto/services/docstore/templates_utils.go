package docstore

import (
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
)

func (s *Server) checkAccessAgainstTemplate(ctx context.Context, id uint64, docAccess *documents.DocumentAccess) (bool, error) {
	tmpl, err := s.getTemplate(ctx, id)
	if err != nil {
		return false, err
	}

	if tmpl.ContentAccess == nil {
		return true, nil
	}

	for _, access := range tmpl.ContentAccess.Jobs {
		if access.Required == nil || !*access.Required {
			continue
		}

		if !slices.ContainsFunc(docAccess.Jobs, func(ja *documents.DocumentJobAccess) bool {
			return ja.Job == access.Job && ja.MinimumGrade == access.MinimumGrade && ja.Access == access.Access
		}) {
			return false, nil
		}
	}

	for _, access := range tmpl.ContentAccess.Users {
		if access.Required == nil || !*access.Required {
			continue
		}

		if !slices.ContainsFunc(docAccess.Users, func(ja *documents.DocumentUserAccess) bool {
			return ja.UserId == access.UserId && ja.Access == access.Access
		}) {
			return false, nil
		}
	}

	return true, nil
}
