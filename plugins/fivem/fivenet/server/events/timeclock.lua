-- Timeclock On/Offduty tracking
local function timeclockTrack(job --[[string]], identifier --[[string]], clockOn --[[bool]])
	if not Config.TimeclockJobs[job] then return end

	if clockOn then
		-- Run select query to see if a timeclock entry needs to be created or updated
		MySQL.query([[
			SELECT u.`id` AS `userId`, fjt.`start_time` AS `startTime`
			FROM `fivenet_jobs_timeclock` fjt INNER JOIN `users` u ON (u.`id` = fjt.`user_id`)
			WHERE u.`identifier` = ? AND fjt.`user_id` = u.`id`
			LIMIT 1
			]],
			{ identifier },
			function(result)
				result = result and result[1] or nil
				if not result or not result.userId then return end

				-- If start time is not null, the entry is (already) active, keep using it
				if result and result.startTime then return end

				MySQL.update([[
					INSERT INTO `fivenet_jobs_timeclock`
					(`job`, `user_id`, `date`)
					VALUES(?, ?, CURRENT_TIMESTAMP)
					ON DUPLICATE KEY UPDATE `start_time` = VALUES(`start_time`)
					]],
					{ job, result.userId })
		end)
	else
		MySQL.update([[
			UPDATE `fivenet_jobs_timeclock` fjt INNER JOIN `users` u ON (u.`id` = fjt.`user_id`)
			SET fjt.`end_time` = CURRENT_TIMESTAMP
			WHERE u.`identifier` = ? AND fjt.`user_id` = u.`id` AND fjt.`start_time` IS NOT NULL AND fjt.`end_time` IS NULL
			]],
			{ identifier })
	end
end

AddEventHandler('esx:setJob', function(playerId)
	local xPlayer = ESX.GetPlayerFromId(playerId)
	if not xPlayer then return end

	-- Check if job is enabled for timeclock tracking
	if not Config.TimeclockJobs[xPlayer.job.name] then return end

	-- If lastJob is nil, user left job's duty
	timeclockTrack(xPlayer.job.name, xPlayer.identifier, xPlayer.job.onDuty)
end)

AddEventHandler('esx:playerDropped', function(playerId)
	local xPlayer = ESX.GetPlayerFromId(playerId)
	if not xPlayer then return end

	-- Check if job is enabled for timeclock tracking
	if not Config.TimeclockJobs[xPlayer.job.name] then return end

	timeclockTrack(xPlayer.job.name, xPlayer.identifier, false)
end)
