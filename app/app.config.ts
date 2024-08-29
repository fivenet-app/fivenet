import type { DiscordConfig, FeatureGates, LoginConfig, Website } from '~/shims';

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
    website: {} as Website,
    featureGates: {} as FeatureGates,
    game: {
        unemployedJobName: 'unemployed',
    },
    statsPage: false,

    // File upload related config
    fileUpload: {
        fileSizes: {
            rector: 5 * 1024 * 1024,
            images: 2 * 1024 * 1024,
        },
        types: {
            images: ['image/jpeg', 'image/jpg', 'image/png'],
        },
    },
    timeouts: {
        grpc: {
            unary: 9000,
        },
        notification: 3500,
    },

    // Nuxt UI app config
    ui: {
        primary: 'blue',
        gray: 'neutral',

        button: {
            default: {
                loadingIcon: 'i-mdi-loading',
            },
        },
        input: {
            default: {
                loadingIcon: 'i-mdi-loading',
            },
        },
        select: {
            default: {
                loadingIcon: 'i-mdi-loading',
                trailingIcon: 'i-octicon-chevron-down-24',
            },
        },
        selectMenu: {
            default: {
                selectedIcon: 'i-mdi-check',
            },
        },
        commandPalette: {
            default: {
                icon: 'i-octicon-search-24',
                loadingIcon: 'i-mdi-loading',
                selectedIcon: 'i-octicon-check-24',
                emptyState: {
                    icon: 'i-octicon-search-24',
                },
            },
        },
        table: {
            default: {
                sortAscIcon: 'i-octicon-sort-asc-24',
                sortDescIcon: 'i-octicon-sort-desc-24',
                sortButton: {
                    icon: 'i-octicon-arrow-switch-24',
                },
                loadingState: {
                    icon: 'i-mdi-loading',
                },
                emptyState: {
                    icon: 'i-octicon-database-24',
                },
            },
            tr: {
                base: 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
            },
            td: {
                padding: 'px-1.5 py-1.5',
            },
        },
        pagination: {
            default: {
                firstButton: {
                    icon: 'i-octicon-chevron-left-24',
                },
                prevButton: {
                    icon: 'i-octicon-arrow-left-24',
                },
                nextButton: {
                    icon: 'i-octicon-arrow-right-24',
                },
                lastButton: {
                    icon: 'i-octicon-chevron-right-24',
                },
            },
        },
        // Nuxt UI Pro Icons
        icons: {
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
        accordion: {
            default: {
                openIcon: 'i-octicon-chevron-down-24',
            },
        },
        breadcrumb: {
            default: {
                divider: 'i-mdi-chevron-right',
            },
        },
        card: {
            footer: {
                padding: 'px-2 py-3 sm:px-4',
            },
        },
        // Nuxt UI Pro
        dashboard: {
            page: {
                wrapper: 'flex w-full min-h-screen min-w-screen max-w-full overflow-y-auto',
            },
        },
    },

    popover: {
        waitTime: 850,
    },

    livemap: {
        userMarkers: {
            activeCharColor: '#fcab10',
            fallbackColor: '#8d81f2',
        },
        markerMarkers: {
            fallbackColor: '#ffffff',
        },
    },
});
