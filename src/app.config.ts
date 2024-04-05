import type { DiscordConfig, FeatureGates, Links, LoginConfig } from '~/shims';

export default defineAppConfig({
    version: '',
    login: {
        signupEnabled: true,
        providers: [],
    } as LoginConfig,
    discord: {
        botInviteURL: '',
    } as DiscordConfig,
    links: {} as Links,
    featureGates: {} as FeatureGates,

    // Nuxt UI app config
    ui: {
        primary: 'sky',
        gray: 'cool',
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
        icons: {
            dynamic: true,
        },
        table: {
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
    },
});
