-- FiveNet Account - Social Login
function addOrSetDiscordIdentifier(license --[[string]], externalId --[[string]], username --[[string]])
	-- Check if user has a FiveNet account
	MySQL.query('SELECT `id` FROM `fivenet_accounts` WHERE `license` = ? LIMIT 1', { getLicenseFromIdentifier(license) },
		function(result)
			result = result and result[1] or nil
			if not result or not result.id  then return end

			-- If the user has an account add or update a Discord Oauth2 "connection"
			MySQL.update([[
				INSERT INTO `fivenet_oauth2_accounts`
				(`account_id`, `provider`, `external_id`, `username`, `avatar`)
				VALUES(?, ?, ? , ?, 'https://cdn.discordapp.com/embed/avatars/0.png')
				ON DUPLICATE KEY UPDATE `external_id` = VALUES(`external_id`)
				]],
				{ result.id, Config.DiscordOAuth2Provider, getLicenseFromIdentifier(externalId), username })
	end)
end

-- User Props
function addUserActivity(sIdentifier --[[string]], tIdentifier --[[string]], type --[[number]], key --[[string]], oldVal --[[string]], newVal --[[string]], reason --[[string]])
	MySQL.update([[
		INSERT INTO `fivenet_user_activity`
		(`source_user_id`, `target_user_id`, `type`, `key`, `old_value`, `new_value`, `reason`)
		VALUES ((SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), (SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), ?, ?, ?, ?, ?)
		]],
		{ sIdentifier, tIdentifier, type, key, oldVal, newVal, reason })
end
exports('addUserActivity', addUserActivity)

function updateOpenFines(tIdentifier --[[string]], fine --[[number]])
	MySQL.update([[
		INSERT INTO `fivenet_user_props`
		(`user_id`, `open_fines`)
		VALUES ((SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), ?)
		ON DUPLICATE KEY UPDATE `open_fines` = CASE WHEN `open_fines` + VALUES(`open_fines`) < 0 THEN 0 ELSE `open_fines` + VALUES(`open_fines`) END;
		]],
		{ tIdentifier, fine })
end
exports('updateOpenFines', updateOpenFines)

function setUserWantedState(tIdentifier --[[string]], wanted --[[bool]])
	MySQL.update([[
		INSERT INTO `fivenet_user_props`
		(`user_id`, `wanted`)
		VALUES ((SELECT id FROM `users` WHERE `identifier` = ? LIMIT 1), ?)
		ON DUPLICATE KEY UPDATE `wanted` = VALUES(`wanted`)
		]],
		{ tIdentifier, wanted })
end
exports('setUserWantedState', setUserWantedState)

function setUserBloodType(tIdentifier --[[string]], bloodType --[[string]])
	MySQL.update([[
		INSERT INTO `fivenet_user_props`
		(`user_id`, `blood_type`)
		VALUES ((SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), ?)
		ON DUPLICATE KEY UPDATE `blood_type` = COALESCE(`blood_type`, VALUES(`blood_type`))
		]],
		{ tIdentifier, bloodType })
end
exports('setUserBloodType', setUserBloodType)

-- Jobs User Activity
-- activityType: 1 = HIRED, 2 = FIRED, 3 = PROMOTED, 4 = DEMOTED
function addJobsUserActivity(job --[[string]], sIdentifier --[[string]], tIdentifier --[[string]], activityType --[[number]], reason --[[string]], data --[[string]])
	MySQL.update([[
		INSERT INTO `fivenet_jobs_user_activity`
		(`job`, `source_user_id`, `target_user_id`, `activity_type`, `reason`, `data`)
		VALUES (?, (SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), (SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1), ?, ?, ?)
		]],
		{ job, sIdentifier, tIdentifier, activityType, reason, data })
end
exports('addJobsUserActivity', addJobsUserActivity)

-- Dispatches
function createDispatch(job --[[string]], message --[[string]], description --[[string]], x --[[number]], y --[[number]], anon --[[bool]], identifier --[[string]])
	MySQL.update([[
		INSERT INTO `fivenet_centrum_dispatches`
		(`job`, `message`, `description`, `x`, `y`, `anon`, `creator_id`)
		VALUES (?, ?, ?, ?, ?, ?, (SELECT `id` FROM `users` WHERE `identifier` = ? LIMIT 1))
		]],
		{ job, message, description, x, y, anon, identifier })
end
exports('createDispatch', createDispatch)

function createCivilProtectionJobDispatch(message --[[string]], description --[[string]], x --[[number]], y --[[number]], anon --[[bool]], identifier --[[string]])
	for job, _ in pairs(Config.Dispatches.CivilProtectionJobs) do
		createDispatch(job, message, description, x, y, anon, identifier)
	end
end
exports('createCivilProtectionJobDispatch', createCivilProtectionJobDispatch)

-- Written by mcnuggets
function loadAnimDict(dict)
    if not HasAnimDictLoaded(dict) then
		RequestAnimDict(dict)

		while not HasAnimDictLoaded(dict) do
			Wait(10)
		end
	end
end

-- Written by mcnuggets
function IsNearVector(source, targetVector, range)
	range = range or 3.0

	local sourcePed = GetPlayerPed(source)
	if sourcePed == 0 then return false end
	local sourceCoords = GetEntityCoords(sourcePed)

	if #(sourceCoords - targetVector) > range then
		return false
	end

	return true
end
