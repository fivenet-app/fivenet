export const setupBypassRoutes = ['/auth/logout', '/auth/account-info', '/auth/character-selector'] as const;

export function isSetupBypassRoute(path: string): boolean {
    return setupBypassRoutes.some((route) => path.startsWith(route));
}
