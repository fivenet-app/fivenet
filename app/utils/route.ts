import type { NuxtRoute, RoutesNamesList } from '@typed-router';

export function unsafeRoute(path: string) {
    return path as NuxtRoute<RoutesNamesList, string, false>;
}
