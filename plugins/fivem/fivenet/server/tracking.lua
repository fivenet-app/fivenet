local playerLocations = {}

local locationUpdateQuery = [[
	INSERT INTO `fivenet_user_locations` (`identifier`, `job`, `x`, `y`, `hidden`) VALUES (@identifier, @job, @x, @y, @hidden)
	ON DUPLICATE KEY UPDATE `job` = VALUES(`job`), `x` = VALUES(`x`), `y` = VALUES(`y`), `hidden` = VALUES(`hidden`);
]]

local function deletePosition(identifier)
	playerLocations[identifier] = nil

	MySQL.update('DELETE FROM `fivenet_user_locations` WHERE `identifier` = ? LIMIT 1', { identifier })
end

local function checkIfPlayerHidden(xPlayer)
	return not xPlayer.job.onDuty or (Config.TrackingItem and not xPlayer.getInventoryItem(Config.TrackingItem))
end

if Config.EnableTracking then
	CreateThread(function()
		while true do
			local queries = {}

			for playerId, xPlayer in pairs(ESX.GetExtendedPlayers()) do
				if Config.TrackingJobs[xPlayer.job.name] then
					local update = true

					if playerLocations[xPlayer.identifier] then
						local curLocation = playerLocations[xPlayer.identifier]
						if IsNearVector(playerId, curLocation, 5.0) then
							update = false
						end
					end

					if update then
						local ped = GetPlayerPed(playerId)
						if ped ~= 0 then
							local coords = GetEntityCoords(ped)
							local hidden = 0

							-- Either players is not on duty and/or doesn't have the tracking item
							if checkIfPlayerHidden(xPlayer) then
								hidden = 1
							end

							playerLocations[xPlayer.identifier] = coords

							table.insert(queries, { locationUpdateQuery,
								{
									["identifier"] = xPlayer.identifier,
									["job"] = xPlayer.job.name,
									["x"] = coords.x,
									["y"] = coords.y,
									["hidden"] = hidden,
								}
							})
						end
					end
				end
			end

			MySQL.transaction(queries)

			Wait(Config.TrackingInterval)
		end
	end)

	AddEventHandler('esx:playerDropped', function(source)
		local xPlayer = ESX.GetPlayerFromId(source)
		if not xPlayer then return end

		deletePosition(xPlayer.identifier)
	end)
end

-- Resource Start
AddEventHandler('onResourceStart', function(resourceName)
	if resourceName == GetCurrentResourceName() and GetConvar('fnet_clear_on_start', 'false') == 'true' then
		-- Truncate user locations table on resource (re-)start
		MySQL.update('DELETE FROM `fivenet_user_locations`')
	end
end)
