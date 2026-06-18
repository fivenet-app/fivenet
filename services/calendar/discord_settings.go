package calendar

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	errorscalendar "github.com/fivenet-app/fivenet/v2026/services/calendar/errors"
	"google.golang.org/protobuf/proto"
)

const maxCalendarDiscordReminderSteps = 2

func (s *Server) prepareCalendarDiscordSettings(
	ctx context.Context,
	cal *calendarresource.Calendar,
) (*calendarresource.CalendarDiscordSettings, *string, error) {
	if cal == nil {
		return nil, nil, nil
	}

	settings := normalizeCalendarDiscordSettings(cal.DiscordSettings)
	if err := s.validateCalendarDiscordSettings(
		ctx,
		cal.GetJob(),
		cal.GetSystemKind(),
		settings,
	); err != nil {
		return nil, nil, err
	}
	if settings == nil {
		return nil, nil, nil
	}

	raw, err := protoutils.MarshalToJSON(settings)
	if err != nil {
		return nil, nil, err
	}

	out := string(raw)
	return settings, &out, nil
}

func normalizeCalendarDiscordSettings(
	settings *calendarresource.CalendarDiscordSettings,
) *calendarresource.CalendarDiscordSettings {
	if settings == nil {
		return nil
	}

	out := proto.Clone(settings).(*calendarresource.CalendarDiscordSettings)
	out.ChannelId = strings.TrimSpace(out.GetChannelId())

	steps := make([]*calendarresource.CalendarDiscordReminderStep, 0, len(out.GetReminderSteps()))
	for i := range out.GetReminderSteps() {
		step := out.GetReminderSteps()[i]
		if step == nil {
			continue
		}

		clone := proto.Clone(step).(*calendarresource.CalendarDiscordReminderStep)
		clone.Message = normalizeOptionalString(clone.Message)
		clone.Embed = normalizeCalendarDiscordReminderEmbed(clone.GetEmbed())
		steps = append(steps, clone)
	}
	out.ReminderSteps = steps

	if !out.GetEnabled() && out.GetChannelId() == "" && len(out.GetReminderSteps()) == 0 {
		return nil
	}

	return out
}

func normalizeCalendarDiscordReminderEmbed(
	embed *calendarresource.CalendarDiscordReminderEmbed,
) *calendarresource.CalendarDiscordReminderEmbed {
	if embed == nil {
		return nil
	}

	out := proto.Clone(embed).(*calendarresource.CalendarDiscordReminderEmbed)
	out.Title = normalizeOptionalString(out.Title)
	out.Description = normalizeOptionalString(out.Description)
	out.Color = normalizeOptionalString(out.Color)

	if !calendarDiscordEmbedHasContent(out) {
		return nil
	}

	return out
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func calendarDiscordEmbedHasContent(embed *calendarresource.CalendarDiscordReminderEmbed) bool {
	if embed == nil {
		return false
	}

	return embed.Title != nil || embed.Description != nil
}

func (s *Server) validateCalendarDiscordSettings(
	ctx context.Context,
	job string,
	systemKind calendarresource.CalendarSystemKind,
	settings *calendarresource.CalendarDiscordSettings,
) error {
	if settings == nil {
		return nil
	}

	if len(settings.GetReminderSteps()) > maxCalendarDiscordReminderSteps {
		return errswrap.NewError(fmt.Errorf(
			"at most %d discord reminder steps are allowed",
			maxCalendarDiscordReminderSteps,
		), errorscalendar.ErrInvalidReminderStep)
	}

	seen := map[int32]struct{}{}
	for idx := range settings.GetReminderSteps() {
		step := settings.GetReminderSteps()[idx]
		if step == nil {
			return errswrap.NewError(
				fmt.Errorf("discord reminder step %d is invalid", idx+1),
				errorscalendar.ErrInvalidReminderStep,
			)
		}

		if _, ok := seen[step.GetAtMinute()]; ok {
			return errswrap.NewError(
				errors.New("discord reminder step times must be unique"),
				errorscalendar.ErrInvalidReminderStep,
			)
		}
		seen[step.GetAtMinute()] = struct{}{}

		if step.GetMessage() == "" && !calendarDiscordEmbedHasContent(step.GetEmbed()) {
			return errswrap.NewError(
				fmt.Errorf("discord reminder step %d needs a message or embed content", idx+1),
				errorscalendar.ErrInvalidReminderStep,
			)
		}
	}

	if strings.TrimSpace(job) == "" ||
		systemKind != calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED {
		if settings.GetEnabled() {
			return errswrap.NewError(
				errors.New("discord reminders can only be enabled on job calendars"),
				errorscalendar.ErrInvalidReminderStep,
			)
		}
		return nil
	}

	if !settings.GetEnabled() {
		return nil
	}

	if s.dc == nil {
		return errswrap.NewError(
			errors.New("discord bot is not enabled"),
			errorscalendar.ErrInvalidReminderStep,
		)
	}
	if strings.TrimSpace(settings.GetChannelId()) == "" {
		return errswrap.NewError(
			errors.New("discord reminder channel is required"),
			errorscalendar.ErrInvalidReminderStep,
		)
	}
	if len(settings.GetReminderSteps()) == 0 {
		return errswrap.NewError(
			errors.New("at least one discord reminder step is required"),
			errorscalendar.ErrInvalidReminderStep,
		)
	}

	return s.validateCalendarReminderChannel(ctx, job, settings.GetChannelId())
}

func (s *Server) validateCalendarReminderChannel(
	ctx context.Context,
	job string,
	channelID string,
) error {
	guildIDRaw, err := s.store.GetCalendarReminderGuildID(ctx, job)
	if err != nil {
		return err
	}

	guildIDNum, err := strconv.ParseUint(guildIDRaw, 10, 64)
	if err != nil {
		return errswrap.NewError(
			errors.New("discord guild id is invalid"),
			errorscalendar.ErrInvalidDiscordChannel,
		)
	}
	channelIDNum, err := strconv.ParseUint(strings.TrimSpace(channelID), 10, 64)
	if err != nil {
		return errswrap.NewError(
			errors.New("discord channel id is invalid"),
			errorscalendar.ErrInvalidDiscordChannel,
		)
	}

	channel, err := s.dc.WithContext(ctx).Channel(discord.ChannelID(channelIDNum))
	if err != nil {
		return errswrap.NewError(
			errors.New("discord channel could not be loaded"),
			errorscalendar.ErrInvalidDiscordChannel,
		)
	}
	if channel.GuildID != discord.GuildID(guildIDNum) {
		return errswrap.NewError(
			errors.New("discord channel does not belong to the job guild"),
			errorscalendar.ErrInvalidDiscordChannel,
		)
	}
	if channel.Type != discord.GuildText {
		return errswrap.NewError(
			errors.New("discord reminder channel must be a text channel"),
			errorscalendar.ErrInvalidDiscordChannel,
		)
	}

	return nil
}
