package types

type State struct {
	GuildID string

	Roles Roles
	Users Users

	UserProcessors []UserProcessorHandler
}

func (s *State) Merge(state *State) {
	if state == nil {
		return
	}

	s.Roles = append(s.Roles, state.Roles...)

	for _, user := range state.Users {
		s.Users.Add(user)
	}

	s.UserProcessors = append(s.UserProcessors, state.UserProcessors...)
}
