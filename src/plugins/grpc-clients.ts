import { useAuthStore } from '../store/auth';
import config from '../config';
import { AuthServiceClient } from '@fivenet/gen/services/auth/AuthServiceClientPb';
import { CitizenStoreServiceClient } from '@fivenet/gen/services/citizenstore/CitizenstoreServiceClientPb';
import { CompletorServiceClient } from '@fivenet/gen/services/completor/CompletorServiceClientPb';
import { DMVServiceClient } from '@fivenet/gen/services/dmv/VehiclesServiceClientPb';
import { DocStoreServiceClient } from '@fivenet/gen/services/docstore/DocstoreServiceClientPb';
import { JobsServiceClient } from '@fivenet/gen/services/jobs/JobsServiceClientPb';
import { LivemapperServiceClient } from '@fivenet/gen/services/livemapper/LivemapServiceClientPb';
import { UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { useLoaderStore } from '../store/loader';
import { RpcError, StatusCode } from 'grpc-web';
import { dispatchNotification } from '../components/notification';

export default defineNuxtPlugin(() => {
    return {
        provide: {
            grpc: new GRPCClients(),
        },
    };
});

export class GRPCClients {
    private authInterceptor: AuthInterceptor;
    private grpcClientOptions: { [index: string]: any };

    constructor() {
        this.authInterceptor = new AuthInterceptor();

        // See https://github.com/jrapoport/grpc-web-devtools#grpc-web-interceptor-support
        this.grpcClientOptions = {
            unaryInterceptors: [this.authInterceptor],
            streamInterceptors: [this.authInterceptor],
        };

        //@ts-ignore GRPCWeb Devtools only exist when the user has the extension installed
        const devInterceptors = window.__GRPCWEB_DEVTOOLS__;
        if (devInterceptors) {
            const { devToolsUnaryInterceptor, devToolsStreamInterceptor } = devInterceptors();

            this.grpcClientOptions.unaryInterceptors.push(devToolsUnaryInterceptor);
            this.grpcClientOptions.streamInterceptors.push(devToolsStreamInterceptor);
        }
    }

    // Handle GRPC errors
    async handleRPCError(err: RpcError): Promise<void> {
        switch (err.code) {
            case StatusCode.UNAUTHENTICATED:
                await useAuthStore().clear();

                dispatchNotification({
                    title: 'Please login again',
                    content: 'You are not signed in anymore',
                    type: 'warning',
                });

                // Only update the redirect query param if it isn't set already
                const router = useRouter();
                const redirect = router.currentRoute.value.query.redirect ?? router.currentRoute.value.fullPath;
                await router.push({ name: 'auth-login', query: { redirect: redirect }, replace: true, force: true });
            case StatusCode.PERMISSION_DENIED:
                dispatchNotification({ title: 'Permission denied', content: err.message, type: 'error' });
                break;
            case StatusCode.INTERNAL:
                dispatchNotification({ title: 'Internal server error occured', content: err.message, type: 'error' });
                break;
            case StatusCode.UNAVAILABLE:
                dispatchNotification({
                    title: 'Unable to reach server',
                    content: 'Unable to reach FiveNet server, please check your internet connection.',
                    type: 'error',
                });
                break;
            default:
                dispatchNotification({ title: 'Unknown error occured', content: err.message, type: 'error' });
                break;
        }
    }

    // GRPC Clients ===============================================================
    // Account / Auth - Unauthorized and authorized clients
    private unAuthClient: undefined | AuthServiceClient;
    getUnAuthClient(): AuthServiceClient {
        if (!this.unAuthClient) {
            this.unAuthClient = new AuthServiceClient(config.apiProtoURL, null, null);
        }

        return this.unAuthClient;
    }

    private authClient: undefined | AuthServiceClient;
    getAuthClient(): AuthServiceClient {
        if (!this.authClient) {
            this.authClient = new AuthServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.authClient;
    }

    // Citizens
    private citizenStoreClient: undefined | CitizenStoreServiceClient;
    getCitizenStoreClient(): CitizenStoreServiceClient {
        if (!this.citizenStoreClient) {
            this.citizenStoreClient = new CitizenStoreServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.citizenStoreClient;
    }

    // Completion
    private completorClient: undefined | CompletorServiceClient;
    getCompletorClient(): CompletorServiceClient {
        if (!this.completorClient) {
            this.completorClient = new CompletorServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.completorClient;
    }

    // DMV (Vehicles)
    private dmvClient: undefined | DMVServiceClient;
    getDMVClient(): DMVServiceClient {
        if (!this.dmvClient) {
            this.dmvClient = new DMVServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.dmvClient;
    }

    // Documents
    private docstoreClient: undefined | DocStoreServiceClient;
    getDocStoreClient(): DocStoreServiceClient {
        if (!this.docstoreClient) {
            this.docstoreClient = new DocStoreServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.docstoreClient;
    }

    // Job
    private jobsClient: undefined | JobsServiceClient;
    getJobsClient(): JobsServiceClient {
        if (!this.jobsClient) {
            this.jobsClient = new JobsServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.jobsClient;
    }

    // Livemap
    private livemapperClient: undefined | LivemapperServiceClient;
    getLivemapperClient(): LivemapperServiceClient {
        if (!this.livemapperClient) {
            this.livemapperClient = new LivemapperServiceClient(config.apiProtoURL, null, this.grpcClientOptions);
        }

        return this.livemapperClient;
    }
}

export class AuthInterceptor implements UnaryInterceptor<any, any> {
    private store;

    constructor() {
        this.store = useAuthStore();
    }

    intercept(request: any, invoker: any) {
        if (this.store.$state.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + this.store.$state.accessToken;
        }
        return invoker(request);
    }
}
