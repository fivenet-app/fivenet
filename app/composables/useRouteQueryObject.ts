import type { RouteLocationNormalized } from 'vue-router';
import { useSettingsStore } from '~/stores/settings';

const _useRouteQueryObject = <T extends object>(key: string, def: T): Ref<T> => {
    const query = useRouteQuery<string, T>(key, JSON.stringify(def), {
        transform: {
            get(val) {
                try {
                    const parsed = JSON.parse(val) as T;
                    return {
                        ...def,
                        ...parsed,
                    };
                } catch (_) {
                    return def;
                }
            },
            set(val) {
                return JSON.stringify(val);
            },
        },
    });

    return query;
};

export const useRouteQueryObject = _useRouteQueryObject;

// List of excluded redirect urls
export const redirectExcludedPages = ['/auth/login', '/auth/logout', '/dereferer', '/auth/character-selector'];

export function getRedirect(route: RouteLocationNormalized): URL {
    const settingsStore = useSettingsStore();
    const { startpage } = storeToRefs(settingsStore);

    const redirect = (route.query.redirect ?? startpage.value ?? '/overview') as string;
    const path = checkRedirectPathValid(redirect) ? redirect : startpage.value;

    return parseRedirectURL(path);
}

export function parseRedirectURL(path: string): URL {
    return new URL('https://example.com' + path);
}

export function checkRedirectPathValid(path: string): boolean {
    return !redirectExcludedPages.some((p) => path.startsWith(p)) && path !== '/';
}

export function getRedirectPath(path: string): string {
    return checkRedirectPathValid(path) ? path : '/overview';
}

export function getRedirectURL(path: string): URL {
    return parseRedirectURL(getRedirectPath(path));
}
