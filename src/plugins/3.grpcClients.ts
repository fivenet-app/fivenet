import { useAuthStore } from '~/store/auth';
import { AuthServiceClient } from '~~/gen/ts/services/auth/auth.client';
import { CitizenStoreServiceClient } from '~~/gen/ts/services/citizenstore/citizenstore.client';
import { CompletorServiceClient } from '~~/gen/ts/services/completor/completor.client';
import { DMVServiceClient } from '~~/gen/ts/services/dmv/vehicles.client';
import { DocStoreServiceClient } from '~~/gen/ts/services/docstore/docstore.client';
import { JobsServiceClient } from '~~/gen/ts/services/jobs/jobs.client';
import { LivemapperServiceClient } from '~~/gen/ts/services/livemapper/livemap.client';
import { RpcError, StatusCode } from 'grpc-web';
import { RectorServiceClient } from '~~/gen/ts/services/rector/rector.client';
import { NotificatorServiceClient } from '~~/gen/ts/services/notificator/notificator.client';
import { useNotificationsStore } from '~/store/notifications';
import { RpcTransport } from '@protobuf-ts/runtime-rpc';
import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import {
    NextServerStreamingFn,
    NextUnaryFn,
    MethodInfo,
    RpcOptions,
    RpcInterceptor,
    UnaryCall,
    ServerStreamingCall,
    RpcStatus,
} from '@protobuf-ts/runtime-rpc/build/types';

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

        switch (err.code) {
            case StatusCode.UNAUTHENTICATED:
                await useAuthStore().clearAuthInfo();

                notifications.dispatchNotification({
                    title: 'notifications.grpc_errors.unauthenticated.title',
                    titleI18n: true,
                    content: 'notifications.grpc_errors.unauthenticated.content',
                    contentI18n: true,
                    type: 'warning',
                });

                // Only update the redirect query param if it isn't already set
                const route = useRoute();
                const redirect = route.query.redirect ?? route.fullPath;
                await navigateTo({ name: 'auth-login', query: { redirect: redirect }, replace: true, force: true });
                break;

            case StatusCode.PERMISSION_DENIED:
                notifications.dispatchNotification({
                    title: 'notifications.grpc_errors.permission_denied.title',
                    titleI18n: true,
                    content: err.message,
                    type: 'error',
                });
                break;

            case StatusCode.INTERNAL:
                let title = 'notifications.grpc_errors.internal.title';
                let content = err.message;
                let contentI18n = false;
                if (err.message.startsWith('errors.')) {
                    contentI18n = true;

                    const errSplits = err.message.split(';');
                    if (errSplits.length > 1) {
                        title = errSplits[1];
                        content = errSplits[0];
                    }
                }

                notifications.dispatchNotification({
                    title: title,
                    titleI18n: true,
                    content: content,
                    contentI18n: contentI18n,
                    type: 'error',
                });
                break;

            case StatusCode.UNAVAILABLE:
                notifications.dispatchNotification({
                    title: 'notifications.grpc_errors.unavailable.title',
                    titleI18n: true,
                    content: 'notifications.grpc_errors.unavailable.content',
                    contentI18n: true,
                    type: 'error',
                });
                break;

            case StatusCode.NOT_FOUND:
                notifications.dispatchNotification({
                    title: 'notifications.grpc_errors.unavailable.title',
                    titleI18n: true,
                    content: err.message,
                    type: 'error',
                });

                break;

            default:
                notifications.dispatchNotification({
                    title: 'notifications.grpc_errors.default.title',
                    titleI18n: true,
                    content: err.message,
                    type: 'error',
                });
                break;
        }

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
