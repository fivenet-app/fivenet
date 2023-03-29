type ArpanetConfig = {
    sentryDSN: string;
    apiProtoURL: string;
};

const config: ArpanetConfig = {
    sentryDSN: '',
    apiProtoURL: '',
};

export default config;

export async function loadConfig() {
    let url = '/api/config';
    await fetch(url, {
        method: 'POST',
    })
        .then((response) => response.json())
        .then((data) => {
            data = data as ArpanetConfig;
            config.sentryDSN = data.sentryDSN;
            config.apiProtoURL = data.apiProtoURL;
        }).catch((err) => {
            console.log("Failed to get aRPaNet config from server!");
        });
}
