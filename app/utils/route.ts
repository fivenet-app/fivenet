import type { NuxtRoute, RoutesNamesList } from '@typed-router';

/**
 * Checks whether a route path is the given path or a child of it.
 */
export function isRoute(path: string, routePath: string): boolean {
    return path === routePath || path.startsWith(`${routePath}/`);
}

export function unsafeRoute(path: string) {
    return path as NuxtRoute<RoutesNamesList, string, false>;
}
