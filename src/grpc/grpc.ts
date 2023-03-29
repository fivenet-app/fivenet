import config from '../config';
import { AuthServiceClient } from '@arpanet/gen/services/auth/AuthServiceClientPb';
import { CitizenStoreServiceClient } from '@arpanet/gen/services/citizenstore/CitizenstoreServiceClientPb';
import { CompletorServiceClient } from '@arpanet/gen/services/completor/CompletorServiceClientPb';
import { DMVServiceClient } from '@arpanet/gen/services/dmv/VehiclesServiceClientPb';
import { DocStoreServiceClient } from '@arpanet/gen/services/docstore/DocstoreServiceClientPb';
import { JobsServiceClient } from '@arpanet/gen/services/jobs/JobsServiceClientPb';
import { LivemapperServiceClient } from '@arpanet/gen/services/livemapper/LivemapServiceClientPb';
import { AuthInterceptor, UnaryLoaderInterceptor } from './interceptors';
import { RpcError, StatusCode } from 'grpc-web';
import { dispatchNotification } from '../components/notification';
import { store } from '../store/store';
import { router } from '../router';

// Handle GRPC errors
export async function handleRPCError(err: RpcError) {
    switch (err.code) {
        case StatusCode.UNAUTHENTICATED:
            await store.dispatch('auth/doLogout');

            dispatchNotification({ title: 'Please login again', content: 'You are not signed in anymore', type: 'warning' });

            // Only update the redirect query param if it isn't set already
            const redirect = router.currentRoute.value.query.redirect ?? router.currentRoute.value.fullPath;
            await router.push({ path: '/login', query: { redirect: redirect }, replace: true, force: true });
            break;
        case StatusCode.PERMISSION_DENIED:
            dispatchNotification({ title: 'Permission denied', content: err.message, type: 'error' });
            break;
        case StatusCode.INTERNAL:
            dispatchNotification({ title: 'Internal server error occured', content: err.message, type: 'error' });
            break;
        case StatusCode.UNAVAILABLE:
            dispatchNotification({
                title: 'Unable to reach server',
                content: 'Unable to reach aRPaNet server, please check your internet connection.',
                type: 'error',
            });
            break;
        default:
            dispatchNotification({ title: 'Unknown error occured', content: err.message, type: 'error' });
            break;
    }
}

const authInterceptor = new AuthInterceptor();
export const unaryErrorHandlerInterceptor = new UnaryLoaderInterceptor();

// See https://github.com/jrapoport/grpc-web-devtools#grpc-web-interceptor-support
export const grpcClientOptions = {
    unaryInterceptors: [authInterceptor, unaryErrorHandlerInterceptor],
    streamInterceptors: [authInterceptor],
} as { [index: string]: any };

//@ts-ignore GRPCWeb Devtools only exist when the user has the extension installed
const devInterceptors = window.__GRPCWEB_DEVTOOLS__;
if (devInterceptors) {
    const { devToolsUnaryInterceptor, devToolsStreamInterceptor } = devInterceptors();

    grpcClientOptions.unaryInterceptors.push(devToolsUnaryInterceptor);
    grpcClientOptions.streamInterceptors.push(devToolsStreamInterceptor);
}

// GRPC Clients ===============================================================
// Account / Auth - Unauthorized and authorized clients
let unAuthClient: AuthServiceClient;
export function getUnAuthClient(): AuthServiceClient {
    if (!unAuthClient) {
        unAuthClient = new AuthServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [unaryErrorHandlerInterceptor],
        });
    }

    return unAuthClient;
}

let authClient: AuthServiceClient;
export function getAuthClient(): AuthServiceClient {
    if (!authClient) {
        authClient = new AuthServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return authClient;
}

// Citizens
let citizenStoreClient: CitizenStoreServiceClient;
export function getCitizenStoreClient(): CitizenStoreServiceClient {
    if (!citizenStoreClient) {
        citizenStoreClient = new CitizenStoreServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return citizenStoreClient;
}

// Completion
let completorClient: CompletorServiceClient;
export function getCompletorClient(): CompletorServiceClient {
    if (!completorClient) {
        completorClient = new CompletorServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return completorClient;
}

// DMV (Vehicles)
let dmvClient: DMVServiceClient;
export function getDMVClient(): DMVServiceClient {
    if (!dmvClient) {
        dmvClient = new DMVServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return dmvClient;
}

// Documents
let docstoreClient: DocStoreServiceClient;
export function getDocStoreClient(): DocStoreServiceClient {
    if (!docstoreClient) {
        docstoreClient = new DocStoreServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return docstoreClient;
}

// Job
let jobsClient: JobsServiceClient;
export function getJobsClient(): JobsServiceClient {
    if (!jobsClient) {
        jobsClient = new JobsServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return jobsClient;
}

// Livemap
let livemapperClient: LivemapperServiceClient;
export function getLivemapperClient(): LivemapperServiceClient {
    if (!livemapperClient) {
        livemapperClient = new LivemapperServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return livemapperClient;
}
