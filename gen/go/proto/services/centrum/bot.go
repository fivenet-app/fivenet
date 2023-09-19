package centrum

type BotManager struct {
	bots map[string]*Bot
}

func NewBotManager() *BotManager {
	return &BotManager{
		bots: map[string]*Bot{},
	}
}

type Bot struct {
	job string
}

// TODO write automatic bot code
