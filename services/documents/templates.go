package documents

import (
	"bytes"
	context "context"
	"errors"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig/v3"
	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	colleagueshydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	htmlnode "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	kindKey      = "kind"
	requiredKey  = "required"
	availableKey = "available"
)

var ErrTemplateActiveChar = errors.New("failed to resolve active character/user")

func isTemplateActionSpan(node *htmlnode.Node) bool {
	if node.Type != htmlnode.ElementNode || node.Data != "span" {
		return false
	}

	for _, attr := range node.Attr {
		switch attr.Key {
		case "data-template-var", "data-template-block", "data-template-block-end":
			return true
		}
	}

	return false
}

func unwrapTemplateActionSpans(node *htmlnode.Node) {
	for child := node.FirstChild; child != nil; {
		next := child.NextSibling
		unwrapTemplateActionSpans(child)

		if isTemplateActionSpan(child) {
			for content := child.FirstChild; content != nil; {
				nextContent := content.NextSibling
				child.RemoveChild(content)
				node.InsertBefore(content, child)
				content = nextContent
			}
			node.RemoveChild(child)
		}

		child = next
	}
}

func stripTemplateActionSpans(content string) (string, error) {
	// Parse the editor HTML as a fragment so action spans can be identified by
	// their element and exact data-* attributes instead of matching HTML with a
	// regular expression. Rendering the cleaned tree may normalize markup such
	// as attribute quoting or equivalent whitespace, which is intentional.
	fragment, err := htmlnode.ParseFragment(strings.NewReader(content), &htmlnode.Node{
		Type:     htmlnode.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	})
	if err != nil {
		return "", err
	}

	root := &htmlnode.Node{Type: htmlnode.DocumentNode}
	for _, node := range fragment {
		root.AppendChild(node)
	}
	unwrapTemplateActionSpans(root)

	var output bytes.Buffer
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		if err := htmlnode.Render(&output, node); err != nil {
			return "", err
		}
	}

	return output.String(), nil
}

var templateSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

func (s *Server) sanitizeTemplateAccess(
	tmpl *documentstemplates.Template,
	sanitizeJobAccess bool,
	sanitizeContentAccess bool,
) error {
	if tmpl == nil {
		return nil
	}

	if sanitizeJobAccess {
		jobAccess, err := access.SanitizeJobAccessEntries(s.jobs, tmpl.GetJobAccess())
		if err != nil {
			return err
		}
		tmpl.JobAccess = access.NormalizeRequiredJobAccessFloors(jobAccess)
	}

	if sanitizeContentAccess {
		contentAccess, err := access.SanitizeAccessJobs(s.jobs, tmpl.GetContentAccess())
		if err != nil {
			return err
		}
		tmpl.ContentAccess = access.NormalizeRequiredAccessFloors(contentAccess)
	}

	return nil
}

func templateJobAccess(jobs []*documentstemplates.TemplateJobAccess) *resourcesaccess.Access {
	return &resourcesaccess.Access{Jobs: jobs}
}

func (s *Server) normalizeTemplateJobAccess(
	userInfo *userinfo.UserInfo,
	jobs []*documentstemplates.TemplateJobAccess,
) (*resourcesaccess.Access, error) {
	highestGrade := userInfo.GetJobGrade()
	if grade, ok := s.enricher.GetHighestJobGrade(userInfo.GetJob()); ok {
		highestGrade = grade
	}

	return access.NormalizeAccess(
		templateJobAccess(jobs),
		nil,
		&resourcesaccess.Access{
			Jobs: []*resourcesaccess.JobAccess{{
				Job:          userInfo.GetJob(),
				MinimumGrade: highestGrade,
				Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			}},
		},
		15,
	)
}

