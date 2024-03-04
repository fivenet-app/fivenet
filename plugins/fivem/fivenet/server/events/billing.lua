-- Fines
local function checkIfBillingEnabledJob(target --[[string]])
	local job = string.gsub(target, 'society_', '')
	return Config.Events.BillingJobs[job]
end

AddEventHandler('esx_billing:sentBill', function(xPlayer, xTarget, type --[['fine'/ 'bill']], label, amount)
	if type ~= 'fine' then return end
	if not checkIfBillingEnabledJob(xPlayer.job.name) then return end

	updateOpenFines(xTarget.identifier, amount)
	addUserActivity(xPlayer.identifier, xTarget.identifier, 0, 'Plugin.Billing.Fines', '', amount, label)
end)

AddEventHandler('esx_billing:removedBill', function(source, type, result)
	if type ~= 'fine' then return end
	if result.target_type ~= 'society' then return end
	if not checkIfBillingEnabledJob(result.target) then return end

	local xPlayer = ESX.GetPlayerFromId(source)
	if not xPlayer then return end

	updateOpenFines(result.identifier, -result.amount)
	addUserActivity(xPlayer.identifier, result.identifier, 0, 'Plugin.Billing.Fines', result.amount, result.amount, result.label)
end)

AddEventHandler('esx_billing:paidBill', function(source, result)
	if result.target_type ~= 'society' then return end
	if not checkIfBillingEnabledJob(result.target) then return end

	local xPlayer = ESX.GetPlayerFromId(source)
	if not xPlayer then return end

	updateOpenFines(xPlayer.identifier, -result.amount)
	addUserActivity(xPlayer.identifier, xPlayer.identifier, 0, 'Plugin.Billing.Fines', result.amount, 0, result.label)
end)
