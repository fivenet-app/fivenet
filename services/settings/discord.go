package settings

import (
	"context"
	"slices"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	pbdiscord "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/discord"
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
	if jp.DiscordGuildId == nil || *jp.DiscordGuildId == "" {
		return nil, errorssettings.ErrDiscordNotEnabled
	}

	resp := &pbsettings.ListDiscordChannelsResponse{
		Channels: []*pbdiscord.Channel{},
	}

	jobGuildID, err := strconv.ParseUint(*jp.DiscordGuildId, 10, 64)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	guildId := discord.GuildID(jobGuildID)

	channels, err := s.dc.WithContext(ctx).Channels(guildId)
	if err != nil {
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