func (s *Server) ListTemplates(
	ctx context.Context,
	req *pbdocuments.ListTemplatesRequest,
) (*pbdocuments.ListTemplatesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	templates, err := s.store.ListTemplates(ctx, false, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	resp := &pbdocuments.ListTemplatesResponse{
		Templates: templates,
	}

	for i := range resp.GetTemplates() {
		s.enricher.EnrichJobName(resp.GetTemplates()[i])
	}

	return resp, nil
}

func (s *Server) GetTemplate(
	ctx context.Context,
	req *pbdocuments.GetTemplateRequest,
) (*pbdocuments.GetTemplateResponse, error) {
	logging.InjectFields(ctx, logging.Fields{templateIDLogFieldKey, req.GetTemplateId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplateId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	resp := &pbdocuments.GetTemplateResponse{}
	resp.Template, err = s.store.GetTemplate(ctx, req.GetTemplateId(), userInfo.GetJobAdmin())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.GetTemplate() == nil {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if req.Render == nil || !req.GetRender() {
		templateAccess, err := s.templateAccess.ListTargetAccess(
			ctx,
			s.db,
			req.GetTemplateId(),
			templateSubjectAccessOptions,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.Template.JobAccess = templateAccess.GetJobs()
		if err := s.sanitizeTemplateAccess(resp.GetTemplate(), true, true); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else if req.Render != nil && req.GetRender() && req.GetSelection() != nil {
		data, err := s.resolveTemplateData(ctx, resp.GetTemplate(), req.GetSelection(), userInfo)
		if err != nil {
			return nil, err
		}
		resp.Template.ContentTitle, resp.Template.State, resp.Template.Content, err = s.renderTemplate(
			resp.GetTemplate(),
			data,
		)
		if err != nil {
			if s.perms.Can(
				userInfo,
				permsdocuments.TemplatesService.CreateTemplate.Perm,
			) {
				return nil, err
			} else {
				return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateInvalid)
			}
		}

		resp.Rendered = true
	}

	s.enricher.EnrichJobName(resp.GetTemplate())

	return resp, nil
}

type resolvedTemplateData struct {
	ActiveChar *jobscolleagues.Colleague
	Documents  []*resourcesdocuments.DocumentShort
	Users      []*users.User
	Vehicles   []*resourcesvehicles.Vehicle
}

func (s *Server) resolveTemplateData(
	ctx context.Context,
	tmpl *documentstemplates.Template,
	selection *documentstemplates.TemplateSelection,
	userInfo *userinfo.UserInfo,
) (*resolvedTemplateData, error) {
	activeChar, err := s.colleagueHydrator.GetBasicByUserID(
		ctx,
		s.db,
		userInfo,
		userInfo.GetUserId(),
		colleagueshydrator.ResolveOpts{},
	)
	if err != nil {
		return nil, err
	}
	if activeChar == nil {
		return nil, errswrap.NewError(
			ErrTemplateActiveChar,
			errorsdocuments.ErrTemplateRenderFailed,
		)
	}

	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateRenderFailed)
	}

	data := &resolvedTemplateData{
		ActiveChar: activeChar,
	}
	if selection == nil {
		return data, validateTemplateRequirements(tmpl, data)
	}

	if len(selection.GetUserIds()) > 0 {
		pageSize := int64(len(selection.GetUserIds()))
		citizensReq := &pbcitizens.ListCitizensRequest{Pagination: &database.PaginationRequest{}}
		citizensReq.GetPagination().SetPageSize(pageSize)

		listOptions := citizensListOptions(fields)
		listOptions.UserIDs = selection.GetUserIds()

		usersResp, err := s.citizensStore.ListCitizens(ctx, citizensReq, listOptions)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateRenderFailed)
		}

		jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
		jobInfoFn(activeChar)
		for _, user := range usersResp.GetUsers() {
			jobInfoFn(user)
		}

		byID := make(map[int32]*users.User, len(usersResp.GetUsers()))
		for _, user := range usersResp.GetUsers() {
			byID[user.GetUserId()] = user
		}
		seen := make(map[int32]struct{}, len(selection.GetUserIds()))
		for _, id := range selection.GetUserIds() {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			if user := byID[id]; user != nil {
				data.Users = append(data.Users, user)
			}
		}
	}

	if len(selection.GetDocumentIds()) > 0 {
		docs, err := s.store.List(ctx, documentsstore.ListQuery{
			DocumentIDs: selection.GetDocumentIds(),
			Limit:       int64(len(selection.GetDocumentIds())),
			UserInfo:    userInfo,
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateRenderFailed)
		}

		byID := make(map[int64]*resourcesdocuments.DocumentShort, len(docs))
		for _, doc := range docs {
			byID[doc.GetId()] = doc
		}
		seen := make(map[int64]struct{}, len(selection.GetDocumentIds()))
		for _, id := range selection.GetDocumentIds() {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			if doc := byID[id]; doc != nil {
				data.Documents = append(data.Documents, doc)
			}
		}
	}

	vehicleFields, err := permsvehicles.VehiclesService.ListVehicles.FieldsTyped.Get(
		s.perms,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateRenderFailed)
	}
	if len(selection.GetPlates()) > 0 {
		vehicles, err := s.vehiclesStore.List(ctx, vehiclesstore.ListQuery{
			Plates:              selection.GetPlates(),
			Limit:               int64(len(selection.GetPlates())),
			IncludePropsUpdated: vehicleFields.Len() > 0,
			IncludeWantedFields: vehicleFields.Contains(
				permsvehicles.VehiclesServiceListVehiclesFieldsPermValueWanted,
			) ||
				userInfo.GetJobAdmin(),
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateRenderFailed)
		}

		byPlate := make(map[string]*resourcesvehicles.Vehicle, len(vehicles))
		for _, vehicle := range vehicles {
			byPlate[vehicle.GetPlate()] = vehicle
		}
		seen := make(map[string]struct{}, len(selection.GetPlates()))
		for _, plate := range selection.GetPlates() {
			if _, ok := seen[plate]; ok {
				continue
			}
			seen[plate] = struct{}{}
			if vehicle := byPlate[plate]; vehicle != nil {
				data.Vehicles = append(data.Vehicles, vehicle)
			}
		}
	}

	return data, validateTemplateRequirements(tmpl, data)
}

