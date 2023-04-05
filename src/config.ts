type ArpanetConfig = {
    sentryDSN: string;
    apiProtoURL: string;
};

const config: ArpanetConfig = {
    sentryDSN: '',
    apiProtoURL: '/grpc',
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
    } catch (_) {
        console.error('Failed to get FiveNet config from server');
    }
}
