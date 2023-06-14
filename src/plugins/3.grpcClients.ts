import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { RpcTransport } from '@protobuf-ts/runtime-rpc';
import {
    MethodInfo,
    NextServerStreamingFn,
    NextUnaryFn,
    RpcError,
    RpcInterceptor,
    RpcOptions,
    ServerStreamingCall,
    UnaryCall,
} from '@protobuf-ts/runtime-rpc/build/types';
import { Notification } from '~/composables/notification/interfaces/Notification.interface';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { AuthServiceClient } from '~~/gen/ts/services/auth/auth.client';
import { CitizenStoreServiceClient } from '~~/gen/ts/services/citizenstore/citizenstore.client';
import { CompletorServiceClient } from '~~/gen/ts/services/completor/completor.client';
import { DMVServiceClient } from '~~/gen/ts/services/dmv/vehicles.client';
import { DocStoreServiceClient } from '~~/gen/ts/services/docstore/docstore.client';
import { JobsServiceClient } from '~~/gen/ts/services/jobs/jobs.client';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';
import { NotificatorServiceClient } from '~~/gen/ts/services/notificator/notificator.client';
import { RectorServiceClient } from '~~/gen/ts/services/rector/rector.client';

export default defineNuxtPlugin(() => {
    return {
        provide: {
            grpc: new GRPCClients(),
        },
    };
});

export class GRPCClients {
    private authInterceptor: AuthInterceptor;
    private transport: GrpcWebFetchTransport;

    constructor() {
        this.authInterceptor = new AuthInterceptor();

        const { $loading } = useNuxtApp();

        // See https://github.com/jrapoport/grpc-web-devtools#grpc-web-interceptor-support
        const interceptors: RpcInterceptor[] = [this.authInterceptor, $loading];

        /* //@ts-ignore GRPCWeb Devtools only exist when the user has the extension installed
        const devInterceptors = window.__GRPCWEB_DEVTOOLS__;
        if (devInterceptors) {
            const { devToolsUnaryInterceptor, devToolsStreamInterceptor } = devInterceptors();

            devToolsUnaryInterceptor;
            devToolsStreamInterceptor;
        } */

        this.transport = new GrpcWebFetchTransport({
            baseUrl: '/grpc',
            format: 'text',
            interceptors: interceptors,
        });
    }

    getTransport(): RpcTransport {
        return this.transport;
    }

    // Handle GRPC errors
    async handleError(err: RpcError): Promise<boolean> {
        const notifications = useNotificationsStore();

        const { $loading } = useNuxtApp();

        const notification = {
            id: '',
            type: 'error',
            title: { key: 'notifications.grpc_errors.internal.title', parameters: [] },
            content: { key: err.message, parameters: [] },
        } as Notification;

        switch (err.code.toLowerCase()) {
            case 'internal':
                if (err.message.startsWith('errors.')) {
                    const errSplits = err.message.split(';');
                    if (errSplits.length > 1) {
                        notification.title = { key: errSplits[1], parameters: [] };
                        notification.content = { key: errSplits[0], parameters: [] };
                    }
                }
                break;

            case 'unavailable':
                notification.title = { key: 'notifications.grpc_errors.unavailable.title', parameters: [] };
                notification.content = { key: 'notifications.grpc_errors.unavailable.content', parameters: [] };
                break;

            case 'unauthenticated':
                await useAuthStore().clearAuthInfo();

                notification.type = 'warning';
                notification.title = { key: 'notifications.grpc_errors.unauthenticated.title', parameters: [] };
                notification.content = { key: 'notifications.grpc_errors.unauthenticated.content', parameters: [] };

                // Only update the redirect query param if it isn't already set
                const route = useRoute();
                const redirect = route.query.redirect ?? route.fullPath;
                navigateTo({
                    name: 'auth-login',
                    query: { redirect: redirect },
                    replace: true,
                    force: true,
                });
                break;

            case 'permission_denied':
                notification.title = { key: 'notifications.grpc_errors.permission_denied.title', parameters: [] };
                break;

            case 'not_found':
                notification.title = { key: 'notifications.grpc_errors.not_found.title', parameters: [] };
                break;

            default:
                notification.title = { key: 'notifications.grpc_errors.default.title', parameters: [] };
                notification.content = {
                    key: 'notifications.grpc_errors.default.content',
                    parameters: [err.message, err.code.valueOf()],
                };
                break;
        }

        notifications.dispatchNotification({
            type: notification.type,
            title: notification.title,
            content: notification.content,
        });

        $loading.errored();
        return true;
    }

    // GRPC Clients ===============================================================
    // Account / Auth - Unauthorized and authorized clients
    getUnAuthClient(): AuthServiceClient {
        return new AuthServiceClient(
            new GrpcWebFetchTransport({
                baseUrl: '/grpc',
                format: 'text',
            })
        );
    }

    getAuthClient(): AuthServiceClient {
        return new AuthServiceClient(this.transport);
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

    // Job
    getJobsClient(): JobsServiceClient {
        return new JobsServiceClient(this.transport);
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
}

export class AuthInterceptor implements RpcInterceptor {
    interceptUnary(next: NextUnaryFn, method: MethodInfo, input: object, options: RpcOptions): UnaryCall {
        if (!options.meta) {
            options.meta = {};
        }

        const { accessToken } = useAuthStore();
        if (accessToken !== null) {
            options.meta['Authorization'] = 'Bearer ' + accessToken;
        }

        return next(method, input, options);
    }

    interceptServerStreaming?(
        next: NextServerStreamingFn,
        method: MethodInfo,
        input: object,
        options: RpcOptions
    ): ServerStreamingCall {
        if (!options.meta) {
            options.meta = {};
        }

        const { accessToken } = useAuthStore();
        if (accessToken !== null) {
            options.meta['Authorization'] = 'Bearer ' + accessToken;
        }

        return next(method, input, options);
    }
}