func citizensListOptions(
	fields *perms.TypedStringList[permscitizens.CitizensServiceListCitizensFieldsPermValue],
) citizensstore.ListCitizensOptions {
	return citizensstore.ListCitizensOptions{
		IncludePhoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		IncludeWanted: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsWanted,
		),
		IncludeJob: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsJob,
		),
		IncludeTrafficInfractionPoints: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsTrafficInfractionPoints,
		),
		IncludeOpenFines: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsOpenFines,
		),
		IncludeBloodType: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType,
		),
		IncludeMugshot: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot,
		),
		IncludeEmail: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsEmail,
		),
	}
}

func validateTemplateRequirements(
	tmpl *documentstemplates.Template,
	data *resolvedTemplateData,
) error {
	reqs := tmpl.GetSchema().GetRequirements()
	checks := []struct {
		spec  *documentstemplates.ObjectSpecs
		kind  string
		count int
	}{
		{reqs.GetUsers(), "users", len(data.Users)},
		{reqs.GetDocuments(), "documents", len(data.Documents)},
		{reqs.GetVehicles(), "vehicles", len(data.Vehicles)},
	}
	for _, check := range checks {
		if check.spec == nil {
			continue
		}
		if check.spec.GetRequired() && check.count == 0 {
			return errorsdocuments.ErrTemplateRequirementsNotMet(map[string]any{
				kindKey:      check.kind,
				requiredKey:  1,
				availableKey: check.count,
			})
		}
		if check.spec.GetMin() > 0 && int64(check.count) < int64(check.spec.GetMin()) {
			return errorsdocuments.ErrTemplateRequirementsNotMet(map[string]any{
				kindKey:      check.kind,
				requiredKey:  check.spec.GetMin(),
				availableKey: check.count,
			})
		}
		if check.spec.GetMax() > 0 && int64(check.count) > int64(check.spec.GetMax()) {
			return errorsdocuments.ErrTemplateRequirementsExceeded(map[string]any{
				kindKey:      check.kind,
				requiredKey:  check.spec.GetMax(),
				availableKey: check.count,
			})
		}
	}

	return nil
}

func (s *Server) renderTemplate(
	docTmpl *documentstemplates.Template,
	data *resolvedTemplateData,
) (string, string, string, error) {
	// Render Title template
	titleTpl, err := template.
		New("title").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.GetContentTitle())
	if err != nil {
		return "", "", "", err
	}
	buf := &bytes.Buffer{}
	err = titleTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	outTitle := buf.String()

	// Render State template
	stateTpl, err := template.
		New("state").
		Funcs(sprig.FuncMap()).
		Parse(docTmpl.GetState())
	if err != nil {
		return "", "", "", err
	}

	buf.Reset()
	err = stateTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	outState := buf.String()

	// Render Content template
	content, err := stripTemplateActionSpans(docTmpl.GetContent())
	if err != nil {
		return "", "", "", err
	}

	contentTpl, err := template.
		New("content").
		Funcs(sprig.FuncMap()).
		Parse(content)
	if err != nil {
		return "", "", "", err
	}

	buf.Reset()
	err = contentTpl.Execute(buf, data)
	if err != nil {
		return "", "", "", err
	}
	out := buf.String()

	return outTitle, outState, out, err
}

