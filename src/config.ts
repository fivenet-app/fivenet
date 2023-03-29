type ArpanetConfig = {
    sentryDSN: string;
    apiProtoURL: string;
};

const config: ArpanetConfig = {
    sentryDSN: '',
    apiProtoURL: '',
};

export default config;

export async function loadConfig(): Promise<void> {
    let url = import.meta.env.DEV ? 'http://localhost:8080/api/config' : '/api/config';

    try {
        const response = await fetch(url, {
            method: 'POST',
        });
        const data = (await response.json()) as ArpanetConfig;
        config.sentryDSN = data.sentryDSN;
        config.apiProtoURL = data.apiProtoURL;
    } catch (_) {
        console.error('Failed to get aRPaNet config from server');
    }
}
