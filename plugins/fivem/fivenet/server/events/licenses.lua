-- These are non standard events, they need to be added by yourself to take advantage of
-- the job user activity tracking

local function lookupLicenseLabel(type)
    local response = MySQL.query.await('SELECT `label` FROM `licenses` WHERE `type` = ? LIMIT 1', { type })
    if response then
        return response[1].label or ''
    end

    return ''
end

AddEventHandler('esx_license:addLicense', function(sourceXPlayer, targetXPlayer, type)
	addUserActivity(sourceXPlayer.identifier, targetXPlayer.identifier, 0, 'Plugin.Licenses', '', type, lookupLicenseLabel(type))
end)

AddEventHandler('esx_license:removeLicense', function(sourceXPlayer, targetXPlayer, type)
    -- You should look up and pass the license name as the last argument
	addUserActivity(sourceXPlayer.identifier, targetXPlayer.identifier, 0, 'Plugin.Licenses', type, '', lookupLicenseLabel(type))
end)
