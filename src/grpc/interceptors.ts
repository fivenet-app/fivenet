import { ClientReadableStream, RpcError, StatusCode, StreamInterceptor, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { dispatchNotification } from '../components/notification';
import { router } from '../router';
import { store } from '../store/store';

export class AuthInterceptor implements UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        if (store.state.auth?.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + store.state.auth?.accessToken;
        }
        return invoker(request);
    }
}

export class UnaryErrorHandlerInterceptor implements UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any): Promise<UnaryResponse<any, any>> {
        store.dispatch('loader/show');

        const res = invoker(request);
        // Add our basic grpc error handler
        res.catch((err: RpcError) => {
            handleGRPCError(err);
        }).finally(() => {
            store.dispatch('loader/hide');
        });

        return res;
    }
}

export class StreamErrorHandlerInterceptor implements StreamInterceptor<any, any> {
    intercept(request: any, invoker: any): ClientReadableStream<any> {
        const res = invoker(request);
        // Add our basic grpc error handler
        res.on('error', (err: RpcError) => {
            handleGRPCError(err);
        }).finally(() => {
            store.dispatch('loader/hide');
        });

        return res;
    }
}

// Handle GRPC errors
function handleGRPCError(err: RpcError) {
    if (err.code == StatusCode.UNAUTHENTICATED) {
        store.dispatch('auth/doLogout');

        dispatchNotification({ title: 'Please login again', content: 'You are not signed in anymore', type: 'warning' });

        router.push({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } });
    } else if (err.code == StatusCode.PERMISSION_DENIED) {
        dispatchNotification({ title: 'Permission denied', content: err.message, type: 'error' });
    } else if (err.code == StatusCode.INTERNAL) {
        dispatchNotification({ title: 'Internal server error occured', content: err.message, type: 'error' });
    }
}
