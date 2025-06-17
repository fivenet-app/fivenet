import { useAuthStore } from '~/stores/auth';
import type { Notification } from '~/utils/notifications';
import type { Error as CommonError } from '~~/gen/ts/resources/common/error';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const logger = useLogger('📡 GRPC');

const throttledErrorCodes = ['internal', 'deadline_exceeded', 'cancelled', 'permission_denied', 'unauthenticated'];

const lastError: { receivedAt: undefined | Date; code: undefined | string } = {
    receivedAt: undefined,
    code: undefined,
};

function addCopyActionToNotification(notification: Notification, err: RpcError, traceId: number): void {
    notification.actions?.push({
        label: { key: 'pages.error.copy_error' },
        click: async () => {
            copyToClipboardWrapper(
                `## Error occured at ${new Date().toISOString()}:
**Service/Method**: \`${err.serviceName}/${err.methodName}\` => \`${err.code}\`
**Message**: \`${err.message}\`
**TraceID**: \`${traceId}\``,
            );

            const notifications = useNotificationsStore();
            notifications.add({
                title: { key: 'notifications.clipboard.copied.title', parameters: {} },
                description: { key: 'notifications.clipboard.copied.content', parameters: {} },
                timeout: 3250,
                type: NotificationType.INFO,
            });
        },
    });
}

// Handle GRPC errors
export async function handleGRPCError(err: RpcError | undefined): Promise<boolean> {
    if (err === undefined) {
        return true;
    }

    const notification = {
        id: 0,
        type: NotificationType.ERROR,
        title: { key: 'notifications.grpc_errors.internal.title', parameters: {} },
        description: {
            key: 'notifications.grpc_errors.internal.content',
            parameters: { msg: err?.message ?? 'Unknown error' },
        },
        actions: [],
    } as Notification;

    const traceId =
        (err?.meta &&
            err.meta['trailer+x-trace-id'] &&
            (typeof err.meta['trailer+x-trace-id'] === 'string'
                ? err.meta['trailer+x-trace-id']
                : (err.meta['trailer+x-trace-id'] as string[]).join(','))) ??
        'N/A';

    const code = err.code?.toUpperCase();
    if (code !== undefined) {
        const route = useRoute();

        // If the error code has already been "handled", skip "handling" them for now
        // Only do this for internal, deadline_exceeded, cancelled, permission_denied, unauthenticated codes
        if (throttledErrorCodes.includes(code)) {
            if (lastError.code === code) {
                if (lastError.receivedAt !== undefined && lastError.receivedAt.getTime() - new Date().getTime() < 10) {
                    return true;
                }
            }
            lastError.code = code;
            lastError.receivedAt = new Date();
        }

        switch (code) {
            case 'INTERNAL':
                break;

            case 'DEADLINE_EXCEEDED':
                notification.title = { key: 'notifications.grpc_errors.deadline_exceeded.title', parameters: {} };
                notification.description = { key: 'notifications.grpc_errors.deadline_exceeded.content', parameters: {} };
                addCopyActionToNotification(notification, err, traceId);
                break;

            case 'CANCELLED':
                // Don't send notifications for cancelled requests
                return true;

            case 'UNAVAILABLE':
                notification.title = { key: 'notifications.grpc_errors.unavailable.title', parameters: {} };
                notification.description = { key: 'notifications.grpc_errors.unavailable.content', parameters: {} };
                break;

            case 'UNAUTHENTICATED':
                useAuthStore().clearAuthInfo();

                notification.type = NotificationType.WARNING;
                notification.title = { key: 'notifications.grpc_errors.unauthenticated.title', parameters: {} };
                notification.description = { key: 'notifications.grpc_errors.unauthenticated.content', parameters: {} };

                // Only update the redirect query param if it isn't already set
                await navigateTo({
                    name: 'auth-login',
                    query: { redirect: route.query.redirect ?? route.fullPath },
                    replace: true,
                    force: true,
                });
                break;

            case 'PERMISSION_DENIED':
                if (!err.message.includes('ErrCharLock')) {
                    notification.title = { key: 'notifications.grpc_errors.permission_denied.title', parameters: {} };
                } else {
                    // In case of a permission denied char lock error, user must re-select their char
                    await navigateTo({
                        name: 'auth-character-selector',
                        query: { redirect: route.query.redirect ?? route.fullPath },
                        replace: true,
                    });
                }
                break;

            case 'NOT_FOUND':
                notification.title = { key: 'notifications.grpc_errors.not_found.title', parameters: {} };
                break;

            default:
                notification.title = { key: 'notifications.grpc_errors.default.title', parameters: {} };
                notification.description = {
                    key: 'notifications.grpc_errors.default.content',
                    parameters: { msg: err.message, code: err.code.valueOf() },
                };
                addCopyActionToNotification(notification, err, traceId);
                break;
        }
    }

    logger.error(
        `Failed request ${err.serviceName}/${err.methodName} (Trace ID: '${traceId}', Code: ${err.code}, Message: ${err.message}`,
    );

    if (err?.message && isTranslatedError(err.message as string)) {
        const parsed = JSON.parse(err.message) as CommonError;

        if (parsed.title) {
            notification.title = parsed.title;
        }
        if (parsed.content) {
            notification.description = parsed.content;
        }
    }

    const notifications = useNotificationsStore();
    notifications.add({
        type: notification.type,
        title: notification.title,
        description: notification.description,
        actions: notification.actions,
    });

    return true;
}
