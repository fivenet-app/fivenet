import store from './store';
import { RpcError, StatusCode, StreamInterceptor, UnaryInterceptor } from 'grpc-web';
import { _RouteLocationBase } from 'vue-router/auto';
import router from './main';
import config from './config';
import { dispatchNotification } from './components/notification';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { DispatchesServiceClient } from '../gen/dispatches/DispatchesServiceClientPb';
import { DocumentsServiceClient } from '../gen/documents/DocumentsServiceClientPb';
import { JobServiceClient } from '../gen/job/JobServiceClientPb';
import { LivemapServiceClient } from '@arpanet/gen/livemap/LivemapServiceClientPb';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';

class AuthInterceptor implements StreamInterceptor<any, any>, UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        if (store.state.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + store.state.accessToken;
        }
        return invoker(request);
    }
}

export const authInterceptor = new AuthInterceptor();

// Handle GRPC errors
export function handleGRPCError(err: RpcError, route: _RouteLocationBase): boolean {
    if (err.code == StatusCode.UNAUTHENTICATED) {
        store.commit('doLogout');

        router.push({ name: 'login', query: { redirect: route.fullPath } });

        dispatchNotification({ title: 'Please login again!', content: 'You are not signed in anymore', type: 'warning' });
        return true;
    } else if (err.code == StatusCode.PERMISSION_DENIED) {
        dispatchNotification({ title: 'Error!', content: err.message, type: 'error' });
        return true;
    } else {
        dispatchNotification({ title: 'Unknown error occured!', content: err.message, type: 'error' });
    }

    return false;
}

// Account / Auth
// Only the authorized client is handled here
let authClient: AccountServiceClient;
export function getAccountClient(): AccountServiceClient {
    if (!authClient) {
        authClient = new AccountServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return authClient;
}

// Dispatches
let dispatchesClient: DispatchesServiceClient;
export function getDispatchesClient(): DispatchesServiceClient {
    if (!dispatchesClient) {
        dispatchesClient = new DispatchesServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return dispatchesClient;
}

// Documents
let documentsClient: DocumentsServiceClient;
export function getDocumentsClient(): DocumentsServiceClient {
    if (!documentsClient) {
        documentsClient = new DocumentsServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return documentsClient;
}

// Job
let jobClient: JobServiceClient;
export function getJobClient(): JobServiceClient {
    if (!jobClient) {
        jobClient = new JobServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return jobClient;
}

// Livemap
let livemapClient: LivemapServiceClient;
export function getLivemapClient(): LivemapServiceClient {
    if (!livemapClient) {
        livemapClient = new LivemapServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return livemapClient;
}

// Users
let usersClient: UsersServiceClient;
export function getUsersClient(): UsersServiceClient {
    if (!usersClient) {
        usersClient = new UsersServiceClient(config.apiProtoURL, null, {
            unaryInterceptors: [authInterceptor],
            streamInterceptors: [authInterceptor],
        });
    }

    return usersClient;
}
