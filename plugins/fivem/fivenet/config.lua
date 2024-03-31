Config = {}

-- Your FiveNet URL without a trailing slash
Config.WebURL = "https://fivenet.app"

Config.EnableTracking = true -- Enable the tracking of players
Config.TrackingJobs = {  -- Those jobs will be tracked
	["ambulance"] = true,
	["doj"] = true,
	["police"] = true,
}
Config.TrackingItem = "radio" -- Players without this item will be updated as 'hidden', set false otherwise
Config.TrackingInterval = 3000 -- Interval in ms until positions will be updated

Config.TimeclockJobs = { -- These jobs will be timeclocked tracked
	["ambulance"] = true,
	["doj"] = true,
    ["police"] = true,
    -- Can also be other jobs that are onDuty enabled
}

Config.Events = {}
Config.Events.BillingJobs = { -- Jobs bills that will cause an user activity to be created for the billing cycle events
	["doj"] = true,
	["police"] = true,
}

Config.DiscordOAuth2Provider = "discord"

Config.Dispatches = {}
Config.Dispatches.CivilProtectionJobs = {
	["police"] = true,
}

Config.UserProps = {}
Config.UserProps.BloodTypes = {
	"A+", "A-",
    "B+", "B-",
    "AB+", "AB-",
    "O+", "O-",
}
