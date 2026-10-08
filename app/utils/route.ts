import type { NuxtRoute, RoutesNamesList } from '@typed-router';

/**
 * Checks whether a route path is the given path or a child of it.
 */
export function isRoute(path: string, routePath: string): boolean {
    const normalize = (value: string): string => (value.length > 1 ? value.replace(/\/+$/, '') : value);
    const normalizedPath = normalize(path);
    const normalizedRoutePath = normalize(routePath);

    return normalizedPath === normalizedRoutePath || normalizedPath.startsWith(`${normalizedRoutePath}/`);
}

export function unsafeRoute(path: string) {
    return path as NuxtRoute<RoutesNamesList, string, false>;
}
