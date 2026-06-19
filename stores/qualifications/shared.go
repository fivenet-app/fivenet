package qualificationsstore

import "github.com/fivenet-app/fivenet/v2026/query/fivenet/table"

const (
	QualificationsPageSize = 10

	QualificationsLabelDefaultFormat = "%abbr%: %name%"
)

var (
	tQuali              = table.FivenetQualifications.AS("qualification")
	tQualiReqs          = table.FivenetQualificationsRequirements.AS("qualification_requirement")
	tQualiResult        = table.FivenetQualificationsResults.AS("qualification_result")
	tQualiResultSuccess = table.FivenetQualificationsResultSuccessMap.AS(
		"qualification_result_success_map",
	)
	tQualiReq      = table.FivenetQualificationsRequests.AS("qualification_request")
	tExamQuestion  = table.FivenetQualificationsExamQuestions.AS("exam_question")
	tExamResponses = table.FivenetQualificationsExamResponses.AS("exam_response")
	tExamUser      = table.FivenetQualificationsExamUsers.AS("exam_user")
)
