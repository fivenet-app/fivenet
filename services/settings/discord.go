package settings

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"

	discordapi "github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	pbdiscord "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/discord"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
)

func (s *Server) ListDiscordChannels(ctx context.Context, req *pbsettings.ListDiscordChannelsRequest) (*pbsettings.ListDiscordChannelsResponse, error) {
	if s.dc == nil {
		return nil, errorssettings.ErrDiscordNotEnabled
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jp, err := s.getJobProps(ctx, userInfo.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp := &pbsettings.ListDiscordChannelsResponse{
		Channels: []*pbdiscord.Channel{},
	}
	// No Guild Id set yet, return empty response
	// This is the case when the job is not linked to a Discord guild yet
	if jp.DiscordGuildId == nil || *jp.DiscordGuildId == "" {
		return resp, nil
	}

	jobGuildID, err := strconv.ParseUint(*jp.DiscordGuildId, 10, 64)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	guildId := discord.GuildID(jobGuildID)

	channels, err := s.dc.WithContext(ctx).Channels(guildId)
	if err != nil {
		if restErr, ok := err.(*httputil.HTTPError); ok {
			if restErr.Status == http.StatusNotFound {
				return resp, nil // Guild not found, return empty response
			}
		}
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	for _, channel := range channels {
		if channel.Type != discord.GuildText {
			continue // Only include text and news channels
		}

		resp.Channels = append(resp.Channels, &pbdiscord.Channel{
			Id:       channel.ID.String(),
			GuildId:  channel.GuildID.String(),
			Name:     channel.Name,
			Type:     uint32(channel.Type),
			Position: int64(channel.Position),
		})
	}

	slices.SortStableFunc(resp.Channels, func(a, b *pbdiscord.Channel) int {
		return int(a.Position - b.Position)
	})

	return resp, nil
}

func (s *Server) ListUserGuilds(ctx context.Context, req *pbsettings.ListUserGuildsRequest) (*pbsettings.ListUserGuildsResponse, error) {
	if s.dc == nil {
		return nil, errorssettings.ErrDiscordNotEnabled
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	acc, err := accounts.RetrieveOAuth2Account(ctx, s.db, s.crypt, userInfo.AccountId, "discord")
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	accessToken, err := accounts.GetAccessToken(ctx, s.db, s.crypt, acc, s.dcOAuth2Provider.ClientID, s.dcOAuth2Provider.ClientSecret, s.dcOAuth2Provider.RedirectURL)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	if accessToken == "" {
		return nil, errorssettings.ErrDiscordConnectRequired
	}

	guilds, err := getUserGuilds(ctx, accessToken)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp := &pbsettings.ListUserGuildsResponse{
		Guilds: []*pbdiscord.Guild{},
	}
	for _, guild := range guilds {
		resp.Guilds = append(resp.Guilds, &pbdiscord.Guild{
			Id:        guild.ID.String(),
			Name:      guild.Name,
			Icon:      guild.IconURL(),
			CreatedAt: timestamp.New(guild.CreatedAt()),
		})
	}

	slices.SortFunc(resp.Guilds, func(a, b *pbdiscord.Guild) int {
		return strings.Compare(a.Name, b.Name)
	})

	return resp, nil
}

func getUserGuilds(ctx context.Context, accessToken string) ([]discord.Guild, error) {
	dc := discordapi.NewClient("Bearer " + accessToken).WithContext(ctx)

	guilds, err := dc.Guilds(200)
	if err != nil {
		return nil, err
	}

	return guilds, nil
}
