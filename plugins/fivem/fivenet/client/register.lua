RegisterNUICallback('exit', function(data, cb)
	SetNuiFocus(false, false)

	cb(true)
end)

RegisterNUICallback('resetPassword', function(data, cb)
	ESX.TriggerServerCallback('fivenet:resetPassword', function(callback)
		if callback then
			SendNUIMessage({
				type = 'token',
				data = callback,
				webUrl = Config.WebURL,
			})

			TriggerEvent('notifications', 'Nutze den Token ~g~'..callback..'~s~ um dein Passwort im FiveNet zurückzusetzen.', 'FIVENET', 'success')
		end
	end)

	cb(true)
end)

RegisterNetEvent('fivenet:registration', function(data)
	if not data then return end

	SetNuiFocus(true, true)
	if type(data) == 'number' then
		SendNUIMessage({
			type = 'token',
			data = data,
			webUrl = Config.WebURL,
		})
	else
		SendNUIMessage({
			type = 'username',
			data = data,
			webUrl = Config.WebURL,
		})
	end
end)

CreateThread(function()
	TriggerEvent('chat:addSuggestion', '/fivenet', 'FiveNet Account Management öffnen')
end)
