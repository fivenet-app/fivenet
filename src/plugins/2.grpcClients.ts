import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { type RpcTransport } from '@protobuf-ts/runtime-rpc';
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
import { NotificatorServiceClient } from '~~/gen/ts/services/notificator/notificator.client';
import { QualificationsServiceClient } from '~~/gen/ts/services/qualifications/qualifications.client';
import { RectorConfigServiceClient } from '~~/gen/ts/services/rector/config.client';
import { RectorFilestoreServiceClient } from '~~/gen/ts/services/rector/filestore.client';
import { RectorLawsServiceClient } from '~~/gen/ts/services/rector/laws.client';
import { RectorServiceClient } from '~~/gen/ts/services/rector/rector.client';

export default defineNuxtPlugin({
    name: 'grpcClients',
    parallel: true,
    async setup(_) {
        return {
            provide: {
                grpc: new GRPCClients(),
            },
        };
    },
});

export class GRPCClients {
    private transport: GrpcWebFetchTransport;

    constructor() {
        this.transport = new GrpcWebFetchTransport({
            baseUrl: '/grpc',
            format: 'text',
            interceptors: [],
            fetchInit: { credentials: 'same-origin' },
        });
    }

    getTransport(): RpcTransport {
        return this.transport;
    }

    // Handle GRPC errors
    async handleError(err: RpcError): Promise<boolean> {
        const notifications = useNotificatorStore();

        const notification = {
            id: '',
            type: NotificationType.ERROR,
            title: { key: 'notifications.grpc_errors.internal.title', parameters: {} },
            description: { key: err.message ?? 'Unknown error', parameters: {} },
        } as Notification;

        let traceId = 'UNKNOWN';
        if (err.meta !== undefined && err.meta['x-trace-id'] !== undefined) {
            traceId = err.meta['x-trace-id'] as string;
        }

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

        if (err.message.startsWith('errors.')) {
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
    getAuthClient(): AuthServiceClient {
        return new AuthServiceClient(this.transport);
    }

    // Centrum
    getCentrumClient(): CentrumServiceClient {
        return new CentrumServiceClient(this.transport);
    }

    // Citizens
    getCitizenStoreClient(): CitizenStoreServiceClient {
        return new CitizenStoreServiceClient(this.transport);
    }

    // Completion
    getCompletorClient(): CompletorServiceClient {
        return new CompletorServiceClient(this.transport);
    }

    // DMV (Vehicles)
    getDMVClient(): DMVServiceClient {
        return new DMVServiceClient(this.transport);
    }

    // Documents
    getDocStoreClient(): DocStoreServiceClient {
        return new DocStoreServiceClient(this.transport);
    }

    // Livemap
    getLivemapperClient(): LivemapperServiceClient {
        return new LivemapperServiceClient(this.transport);
    }

    // Notificator
    getNotificatorClient(): NotificatorServiceClient {
        return new NotificatorServiceClient(this.transport);
    }

    // Rector
    getRectorClient(): RectorServiceClient {
        return new RectorServiceClient(this.transport);
    }

    getRectorConfigClient(): RectorConfigServiceClient {
        return new RectorConfigServiceClient(this.transport);
    }

    getRectorFilestoreClient(): RectorFilestoreServiceClient {
        return new RectorFilestoreServiceClient(this.transport);
    }

    getRectorLawsClient(): RectorLawsServiceClient {
        return new RectorLawsServiceClient(this.transport);
    }

    // Jobs
    getJobsClient(): JobsServiceClient {
        return new JobsServiceClient(this.transport);
    }

    getJobsConductClient(): JobsConductServiceClient {
        return new JobsConductServiceClient(this.transport);
    }

    getJobsTimeclockClient(): JobsTimeclockServiceClient {
        return new JobsTimeclockServiceClient(this.transport);
    }

    // Qualifications
    getQualificationsClient(): QualificationsServiceClient {
        return new QualificationsServiceClient(this.transport);
    }

    // Calendar
    getCalendarClient(): CalendarServiceClient {
        return new CalendarServiceClient(this.transport);
    }
}
