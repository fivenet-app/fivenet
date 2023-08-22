type AppConfig = {
    sentryDSN: string;
    login: LoginConfig;
    discord: DiscordConfig;
};

type LoginConfig = {
    signupEnabled: boolean;
    providers: ProviderConfig[];
};

type ProviderConfig = {
    name: string;
    label: string;
};

type DiscordConfig = {
    botInviteURL?: string;
};

const config: AppConfig = {
    sentryDSN: '',
    login: {
        signupEnabled: true,
        providers: [],
    },
    discord: {
        botInviteURL: '',
    },
};

export default config;

export async function loadConfig(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const resp = await fetch('/api/config', {
                method: 'POST',
            });
            if (!resp.ok) {
                const text = await resp.text();
                throw createError({
                    statusCode: 500,
                    statusMessage: 'Failed to get FiveNet config from backend',
                    message: text,
                    fatal: true,
                    unhandled: false,
                });
            }
            const data = (await resp.json()) as AppConfig;
            config.sentryDSN = data.sentryDSN;
            config.login = data.login;

            return res();
        } catch (e) {
            return rej(e);
        }
    });
}

type ClientConfig = {
    NUIEnabled: boolean;
    NUIResourceName?: string;
};

export const clientConfig: ClientConfig = {
    NUIEnabled: false,
};
