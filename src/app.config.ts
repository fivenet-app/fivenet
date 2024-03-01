import type { AppConfig } from '~/shims';

export default defineAppConfig({
    version: '',
    login: {
        signupEnabled: true,
        providers: [],
    },
    discord: {
        botInviteURL: '',
    },
    links: {},
    featureGates: {},
} as AppConfig);
