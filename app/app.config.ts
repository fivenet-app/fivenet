import type { DiscordConfig, FeatureGates, GameConfig, LoginConfig, SystemConfig, WebsiteConfig } from '~/typings';

export default defineAppConfig({
    // Server provided App Config
    version: '',

    defaultLocale: 'en',

    login: {
        signupEnabled: true,
        lastCharLock: false,
        providers: [],
    } as LoginConfig,
    discord: {
        botInviteURL: '',
    } as DiscordConfig,
    website: {} as WebsiteConfig,
    featureGates: {
        imageProxy: false,
    } as FeatureGates,
    game: {
        unemployedJobName: 'unemployed',
        startJobGrade: 0,
    } as GameConfig,
    system: {} as SystemConfig,

    // File upload related config
    fileUpload: {
        fileSizes: {
            rector: 5 * 1024 * 1024,
            images: 2 * 1024 * 1024,
        },
        types: {
            images: ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'],
        },
    },
    // Request timeouts
    timeouts: {
        grpc: {
            unary: 9000,
        },
        notification: 3500,
    },
    maxAccessEntries: 12,

    fallbackColor: 'blue',

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

    // Nuxt UI and UI Pro config
    ui: {
        primary: 'blue',
        gray: 'neutral',

        avatar: {
            default: {
                icon: 'i-mdi-file-image-box',
            },
        },
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
        inputMenu: {
            default: {
                selectedIcon: 'i-mdi-check',
                trailingIcon: 'i-mdi-chevron-down',
            },
        },
        select: {
            default: {
                loadingIcon: 'i-mdi-loading',
                trailingIcon: 'i-mdi-chevron-down',
            },
        },
        selectMenu: {
            default: {
                selectedIcon: 'i-mdi-check',
            },
        },
        notification: {
            default: {
                closeButton: {
                    icon: 'i-mdi-window-close',
                },
            },
        },
        commandPalette: {
            default: {
                icon: 'i-mdi-search',
                loadingIcon: 'i-mdi-loading',
                selectedIcon: 'i-mdi-check',
                emptyState: {
                    icon: 'i-mdi-search',
                },
                closeButton: {
                    icon: 'i-mdi-window-close',
                },
            },
        },
        table: {
            default: {
                sortAscIcon: 'i-mdi-sort-ascending',
                sortDescIcon: 'i-mdi-sort-descending',
                sortButton: {
                    icon: 'i-mdi-sort',
                },
                loadingState: {
                    icon: 'i-mdi-loading',
                },
                emptyState: {
                    icon: 'i-mdi-database',
                },
                expandButton: {
                    icon: 'i-mdi-chevron-down',
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
                    icon: 'i-mdi-chevron-left-first',
                },
                prevButton: {
                    icon: 'i-mdi-chevron-left',
                },
                nextButton: {
                    icon: 'i-mdi-chevron-right',
                },
                lastButton: {
                    icon: 'i-mdi-chevron-right-last',
                },
            },
        },
        accordion: {
            default: {
                openIcon: 'i-mdi-chevron-down',
            },
        },
        breadcrumb: {
            default: {
                divider: 'i-mdi-chevron-right',
            },
        },
        carousel: {
            default: {
                prevButton: {
                    icon: 'i-mdi-chevron-left',
                },
                nextButton: {
                    icon: 'i-mdi-chevron-right',
                },
            },
        },
        card: {
            header: {
                padding: 'px-2 py-3 sm:px-4',
            },
            body: {
                padding: 'px-2 py-3 sm:px-4',
            },
            footer: {
                padding: 'px-2 py-3 sm:px-4',
            },
        },
        alert: {
            body: { padding: 'px-2 py-2 sm:p-2' },
            header: { padding: 'px-2 py-2 sm:p-2' },
            footer: { padding: 'px-2 py-2 sm:p-2' },
        },
        content: {
            search: {
                fileIcon: {
                    name: 'i-mdi-file-document',
                },
            },
            surround: {
                icon: {
                    prev: 'i-ph-arrow-left',
                    next: 'i-ph-arrow-right',
                },
            },
            prose: {
                code: {
                    button: {
                        icon: {
                            copy: 'i-mdi-content-copy',
                            copied: 'i-mdi-checkbox-marked',
                        },
                    },
                    icon: {
                        terminal: 'i-mdi-console',
                    },
                },
            },
        },

        // Nuxt UI Pro
        icons: {
            // Icons
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
        dashboard: {
            panel: {
                content: {
                    wrapper: 'pb-24 sm:pb-4',
                },
            },
        },
        page: {
            grid: {
                wrapper: 'gap-4',
            },
        },
    },
});
