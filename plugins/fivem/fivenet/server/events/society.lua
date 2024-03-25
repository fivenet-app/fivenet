-- These are non standard events, they need to be added by yourself to take advantage of
-- the job user activity tracking

AddEventHandler('esx_society:hired', function(xPlayer, xTarget)
	local sIdentifier = ESX.GetPlayerInfo(xPlayer, "identifier")
	if not sIdentifier then return end
	local tIdentifier = ESX.GetPlayerInfo(xTarget, "identifier")
	if not tIdentifier then return end

	local job = ESX.GetPlayerInfo(xTarget, "job")
	if not job then return end

	addJobsUserActivity(job.name, sIdentifier, tIdentifier, 1, nil, "{}")
end)

AddEventHandler('esx_society:gradeChanged', function(xPlayer, xTarget, promoted)
	local sIdentifier = ESX.GetPlayerInfo(xPlayer, "identifier")
	if not sIdentifier then return end
	local tIdentifier = ESX.GetPlayerInfo(xTarget, "identifier")
	if not tIdentifier then return end

	local job = ESX.GetPlayerInfo(xTarget, "job")
	if not job then return end

	local data = { gradeChange = { grade = job.grade, gradeLabel = job.grade_label }}
	addJobsUserActivity(job.name, sIdentifier, tIdentifier, promoted and 3 or 4, nil, json.encode(data))
end)

AddEventHandler('esx_society:fired', function(xPlayer, xTarget)
	local sIdentifier = ESX.GetPlayerInfo(xPlayer, "identifier")
	if not sIdentifier then return end
	local tIdentifier = ESX.GetPlayerInfo(xTarget, "identifier")
	if not tIdentifier then return end

	local job = ESX.GetPlayerInfo(xPlayer, "job")
	if not job then return end

	addJobsUserActivity(job.name, sIdentifier, tIdentifier, 2, nil, "{}")
end)
