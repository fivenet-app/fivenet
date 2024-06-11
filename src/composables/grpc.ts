import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { GrpcWSTransport } from '~/composables/grpcws/bridge';
import { type Notification } from '~/composables/notifications';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AuthServiceClient } from '~~/gen/ts/services/auth/auth.client';
import { CalendarServiceClient } from '~~/gen/ts/services/calendar/calendar.client';
import { CentrumServiceClient } from '~~/gen/ts/services/centrum/centrum.client';
import { CitizenStoreServiceClient } from '~~/gen/ts/services/citizenstore/citizenstore.client';
import { CompletorServiceClient } from '~~/gen/ts/services/completor/completor.client';
import { DMVServiceClient } from '~~/gen/ts/services/dmv/vehicles.client';
import { DocStoreServiceClient } from '~~/gen/ts/services/docstore/docstore.client';
import { JobsConductServiceClient } from '~~/gen/ts/services/jobs/conduct.client';
import { JobsServiceClient } from '~~/gen/ts/services/jobs/jobs.client';
import { JobsTimeclockServiceClient } from '~~/gen/ts/services/jobs/timeclock.client';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';
import { MessengerServiceClient } from '~~/gen/ts/services/messenger/messenger.client';
import { NotificatorServiceClient } from '~~/gen/ts/services/notificator/notificator.client';
import { QualificationsServiceClient } from '~~/gen/ts/services/qualifications/qualifications.client';
import { RectorConfigServiceClient } from '~~/gen/ts/services/rector/config.client';
import { RectorFilestoreServiceClient } from '~~/gen/ts/services/rector/filestore.client';
import { RectorLawsServiceClient } from '~~/gen/ts/services/rector/laws.client';
import { RectorServiceClient } from '~~/gen/ts/services/rector/rector.client';

const grpcWebTransport = new GrpcWebFetchTransport({
    baseUrl: '/api/grpc',
    format: 'text',
    fetchInit: {
        credentials: 'same-origin',
    },
    timeout: 8500,
});

const grpcWSTransport = new GrpcWSTransport({
    url: `${window.location.protocol}//${window.location.hostname}:${!import.meta.dev ? window.location.port : 8080}/api/grpc`,
    interceptors: [],
    debug: import.meta.dev,
});

// Handle GRPC errors
export async function handleGRPCError(err: RpcError): Promise<boolean> {
    const notifications = useNotificatorStore();

    const notification = {
        id: '',
        type: NotificationType.ERROR,
        title: { key: 'notifications.grpc_errors.internal.title', parameters: {} },
        description: { key: err?.message ?? 'Unknown error', parameters: {} },
    } as Notification;

    const traceId = (err?.meta && (err?.meta['trailer+x-trace-id'] as string)) ?? 'UNKNOWN';

    if (err.code !== undefined) {
        const route = useRoute();
        switch (err.code.toLowerCase()) {
            case 'internal':
                break;

            case 'unavailable':
                notification.title = { key: 'notifications.grpc_errors.unavailable.title', parameters: {} };
                notification.description = { key: 'notifications.grpc_errors.unavailable.content', parameters: {} };
                break;

            case 'unauthenticated':
                useAuthStore().clearAuthInfo();

                notification.type = NotificationType.WARNING;
                notification.title = { key: 'notifications.grpc_errors.unauthenticated.title', parameters: {} };
                notification.description = { key: 'notifications.grpc_errors.unauthenticated.content', parameters: {} };

                // Only update the redirect query param if it isn't already set
                navigateTo({
                    name: 'auth-login',
                    query: { redirect: route.query.redirect ?? route.fullPath },
                    replace: true,
                    force: true,
                });
                break;

            case 'permission_denied':
                if (!err.message.includes('ErrCharLock')) {
                    notification.title = { key: 'notifications.grpc_errors.permission_denied.title', parameters: {} };
                } else {
                    // In case of a permission denied char lock error, user must re-select their char
                    navigateTo({
                        name: 'auth-character-selector',
                        query: { redirect: route.query.redirect ?? route.fullPath },
                        replace: true,
                        force: true,
                    });
                }
                break;

            case 'not_found':
                notification.title = { key: 'notifications.grpc_errors.not_found.title', parameters: {} };
                break;

            default:
                notification.title = { key: 'notifications.grpc_errors.default.title', parameters: {} };
                notification.description = {
                    key: 'notifications.grpc_errors.default.content',
                    parameters: { msg: err.message, code: err.code.valueOf() },
                };
                notification.onClick = async () =>
                    copyToClipboardWrapper(
                        `## Error occured at ${new Date().toLocaleDateString()}:
**Service/Method**: ${err.serviceName}/${err.methodName} => ${err.code}
**Message**: ${err.message}
**TraceID**: ${traceId}`,
                    );
                break;
        }
    }

    console.error(
        `GRPC Client: Failed request ${err.serviceName}/${err.methodName} (Trace ID: '${traceId}', Message: ${err.message}`,
    );

    if (err.message?.startsWith('errors.')) {
        const errSplits = err.message.split(';');
        if (errSplits.length > 1) {
            notification.title = { key: errSplits[0], parameters: {} };
            notification.description = { key: errSplits[1], parameters: {} };
        } else {
            notification.description = { key: err.message, parameters: {} };
        }
    }

    notifications.add({
        type: notification.type,
        title: notification.title,
        description: notification.description,
        onClick: notification.onClick,
        onClickText: { key: 'pages.error.copy_error', parameters: {} },
    });

    return true;
}

