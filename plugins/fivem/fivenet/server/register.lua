function getLicenseFromIdentifier(identifier)
	local start = string.find(identifier, ':', 1, true)
	if not start then return identifier end

	return string.sub(identifier, start + 1, -1)
end

local function isRegistered(xPlayer)
	local query = MySQL.single.await('SELECT username FROM fivenet_accounts WHERE license = ?', {getLicenseFromIdentifier(xPlayer.identifier)})
	if query then
		return query.username
	end

	return false
end

local TokenLength = 6

local function createToken(xPlayer)
	local token = ''
	for i = 1, TokenLength do
		token = token..math.random(1, 9)
	end

	MySQL.update('INSERT INTO fivenet_accounts (enabled, license, reg_token) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE reg_token = ?', {1, getLicenseFromIdentifier(xPlayer.identifier), token, token})
	return token
end

local function forgotPassword(xPlayer)
	local token = ''
	for i = 1, TokenLength do
		token = token..math.random(1, 9)
	end

	MySQL.update('UPDATE fivenet_accounts SET password = NULL, reg_token = ? WHERE license = ?', {token, getLicenseFromIdentifier(xPlayer.identifier)})
	return token
end

ESX.RegisterServerCallback('fivenet:resetPassword', function(source, cb)
	local xPlayer = ESX.GetPlayerFromId(source)
	if not xPlayer then return end

	if isRegistered(xPlayer) then
		local token = forgotPassword(xPlayer)
		cb(token)
	else
		cb(false)
	end
end)

RegisterCommand('fivenet', function(source, args)
	local xPlayer = ESX.GetPlayerFromId(source)
	if not xPlayer then return end

	local registered = isRegistered(xPlayer)
	if not registered then
		registered = tonumber(createToken(xPlayer))
	end

	TriggerClientEvent('fivenet:registration', source, registered)
end)
