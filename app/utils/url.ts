export function generateDerefURL(target: string): string {
    return (
        '/dereferer?' +
        new URLSearchParams({
            target: target,
            source: window.location.href,
        })
    );
}

export function generateDiscordConnectURL(provider: string, redirect?: string, params?: Record<string, string>): string {
    const url = new URL(`https://example.com/api/oauth2/login/${provider}`);

    if (params != undefined) Object.keys(params).forEach((key) => url.searchParams.set(key, params[key]!));

    url.searchParams.set('connect-only', 'true');
    if (redirect !== undefined) {
        url.searchParams.set('redirect', redirect);
    }

    return url.pathname + url.search;
}
