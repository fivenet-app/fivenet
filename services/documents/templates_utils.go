package documents

import (
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
)

func (s *Server) checkAccessAgainstTemplate(
	tmpl *documents.Template,
	docAccess *documents.DocumentAccess,
) bool {
	if tmpl.GetContentAccess() == nil {
		return true
	}

	for _, access := range tmpl.GetContentAccess().GetJobs() {
		if access.Required == nil || !access.GetRequired() {
			continue
		}

		if !slices.ContainsFunc(docAccess.GetJobs(), func(ja *documents.DocumentJobAccess) bool {
			return ja.GetJob() == access.GetJob() &&
				ja.GetMinimumGrade() == access.GetMinimumGrade() &&
				ja.GetAccess() == access.GetAccess()
		}) {
			return false
		}
	}

	for _, access := range tmpl.GetContentAccess().GetUsers() {
		if access.Required == nil || !access.GetRequired() {
			continue
		}

		if !slices.ContainsFunc(docAccess.GetUsers(), func(ja *documents.DocumentUserAccess) bool {
			return ja.GetUserId() == access.GetUserId() && ja.GetAccess() == access.GetAccess()
		}) {
			return false
		}
	}

	return true
}
