import type { DiscordConfig, FeatureGates, Links, LoginConfig } from '~/shims';

export default defineAppConfig({
    version: '',
    login: {
        signupEnabled: true,
        lastCharLock: false,
        providers: [],
    } as LoginConfig,
    discord: {
        botInviteURL: '',
    } as DiscordConfig,
    links: {} as Links,
    featureGates: {} as FeatureGates,
    game: {
        unemployedJobName: 'unemployed',
    },

    // Nuxt UI app config
    ui: {
        primary: 'blue',
        gray: 'neutral',
        tooltip: {
            default: {
                openDelay: 500,
            },
        },
        button: {
            default: {
                loadingIcon: 'i-mdi-loading',
            },
        },
        table: {
            tr: {
                base: 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
            },
            td: {
                padding: 'px-1.5 py-1.5',
            },
        },
        selectMenu: {
            default: {
                selectedIcon: 'i-mdi-check',
            },
        },
        card: {
            footer: {
                padding: 'px-2 py-3 sm:px-4',
            },
        },
        icons: {
            dynamic: false,
            // Nuxt UI Pro Icons
            dark: 'i-mdi-moon-and-stars',
            light: 'i-mdi-weather-sunny',
            system: 'i-mdi-computer',
            search: 'i-mdi-search',
            external: 'i-mdi-external-link',
            chevron: 'i-mdi-chevron-down',
            hash: 'i-mdi-hashtag',
            menu: 'i-mdi-menu',
            close: 'i-mdi-window-close',
            check: 'i-mdi-check-circle',
        },
        // Nuxt UI Pro
        dashboard: {
            page: {
                wrapper: 'flex w-full min-h-screen min-w-screen max-w-full overflow-y-auto',
            },
        },
    },
    nuxtIcon: {
        iconifyApiOptions: {
            url: '/api/icons',
            publicApiFallback: false,
        },
    },

    filestore: {
        fileSizes: {
            rector: 5 * 1024 * 1024,
            images: 2 * 1024 * 1024,
        },
        types: {
            images: ['image/jpeg', 'image/jpg', 'image/png'],
        },
    },
    popover: {
        waitTime: 1000,
    },
});
