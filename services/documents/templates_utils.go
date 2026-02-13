package documents

import (
	"slices"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
)

func (s *Server) checkAccessAgainstTemplate(
	tmpl *documentstemplates.Template,
	docAccess *documentsaccess.DocumentAccess,
) bool {
	if tmpl.GetContentAccess() == nil {
		return true
	}

	for _, access := range tmpl.GetContentAccess().GetJobs() {
		if access.Required == nil || !access.GetRequired() {
			continue
		}

		if !slices.ContainsFunc(
			docAccess.GetJobs(),
			func(ja *documentsaccess.DocumentJobAccess) bool {
				return ja.GetJob() == access.GetJob() &&
					ja.GetMinimumGrade() == access.GetMinimumGrade() &&
					ja.GetAccess() == access.GetAccess()
			},
		) {
			return false
		}
	}

	for _, access := range tmpl.GetContentAccess().GetUsers() {
		if access.Required == nil || !access.GetRequired() {
			continue
		}

		if !slices.ContainsFunc(
			docAccess.GetUsers(),
			func(ja *documentsaccess.DocumentUserAccess) bool {
				return ja.GetUserId() == access.GetUserId() && ja.GetAccess() == access.GetAccess()
			},
		) {
			return false
		}
	}

	return true
}