func (s *Server) CreateTemplate(
	ctx context.Context,
	req *pbdocuments.CreateTemplateRequest,
) (*pbdocuments.CreateTemplateResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.sanitizeTemplateAccess(req.GetTemplate(), true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if documentstemplates.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documentsaccess.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	var categoryId *int64
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			id := cat.GetId()
			categoryId = &id
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	sortRank, err := s.store.NextTemplateGroupRank(ctx, tx, userInfo.GetJob(), 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	lastId, err := s.store.CreateTemplate(
		ctx,
		tx,
		req.GetTemplate(),
		userInfo.GetJob(),
		categoryId,
		sortRank,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	templateAccess := templateJobAccess(req.GetTemplate().GetJobAccess())
	highestGrade := userInfo.GetJobGrade()
	if grade, ok := s.enricher.GetHighestJobGrade(userInfo.GetJob()); ok {
		highestGrade = grade
	}
	templateAccess = access.EnsureJobAccessEntries(
		templateAccess,
		&resourcesaccess.JobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		},
		&resourcesaccess.JobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: highestGrade,
			Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			Required:     new(true),
		},
	)
	normalizedAccess, err := s.normalizeTemplateJobAccess(userInfo, templateAccess.GetJobs())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		lastId,
		normalizedAccess,
		templateSubjectAccessOptions,
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbdocuments.CreateTemplateResponse{
		Id: lastId,
	}, nil
}

func (s *Server) UpdateTemplate(
	ctx context.Context,
	req *pbdocuments.UpdateTemplateRequest,
) (*pbdocuments.UpdateTemplateResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{templateIDLogFieldKey, req.GetTemplate().GetId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplate().GetId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if err := s.sanitizeTemplateAccess(req.GetTemplate(), true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if documentstemplates.TemplateAccessHasDuplicates(req.GetTemplate().GetJobAccess()) ||
		documentsaccess.DocumentAccessHasDuplicates(req.GetTemplate().GetContentAccess()) {
		return nil, errorsdocuments.ErrTemplateAccessDuplicate
	}

	var categoryId *int64
	if req.GetTemplate().GetCategory() != nil {
		cat, err := s.getCategory(ctx, req.GetTemplate().GetCategory().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if cat != nil {
			id := cat.GetId()
			categoryId = &id
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpdateTemplate(ctx, tx, req.GetTemplate(), categoryId); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	normalizedAccess, err := s.normalizeTemplateJobAccess(
		userInfo,
		req.GetTemplate().GetJobAccess(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := s.templateAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		req.GetTemplate().GetId(),
		normalizedAccess,
		templateSubjectAccessOptions,
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrTemplateAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	tmpl, err := s.store.GetTemplate(ctx, req.GetTemplate().GetId(), false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	templateAccess, err := s.templateAccess.ListTargetAccess(
		ctx,
		s.db,
		req.GetTemplate().GetId(),
		templateSubjectAccessOptions,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	tmpl.JobAccess = templateAccess.GetJobs()
	if err := s.sanitizeTemplateAccess(tmpl, true, true); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.UpdateTemplateResponse{
		Template: tmpl,
	}, nil
}

func (s *Server) DeleteTemplate(
	ctx context.Context,
	req *pbdocuments.DeleteTemplateRequest,
) (*pbdocuments.DeleteTemplateResponse, error) {
	logging.InjectFields(ctx, logging.Fields{templateIDLogFieldKey, req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	dTmpl, err := s.store.GetTemplate(ctx, req.GetId(), true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if !check && !userInfo.GetJobAdmin() {
		if dTmpl.GetCreatorJob() == "" {
			return nil, errorsdocuments.ErrTemplateNoPerms
		}

		// Make sure the highest job grade can delete the template
		grade := s.jobs.GetHighestJobGrade(userInfo.GetJob())
		if grade == nil ||
			(userInfo.GetJob() == dTmpl.GetCreatorJob() && grade.GetGrade() != userInfo.GetJobGrade()) {
			return nil, errorsdocuments.ErrTemplateNoPerms
		}
	}

	var deletedAtTime *timestamp.Timestamp
	// Check if page has any un-deleted child pages
	if dTmpl.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteTemplate(
		ctx,
		s.db,
		req.GetId(),
		userInfo.GetJob(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteTemplateResponse{}, nil
}

func (s *Server) MoveTemplate(
	ctx context.Context,
	req *pbdocuments.MoveTemplateRequest,
) (*pbdocuments.MoveTemplateResponse, error) {
	logging.InjectFields(ctx, logging.Fields{templateIDLogFieldKey, req.GetTemplateId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.templateAccess.CanUserAccessTarget(
		ctx,
		req.GetTemplateId(),
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	if req.GetBeforeId() > 0 && req.GetAfterId() > 0 {
		return nil, errorsdocuments.ErrFailedQuery
	}

	templateOrder, err := s.store.GetTemplateOrderInfo(ctx, s.db, req.GetTemplateId())
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorsdocuments.ErrTemplateNoPerms
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if templateOrder.CreatorJob != userInfo.GetJob() && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrTemplateNoPerms
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	sortRank, err := s.store.InsertTemplateGroupRank(
		ctx,
		tx,
		templateOrder.CreatorJob,
		req.GetTemplateId(),
		req.BeforeId,
		req.AfterId,
	)
	if err != nil {
		if errors.Is(err, errorsdocuments.ErrTemplateNoPerms) {
			return nil, err
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.store.UpdateTemplateSortRank(
		ctx,
		tx,
		req.GetTemplateId(),
		templateOrder.CreatorJob,
		sortRank,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.MoveTemplateResponse{}, nil
}
