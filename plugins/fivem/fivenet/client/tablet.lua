local usingTablet = false -- false = closed, true = open
local blockInputs = false

-- Objects
local tablet

function IsInTablet()
	return usingTablet and true or false
end

AddEventHandler('fivenet:viewTablet', function(state)
	if not state then
		blockInputs = false
		CloseTablet()
	elseif state then
		blockInputs = false
		OpenTablet()
	end
end)

local dict, anim = 'amb@code_human_in_bus_passenger_idles@female@tablet@idle_a', 'idle_a'

function OpenTablet()
	if IsInTablet() then return end

	local ped = ESX.PlayerData.ped

	-- Don't use tablet animation when in vehicle
	if not IsPedInAnyVehicle(ped) then
		local hash = `prop_cs_tablet`

		local object = CreateObject(hash, 0, 0, 0, true, true, false)
		AttachEntityToEntity(object, ped, GetPedBoneIndex(ped, 28422), -0.05, 0.0, 0.0, 0.0, 0.0, 0.0, true, true, false, true, 1, true)

		tablet = object
    else
        loadAnimDict(dict)
        TaskPlayAnim(ped, anim, 8.0, 1, 50, 0, 0, 0, 0, 0)
	end

	usingTablet = true

	CreateThread(function()
		while usingTablet do
			BlockWeaponWheelThisFrame()

			DisableControlAction(27, 75, true) -- Exit vehicle when driving
            DisableControlAction(0, 0, true)  -- Next Camera
            DisableControlAction(0, 1, true)  -- Look Left/Right
            DisableControlAction(0, 2, true)  -- Look up/Down
            DisableControlAction(0, 36, true) -- Input Duck/Sneak
            DisableControlAction(0, 75, true) -- Exit Vehicle
            DisableControlAction(0, 81, true) -- Next Radio (Vehicle)
            DisableControlAction(0, 82, true) -- Previous Radio (Vehicle)

			if blockInputs then
				DisableAllControlActions(0)
			end

			Wait(0)
		end

		HudWeaponWheelIgnoreControlInput(false)
	end)

	-- Handles pause menu state for tablet
	CreateThread(function()
		while usingTablet do
			Wait(500)
			local isPauseOpen = IsPauseMenuActive() ~= false
			-- Handle if the phone is already visible and escape menu is opened
			if isPauseOpen and IsInTablet() then
				CloseTablet()
                break
			end
		end
	end)

	SendNUIMessage({type = 'openTablet', state = usingTablet, webUrl = Config.WebURL})

	SetNuiFocus(true, true)
	SetNuiFocusKeepInput(true)
end

function CloseTablet()
	if not IsInTablet() then return end

	SetNuiFocus(false, false)
	SetNuiFocusKeepInput(false)

	SendNUIMessage({type = 'closeTablet', state = not usingTablet})

	-- Stop animation
	StopAnimTask(GetPlayerPed(-1), dict, anim, 1.0)

	-- Unblock with delay so escape key isn't handled by the game
	CreateThread(function()
		Wait(100)
		usingTablet = false
	end)

	if DoesEntityExist(tablet) then
		DeleteEntity(tablet)
	end
end

RegisterNUICallback('openTablet', function(data, cb)
	TriggerEvent('fivenet:viewTablet', true)

	cb(true)
end)

RegisterNUICallback('closeTablet', function(data, cb)
	TriggerEvent('fivenet:viewTablet', false)

	cb(true)
end)

RegisterNUICallback('focusTablet', function(data, cb)
	blockInputs = data.state or false

	cb(true)
end)

RegisterCommand('tablet', function()
	TriggerEvent('fivenet:viewTablet', not usingTablet)
end)

RegisterCommand('tabletfix', function()
	TriggerEvent('fivenet:viewTablet', false)
	blockInputs = false

	SendNUIMessage({type = 'fixTablet', webUrl = Config.WebURL})
end)

CreateThread(function()
	TriggerEvent('chat:addSuggestion', '/tablet', 'FiveNet Tablet öffnen')
	TriggerEvent('chat:addSuggestion', '/tabletfix', 'Probleme mit FiveNet Tablet lösen')
end)

RegisterKeyMapping('tablet', 'Tablet öffnen', 'keyboard', 'F5')

-- Hide tablet on resource stop and player relog
AddEventHandler('onResourceStop', function(resourceName)
	if resourceName == GetCurrentResourceName() then
		if usingTablet then
			TriggerEvent('fivenet:viewTablet', false)
		end
	end
end)

RegisterNetEvent('esx:onPlayerLogout', function()
	if IsInTablet() then
		CloseTablet()
	end
end)

-- NUI Callback Handlers for FiveNet actions
RegisterNUICallback('setWaypoint', function(data, cb)
	SetNewWaypoint(data.x, data.y)

	cb(true)
end)

RegisterNUICallback('phoneCallNumber', function(data, cb)
	-- Check if the user has a phone item
	if not ESX.getInventoryItem('phone') and not ESX.getInventoryItem('phone_jailbreak') then
		cb(true)
		return
	end

	if IsInTablet() then
		CloseTablet()
	end

    -- Your phone plugin call a number code here: data.phoneNumber

	cb(true)
end)

RegisterNUICallback('copyToClipboard', function(data, cb)
	SendNUIMessage({type = 'copyToClipboard', text = data.text})

	cb(true)
end)

RegisterNUICallback('setRadioFrequency', function(data, cb)
	local frequency = tonumber(data.frequency)
	if frequency then
        -- This is for pma-voice
		local currentChannel = exports['pma-voice']:getRadioChannel()

		if currentChannel ~= frequency then
			TriggerEvent('tgiann-radio:t', frequency)
		end
	end

	cb(true)
end)

RegisterNUICallback('setWaypointPLZ', function(data, cb)
	ExecuteCommand('plz ' .. data.plz)

	cb(true)
end)
