type ArpanetConfig = {
    sentryDSN: string;
    apiProtoURL: string;
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
    apiProtoURL: '/grpc',
    login: {
        providers: [],
    },
};

export default config;

export async function loadConfig(): Promise<void> {
    try {
        const response = await fetch('/api/config', {
            method: 'POST',
        });
        const data = (await response.json()) as ArpanetConfig;
        config.sentryDSN = data.sentryDSN;
        config.apiProtoURL = data.apiProtoURL;
        config.login = data.login;
    } catch (_) {
        console.error('Failed to get FiveNet config from server');
    }
}
