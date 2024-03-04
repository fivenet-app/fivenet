fx_version 'cerulean'
game 'gta5'

author 'Galexrt'

lua54 'yes'

ui_page "html/index.html"

files {
	'html/index.html',
	'html/index.js',
	'html/style.css',
	'html/images/*.png',
}

server_scripts {
	'@es_extended/imports.lua',
	'@oxmysql/lib/MySQL.lua',
	'server/*.lua',
    'server/events/*.lua'
}

shared_scripts {
	'@es_extended/imports.lua',
	'config.lua'
}

client_scripts {
	'client/*.lua'
}

convar_category 'FiveNet' {
	"Configuration Options",
	{
		{ "Tracking - Clear Table on resource start", "$fnet_clear_on_start", "CV_BOOL", "false" },
	}
}
