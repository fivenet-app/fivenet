package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BCryptCost parameter for bcrypt hashing.
const BCryptCost = 14

func (s *Server) getCookieBase(name string, value string) http.Cookie {
	//nolint:gosec // `Same-Site: None` is required because otherwise users can't login in the in-game tablet (iframe).
	return http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().UTC().Add(auth.TokenExpireTime),
		MaxAge:   int(auth.TokenExpireTime.Seconds()),
		Domain:   s.domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

func (s *Server) setCookies(
	ctx context.Context,
	accClaims *authclaims.AccountInfoClaims,
) error {
	//nolint:gosec // getCookieBase returns a secure pre-configured cookie "base"
	authedCookie := s.getCookieBase(auth.AuthedCookieName, "true")
	authedCookie.HttpOnly = false
	header := metadata.Pairs("set-cookie", authedCookie.String())

	accToken, err := s.tm.FromAccClaims(accClaims)
	if err != nil {
		return err
	}

	accCookie := s.getCookieBase(auth.AccCookieName, accToken)
	header.Append("set-cookie", accCookie.String())

	// Send the cookies back to the client
	return grpc.SendHeader(ctx, header)
}

func (s *Server) destroyCookies(ctx context.Context) error {
	//nolint:gosec // getCookieBase returns a secure pre-configured cookie "base"
	authedCookie := s.getCookieBase(auth.AuthedCookieName, "false")
	authedCookie.HttpOnly = false

	//nolint:gosec // getCookieBase returns a secure pre-configured cookie "base"
	accCookie := s.getCookieBase(auth.AccCookieName, "")
	accCookie.Expires = time.Unix(0, 0).UTC()
	accCookie.MaxAge = -1

	// Send the cookies back to the client
	header := metadata.Pairs(
		"set-cookie", accCookie.String(),
		"set-cookie", authedCookie.String(),
	)
	return grpc.SendHeader(ctx, header)
}

func (s *Server) getCharacter(
	ctx context.Context,
	charID int32,
) (*users.User, *jobsprops.JobProps, error) {
	char, jProps, err := s.store.GetCharacter(ctx, charID)
	if err != nil {
		return nil, nil, err
	}

	if jProps != nil {
		s.enricher.EnrichJobName(jProps)
	}
	s.enricher.EnrichJobInfo(char)
	s.enrichCharacterJobs(char)

	return char, jProps, nil
}

// enrichCharacterJobs supplies display labels for every server-synchronized
// work-context job. These are returned to the client as part of the character
// so the context selector never needs to resolve internal job names itself.
func (s *Server) enrichCharacterJobs(char *users.User) {
	for _, userJob := range char.GetJobs() {
		job := s.enricher.GetJobByName(userJob.GetJob())
		if job == nil {
			userJob.SetJobLabel(mstlystcdata.NotAvailablePlaceholder)
			userJob.SetGradeLabel(mstlystcdata.NotAvailablePlaceholder)
			continue
		}

		userJob.SetJobLabel(job.GetLabel())
		gradeLabel := strconv.FormatInt(int64(userJob.GetGrade()), 10)
		for _, grade := range job.GetGrades() {
			if grade.GetGrade() == userJob.GetGrade() {
				gradeLabel = grade.GetLabel()
				break
			}
		}
		userJob.SetGradeLabel(gradeLabel)
	}
}

func (s *Server) getJobWithProps(
	ctx context.Context,
	jobName string,
) (*jobs.Job, int32, *jobsprops.JobProps, error) {
	job, jobGrade, jProps, err := s.store.GetJobWithProps(ctx, jobName)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, 0, nil, err
		}

		job = &jobs.Job{Name: jobName, Label: jobName}
		jProps = &jobsprops.JobProps{Job: jobName}
		jProps.JobLabel = &jobName
		jProps.Default(jobName)
	}

	if jProps != nil {
		s.enricher.EnrichJobName(jProps)
	}

	return job, jobGrade, jProps, nil
}

// Helper for password hashing.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	return string(hashed), err
}

// Helper for password validation.
func checkPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// Helper for username normalization.
func normalizeUsername(username string) string {
	return strings.TrimSpace(username)
}
