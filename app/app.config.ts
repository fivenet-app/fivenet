import type { Discord, FeatureGates, Game, LoginConfig, System, Website } from '~~/gen/ts/resources/clientconfig/clientconfig';

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
        botEnabled: false,
    } as Discord,
    website: {
        links: undefined,
    } as Website,
    featureGates: {
        imageProxy: false,
    } as FeatureGates,
    game: {
        unemployedJobName: 'unemployed',
        startJobGrade: 0,
    } as Game,
    system: {
        bannerMessageEnabled: false,

        otlp: {
            enabled: false,
        },
    } as System,

    // File upload related config
    fileUpload: {
        fileSizes: {
            fileStore: 5 * 1024 * 1024,
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

    custom: {
        icons: {
            // Custom Icons
            sort: 'i-mdi-sort',
            sortAsc: 'i-mdi-sort-ascending',
            sortDesc: 'i-mdi-sort-descending',
        },
    },

    // Nuxt UI and UI Pro config
    ui: {
        colors: {
            primary: 'blue',
            secondary: 'neutral',
            info: 'blue',
            success: 'green',
            warning: 'yellow',
            error: 'red',
            // Palette colors
            amber: 'amber',
            blue: 'blue',
            cyan: 'cyan',
            emerald: 'emerald',
            fuchsia: 'fuchsia',
            green: 'green',
            indigo: 'indigo',
            lime: 'lime',
            orange: 'orange',
            pink: 'pink',
            purple: 'purple',
            red: 'red',
            rose: 'rose',
            sky: 'sky',
            teal: 'teal',
            violet: 'violet',
            white: 'white',
            yellow: 'yellow',
            // Gray Colors
            gray: 'gray',
            neutral: 'neutral',
            slate: 'slate',
            stone: 'stone',
            zinc: 'zinc',
        },

        icons: {
            arrowLeft: 'mdi-arrow-left',
            arrowRight: 'mdi-arrow-right',
            check: 'mdi-check',
            chevronDoubleLeft: 'mdi-chevron-double-left',
            chevronDoubleRight: 'mdi-chevron-double-right',
            chevronDown: 'mdi-chevron-down',
            chevronLeft: 'mdi-chevron-left',
            chevronRight: 'mdi-chevron-right',
            chevronUp: 'mdi-chevron-up',
            close: 'mdi-close',
            ellipsis: 'mdi-dots-horizontal',
            external: 'mdi-arrow-top-right',
            file: 'mdi-file-document',
            folder: 'mdi-folder',
            folderOpen: 'mdi-folder-open',
            loading: 'mdi-loading',
            minus: 'mdi-minus',
            plus: 'mdi-plus',
            search: 'mdi-magnify',
            upload: 'mdi-upload',
            arrowUp: 'mdi-arrow-up',
            arrowDown: 'mdi-arrow-down',
            caution: 'mdi-alert-circle',
            copy: 'mdi-content-copy',
            copyCheck: 'mdi-check-circle-outline',
            dark: 'mdi-moon-waning-crescent',
            error: 'mdi-close-circle',
            eye: 'mdi-eye',
            eyeOff: 'mdi-eye-off',
            hash: 'mdi-pound',
            info: 'mdi-information',
            light: 'mdi-white-balance-sunny',
            menu: 'mdi-menu',
            panelClose: 'mdi-menu-close',
            panelOpen: 'mdi-menu-open',
            reload: 'mdi-reload',
            stop: 'mdi-stop',
            success: 'mdi-check-circle',
            system: 'mdi-monitor',
            tip: 'mdi-lightbulb-variant',
            warning: 'mdi-alert',
        },

        inputMenu: {
            slots: {
                content: 'min-w-fit',
            },
        },

        selectMenu: {
            slots: {
                content: 'min-w-fit',
            },
        },

        table: {
            slots: {
                tr: 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
                th: 'px-4 py-1.5',
                td: 'p-1.5',
            },
        },
    },
});
