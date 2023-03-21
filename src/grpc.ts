import { store } from './store/store';
import { RpcError, StatusCode, StreamInterceptor, UnaryInterceptor } from 'grpc-web';
import { _RouteLocationBase } from 'vue-router/auto';
import { router } from './router';
import config from './config';
import { dispatchNotification } from './components/notification';
import { AuthServiceClient } from '@arpanet/gen/services/auth/AuthServiceClientPb';
import { CitizenStoreServiceClient } from '@arpanet/gen/services/citizenstore/CitizenstoreServiceClientPb';
import { CompletorServiceClient } from '@arpanet/gen/services/completor/CompletorServiceClientPb';
import { DispatcherServiceClient } from '@arpanet/gen/services/dispatcher/DispatcherServiceClientPb';
import { DocStoreServiceClient } from '@arpanet/gen/services/docstore/DocstoreServiceClientPb';
import { JobsServiceClient } from '@arpanet/gen/services/jobs/JobsServiceClientPb';
import { LivemapperServiceClient } from '@arpanet/gen/services/livemapper/LivemapServiceClientPb';
import { DMVServiceClient } from '@arpanet/gen/services/dmv/VehiclesServiceClientPb';

class AuthInterceptor implements StreamInterceptor<any, any>, UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        if (store.state.auth?.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + store.state.auth?.accessToken;
        }
        return invoker(request);
    }
}

export const authInterceptor = new AuthInterceptor();

class LoaderInterceptor implements StreamInterceptor<any, any>, UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        console.log("SHOWING LOADER");
        store.dispatch('loader/show');
        return invoker(request);
    }
}

const loaderInterceptor = new LoaderInterceptor();

// See https://github.com/jrapoport/grpc-web-devtools#grpc-web-interceptor-support
export const clientAuthOptions = {
    unaryInterceptors: [authInterceptor, loaderInterceptor],
    streamInterceptors: [authInterceptor, loaderInterceptor],
} as { [index: string]: any };

//@ts-ignore GRPCWeb Devtools only exist when the user has the extension installed
const devInterceptors = window.__GRPCWEB_DEVTOOLS__;
if (devInterceptors) {
    const { devToolsUnaryInterceptor, devToolsStreamInterceptor } = devInterceptors();

    clientAuthOptions.unaryInterceptors.push(devToolsUnaryInterceptor);
    clientAuthOptions.streamInterceptors.push(devToolsStreamInterceptor);
}

// Handle GRPC errors
export function handleGRPCError(err: RpcError, route: _RouteLocationBase): boolean {
    if (err.code == StatusCode.UNAUTHENTICATED) {
        store.dispatch('auth/doLogout');

        dispatchNotification({ title: 'Please login again!', content: 'You are not signed in anymore', type: 'warning' });

        router.push({ path: '/login', query: { redirect: route.fullPath } });
        return true;
    } else if (err.code == StatusCode.PERMISSION_DENIED) {
        dispatchNotification({ title: 'Error!', content: err.message, type: 'error' });
        return true;
    } else {
        dispatchNotification({ title: 'Unknown error occured!', content: err.message, type: 'error' });
    }

    return false;
}

// GRPC Clients ===============================================================
// Account / Auth - Only the authorized client is kept here
let authClient: AuthServiceClient;
export function getAuthClient(): AuthServiceClient {
    if (!authClient) {
        authClient = new AuthServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return authClient;
}

// Citizens
let citizenStoreClient: CitizenStoreServiceClient;
export function getCitizenStoreClient(): CitizenStoreServiceClient {
    if (!citizenStoreClient) {
        citizenStoreClient = new CitizenStoreServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return citizenStoreClient;
}

// Completion
let completorClient: CompletorServiceClient;
export function getCompletorClient(): CompletorServiceClient {
    if (!completorClient) {
        completorClient = new CompletorServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return completorClient;
}

// Dispatches
let dispatcherClient: DispatcherServiceClient;
export function getDispatcherClient(): DispatcherServiceClient {
    if (!dispatcherClient) {
        dispatcherClient = new DispatcherServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return dispatcherClient;
}

// DMV (Vehicles)
let dmvClient: DMVServiceClient;
export function getDMVClient(): DMVServiceClient {
    if (!dmvClient) {
        dmvClient = new DMVServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return dmvClient;
}

// Documents
let docstoreClient: DocStoreServiceClient;
export function getDocStoreClient(): DocStoreServiceClient {
    if (!docstoreClient) {
        docstoreClient = new DocStoreServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return docstoreClient;
}

// Job
let jobsClient: JobsServiceClient;
export function getJobsClient(): JobsServiceClient {
    if (!jobsClient) {
        jobsClient = new JobsServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return jobsClient;
}

// Livemap
let livemapperClient: LivemapperServiceClient;
export function getLivemapperClient(): LivemapperServiceClient {
    if (!livemapperClient) {
        livemapperClient = new LivemapperServiceClient(config.apiProtoURL, null, clientAuthOptions);
    }

    return livemapperClient;
}
