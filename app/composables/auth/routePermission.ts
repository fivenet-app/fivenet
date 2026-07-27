import { titleCase } from 'scule';
import type { RouteLocationNormalized, RouteLocationNormalizedLoaded } from 'vue-router';
import type { Notification } from '~/types/notifications';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

type RouteWithPermission = Pick<RouteLocationNormalized | RouteLocationNormalizedLoaded, 'meta' | 'name' | 'path'>;

export function canAccessRoute(route: RouteWithPermission): boolean {
    if (!route.meta.permission) return true;

    const { can } = useAuth();
    return can(route.meta.permission).value;
}

export function getRoutePermissionDeniedNotification(route: RouteWithPermission): Notification {
    return {
        title: { key: 'notifications.auth.no_permission.title', parameters: {} },
        description: {
            key: 'notifications.auth.no_permission.content',
            parameters: {
                path: route.name ? titleCase(route.name.toString().replaceAll('-', ' ')) + ` (${route.path})` : route.path,
            },
        },
        type: NotificationType.WARNING,
        actions: [],
    };
}

export async function revalidateCurrentRoutePermission(addNotification: (notification: Notification) => void): Promise<void> {
    const route = useRoute();

    if (canAccessRoute(route)) return;

    addNotification(getRoutePermissionDeniedNotification(route));

    const { username } = useAuth();
    if (username.value !== null) {
        await navigateTo({
            name: 'overview',
        });
    }
}