// GRPC Clients ===============================================================
// Auth
export function getGRPCAuthClient(): AuthServiceClient {
    return new AuthServiceClient(grpcWebTransport);
}

// Centrum
export function getGRPCCentrumClient(): CentrumServiceClient {
    return new CentrumServiceClient(grpcWSTransport);
}

// Citizens
export function getGRPCCitizenStoreClient(): CitizenStoreServiceClient {
    return new CitizenStoreServiceClient(grpcWSTransport);
}

// Completion
export function getGRPCCompletorClient(): CompletorServiceClient {
    return new CompletorServiceClient(grpcWSTransport);
}

// DMV (Vehicles)
export function getGRPCDMVClient(): DMVServiceClient {
    return new DMVServiceClient(grpcWSTransport);
}

// Documents
export function getGRPCDocStoreClient(): DocStoreServiceClient {
    return new DocStoreServiceClient(grpcWSTransport);
}

// Livemap
export function getGRPCLivemapperClient(): LivemapperServiceClient {
    return new LivemapperServiceClient(grpcWSTransport);
}

// Notificator
export function getGRPCNotificatorClient(): NotificatorServiceClient {
    return new NotificatorServiceClient(grpcWSTransport);
}

// Rector
export function getGRPCRectorClient(): RectorServiceClient {
    return new RectorServiceClient(grpcWSTransport);
}

export function getGRPCRectorConfigClient(): RectorConfigServiceClient {
    return new RectorConfigServiceClient(grpcWSTransport);
}

export function getGRPCRectorFilestoreClient(): RectorFilestoreServiceClient {
    return new RectorFilestoreServiceClient(grpcWSTransport);
}

export function getGRPCRectorLawsClient(): RectorLawsServiceClient {
    return new RectorLawsServiceClient(grpcWSTransport);
}

// Jobs
export function getGRPCJobsClient(): JobsServiceClient {
    return new JobsServiceClient(grpcWSTransport);
}

export function getGRPCJobsConductClient(): JobsConductServiceClient {
    return new JobsConductServiceClient(grpcWSTransport);
}

export function getGRPCJobsTimeclockClient(): JobsTimeclockServiceClient {
    return new JobsTimeclockServiceClient(grpcWSTransport);
}

// Qualifications
export function getGRPCQualificationsClient(): QualificationsServiceClient {
    return new QualificationsServiceClient(grpcWSTransport);
}

// Calendar
export function getGRPCCalendarClient(): CalendarServiceClient {
    return new CalendarServiceClient(grpcWSTransport);
}

// Messenger
export function getGRPCMessengerClient(): MessengerServiceClient {
    return new MessengerServiceClient(grpcWSTransport);
}
