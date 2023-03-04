import store from './store';
import * as grpcweb from 'grpc-web';
import { _RouteLocationBase } from 'vue-router/auto';
import router from './main';
import { dispatchNotification } from './components/notification';

class AuthInterceptor implements grpcweb.StreamInterceptor<any, any>, grpcweb.UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        if (store.state.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + store.state.accessToken;
        }
        return invoker(request);
    }
}

export const authInterceptor = new AuthInterceptor();

export const clientAuthOptions = {
    unaryInterceptors: [authInterceptor],
    streamInterceptors: [authInterceptor],
} as { [index: string]: any };

// Handle GRPC errors
export function handleGRPCError(err: grpcweb.RpcError, route: _RouteLocationBase): boolean {
    if (err.code == grpcweb.StatusCode.UNAUTHENTICATED) {
        store.commit('doLogout');

        router.push({ name: 'login', query: { redirect: route.fullPath } });

        dispatchNotification({ title: 'Please login again!', content: 'You are not signed in anymore', type: 'warning' });
        return true;
    } else if (err.code == grpcweb.StatusCode.PERMISSION_DENIED) {
        dispatchNotification({ title: 'Error!', content: err.message, type: 'error' });
        return true;
    } else {
        dispatchNotification({ title: 'Unknown error occured!', content: err.message, type: 'error' });
    }

    return false;
}
