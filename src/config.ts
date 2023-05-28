type ArpanetConfig = {
    sentryDSN: string;
    login: LoginConfig;
};

type LoginConfig = {
    providers: ProviderConfig[];
};

type ProviderConfig = {
    name: string;
    label: string;
};

const config: ArpanetConfig = {
    sentryDSN: '',
    login: {
        providers: [],
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
            const data = (await resp.json()) as ArpanetConfig;
            config.sentryDSN = data.sentryDSN;
            config.login = data.login;

            return res();
        } catch (e) {
            return rej(e);
        }
    });
}
