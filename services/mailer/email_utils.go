package mailer

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorsmailer "github.com/fivenet-app/fivenet/v2026/services/mailer/errors"
)

const defaultDomain = "fivenet.ls"

func (s *Server) validateEmail(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	input string,
	forJob bool,
) error {
	emails, domains, err := s.generateEmailProposals(ctx, userInfo, forJob)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}

	email, domain, found := strings.Cut(input, "@")
	if !found {
		return errorsmailer.ErrAddresseInvalid
	}

	if len(emails) > 0 && !slices.Contains(emails, email) {
		return errorsmailer.ErrAddresseInvalid
	}

	if !slices.Contains(domains, domain) {
		return errorsmailer.ErrAddresseInvalid
	}

	return nil
}

func (s *Server) generateEmailProposals(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	forJob bool,
) ([]string, []string, error) {
	emails := []string{}
	domains := []string{}

	if forJob {
		// Job's email
		job := s.enricher.GetJobByName(userInfo.GetJob())
		if job == nil {
			return nil, nil, errorsmailer.ErrFailedQuery
		}

		domains = append(domains, fmt.Sprintf("%s.%s", utils.Slug(job.GetName()), defaultDomain))
		domains = append(domains, fmt.Sprintf("%s.%s", utils.Slug(job.GetLabel()), defaultDomain))
		if strings.Contains(job.GetLabel(), " ") {
			labelSplit := strings.Split(job.GetLabel(), " ")
			if len(labelSplit) < 3 {
				for _, split := range labelSplit {
					domains = append(
						domains,
						fmt.Sprintf("%s.%s", utils.Slug(split), defaultDomain),
					)
				}
			} else {
				for idx, split := range labelSplit {
					if idx > 0 && len(labelSplit)-1 >= idx+1 {
						domains = append(
							domains,
							fmt.Sprintf(
								"%s.%s",
								utils.Slug(split+"."+labelSplit[idx+1]),
								defaultDomain,
							),
						)
					}
				}
			}
		}
	} else {
		// User's private email
		user, err := s.store.GetUserShort(ctx, s.db, userInfo.GetUserId())
		if err != nil {
			return nil, nil, errswrap.NewError(err, errorsmailer.ErrFailedQuery)
		}

		// Cleanup name
		user.Firstname = strings.TrimSpace(user.GetFirstname())
		user.Lastname = strings.TrimSpace(user.GetLastname())

		domains = append(domains, defaultDomain)
		emails = append(emails, getBasicNameEmails(user.GetFirstname(), user.GetLastname())...)

		// Generate version with "prefixes" (e.g., `Dr.`) removed
		firstname := utils.RemoveTitlePrefixes(user.GetFirstname())
		if firstname != user.GetFirstname() {
			emails = append(emails, getBasicNameEmails(firstname, user.GetLastname())...)
		}

		// Generate names with birth year added
		dateOfBirth, err := time.Parse("02.01.2006", user.GetDateofbirth())
		if err == nil {
			for _, email := range emails {
				emails = append(emails, fmt.Sprintf("%s%d", email, dateOfBirth.Year()))
			}
		}
	}

	slices.Sort(emails)
	utils.RemoveSliceDuplicates(emails)

	slices.Sort(domains)
	utils.RemoveSliceDuplicates(domains)

	return emails, domains, nil
}

func getBasicNameEmails(firstname string, lastname string) []string {
	return []string{ // User fullname: Erika Mustermann
		utils.Slug(fmt.Sprintf("%s.%s", firstname, lastname)), // erika.mustermann
		utils.Slug(firstname), // erika
		utils.Slug(lastname),  // mustermann
		utils.Slug(
			fmt.Sprintf("%s%s", utils.StringFirstN(firstname, 1), lastname),
		), // emustermann
		utils.Slug(
			fmt.Sprintf("%s%s", firstname, utils.StringFirstN(lastname, 1)),
		), // erikam
		utils.Slug(
			fmt.Sprintf("%s.%s", utils.StringFirstN(firstname, 1), utils.StringFirstN(lastname, 3)),
		), // eri.mus
	}
}
