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
    },
});
