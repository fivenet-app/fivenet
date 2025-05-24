package documents

import (
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
)

func (s *Server) checkAccessAgainstTemplate(tmpl *documents.Template, docAccess *documents.DocumentAccess) bool {
	if tmpl.ContentAccess == nil {
		return true
	}

	for _, access := range tmpl.ContentAccess.Jobs {
		if access.Required == nil || !*access.Required {
			continue
		}

		if !slices.ContainsFunc(docAccess.Jobs, func(ja *documents.DocumentJobAccess) bool {
			return ja.Job == access.Job && ja.MinimumGrade == access.MinimumGrade && ja.Access == access.Access
		}) {
			return false
		}
	}

	for _, access := range tmpl.ContentAccess.Users {
		if access.Required == nil || !*access.Required {
			continue
		}

		if !slices.ContainsFunc(docAccess.Users, func(ja *documents.DocumentUserAccess) bool {
			return ja.UserId == access.UserId && ja.Access == access.Access
		}) {
			return false
		}
	}

	return true
}
