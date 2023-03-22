import config from '../config';
import { AuthServiceClient } from '@arpanet/gen/services/auth/AuthServiceClientPb';
import { CitizenStoreServiceClient } from '@arpanet/gen/services/citizenstore/CitizenstoreServiceClientPb';
import { CompletorServiceClient } from '@arpanet/gen/services/completor/CompletorServiceClientPb';
import { DispatcherServiceClient } from '@arpanet/gen/services/dispatcher/DispatcherServiceClientPb';
import { DocStoreServiceClient } from '@arpanet/gen/services/docstore/DocstoreServiceClientPb';
import { JobsServiceClient } from '@arpanet/gen/services/jobs/JobsServiceClientPb';
import { LivemapperServiceClient } from '@arpanet/gen/services/livemapper/LivemapServiceClientPb';
import { DMVServiceClient } from '@arpanet/gen/services/dmv/VehiclesServiceClientPb';
import { AuthInterceptor, UnaryErrorHandlerInterceptor } from './interceptors';

const authInterceptor = new AuthInterceptor();

// See https://github.com/jrapoport/grpc-web-devtools#grpc-web-interceptor-support
export const grpcClientOptions = {
    unaryInterceptors: [authInterceptor, new UnaryErrorHandlerInterceptor()],
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
// Account / Auth - Only the authorized client is kept here
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

// Dispatches
let dispatcherClient: DispatcherServiceClient;
export function getDispatcherClient(): DispatcherServiceClient {
    if (!dispatcherClient) {
        dispatcherClient = new DispatcherServiceClient(config.apiProtoURL, null, grpcClientOptions);
    }

    return dispatcherClient;
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
