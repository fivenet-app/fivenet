import type { AppConfig } from '~/shims';

export default defineAppConfig({
    version: '',
    sentryDSN: '',
    login: {
        signupEnabled: true,
        providers: [],
    },
    discord: {
        botInviteURL: '',
    },
    links: {},
} as AppConfig);
