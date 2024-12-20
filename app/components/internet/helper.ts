import HomePage from './pages/HomePage.vue';

export const urlHomePage = 'internet.search';

export const localPages = {
    'internet.search': HomePage,
};

export function splitURL(url: string): undefined | { domain: string; path?: string } {
    const split = url.split('/');
    if (split.length === 0) {
        return undefined;
    }
    const path = '/' + (split[1] ? split.slice(1).join('/') : '');

    return { domain: split[0]!, path: path };
}

export function joinURL(domain: string, path?: string): string {
    return domain + (path && path !== '' ? path : '/');
}
