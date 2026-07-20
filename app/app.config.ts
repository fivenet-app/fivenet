import type { Auth, Discord, FeatureGates, Game, System } from '../gen/ts/resources/clientconfig/clientconfig';
import type { Display, QuickButtons, Website } from '../gen/ts/resources/settings/config';
import { DataMode, type Data } from '../gen/ts/resources/settings/data';

export default defineAppConfig({
    // Server provided App Config
    version: '',

    defaultLocale: 'en',

    auth: {
        signupEnabled: true,
        lastCharLock: false,
        providers: [],
    } as Auth,
    discord: {
        botEnabled: false,
    } as Discord,
    website: {
        links: undefined,
    } as Website,
    featureGates: {} as FeatureGates,
    game: {
        unemployedJobName: 'unemployed',
        startJobGrade: 0,

        livemap: {
            enableCayoPerico: true,
        },

        maxWantedDurationUserEnabled: false,
        maxWantedDurationVehicleEnabled: false,
    } as Game,
    system: {
        bannerMessageEnabled: false,

        otlp: {
            enabled: false,
            headers: {},
            url: '',
        },
    } as System,
    display: {
        intlLocale: 'en-US',
        currencyName: 'USD',
    } as Display,
    quickButtons: {
        penaltyCalculator: {},
    } as QuickButtons,
    data: {
        mode: DataMode.UNAVAILABLE,
    } as Data,

    maxContentLength: 40_000, // Characters

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
        loadWaitTime: 850,
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

    // Nuxt UI config
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
            taupe: 'taupe',
            mauve: 'mauve',
            mist: 'mist',
            olive: 'olive',
        },

        icons: {
            arrowDown: 'i-mdi-arrow-down',
            arrowLeft: 'i-mdi-arrow-left',
            arrowRight: 'i-mdi-arrow-right',
            arrowUp: 'i-mdi-arrow-up',
            caution: 'i-mdi-alert-circle',
            check: 'i-mdi-check',
            chevronDoubleLeft: 'i-mdi-chevron-double-left',
            chevronDoubleRight: 'i-mdi-chevron-double-right',
            chevronDown: 'i-mdi-chevron-down',
            chevronLeft: 'i-mdi-chevron-left',
            chevronRight: 'i-mdi-chevron-right',
            chevronUp: 'i-mdi-chevron-up',
            close: 'i-mdi-close',
            copy: 'i-mdi-content-copy',
            copyCheck: 'i-mdi-check-circle-outline',
            dark: 'i-mdi-moon-waning-crescent',
            drag: 'i-mdi-drag-vertical',
            ellipsis: 'i-mdi-dots-horizontal',
            error: 'i-mdi-close-circle',
            external: 'i-mdi-arrow-top-right',
            eye: 'i-mdi-eye',
            eyeOff: 'i-mdi-eye-off',
            file: 'i-mdi-file-document',
            folder: 'i-mdi-folder',
            folderOpen: 'i-mdi-folder-open',
            hash: 'i-mdi-pound',
            info: 'i-mdi-information',
            light: 'i-mdi-white-balance-sunny',
            loading: 'i-mdi-loading',
            menu: 'i-mdi-menu',
            minus: 'i-mdi-minus',
            panelClose: 'i-mdi-menu-close',
            panelOpen: 'i-mdi-menu-open',
            plus: 'i-mdi-plus',
            reload: 'i-mdi-reload',
            search: 'i-mdi-magnify',
            star: 'i-mdi-star-outline',
            stop: 'i-mdi-stop',
            success: 'i-mdi-check-circle',
            system: 'i-mdi-monitor',
            tip: 'i-mdi-lightbulb-variant',
            upload: 'i-mdi-upload',
            warning: 'i-mdi-alert',
        },

        alert: {
            slots: {
                root: 'p-2',
            },
        },

        avatar: {
            slots: {
                root: 'rounded-md',
            },
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

        modal: {
            variants: {
                fullscreen: {
                    false: {
                        content: 'max-w-3xl',
                    },
                },
            },
        },

        slideover: {
            variants: {
                side: {
                    right: {
                        content: 'max-w-xl',
                    },
                    left: {
                        content: 'max-w-xl',
                    },
                },
            },
        },

        table: {
            slots: {
                tr: 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-default',
                th: 'px-4 py-1.5',
                td: 'p-1.5',
            },
        },

        pageGrid: {
            base: 'gap-4',
        },
    },
});
