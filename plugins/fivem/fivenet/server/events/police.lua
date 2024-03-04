-- Jail
AddEventHandler('esx_prison:jailPlayer', function(pPlayer, xPlayer, time --[[number]])
	addUserActivity(pPlayer.identifier, xPlayer.identifier, 0, 'Plugin.Jail', '', time, '')
end)

AddEventHandler('esx_prison:unjailedByPlayer', function(xPlayer, pPlayer, _, type --[[ 'police'/ 'admin']])
	addUserActivity(pPlayer.identifier, xPlayer.identifier, 0, 'Plugin.Jail', '', '0', type)
end)

AddEventHandler('esx_prison:escapePoliceNotify', function(xPlayer)
	addUserActivity(xPlayer.identifier, xPlayer.identifier, 0, 'Plugin.Jail', '0', '', '')
	-- Set user wanted and add user activity item
	setUserWantedState(xPlayer.identifier, true)
	addUserActivity(xPlayer.identifier, xPlayer.identifier, 0, 'UserProps.Wanted', 'false', 'true', 'Gefängnisausbruch')
end)

-- Panicbutton
AddEventHandler('esx_policeJob:panicButton', function(source, x --[[number]], y --[[number]], _, name)
	local xPlayer = ESX.GetPlayerFromId(source)
	if not xPlayer then return end

	-- Send panic button dispatches to source user's job only for now
	createDispatch(xPlayer.job.name, 'Panikknopf ausgelöst', name, x, y, false, xPlayer.identifier)
end)
