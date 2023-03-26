import { RpcError, StatusCode, UnaryInterceptor, UnaryResponse } from 'grpc-web';
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

        // Add our basic grpc error handler
        return invoker(request)
            .catch(async (err: RpcError) => {
                await handleGRPCError(err);
            })
            .finally(() => {
                store.dispatch('loader/hide');
            });
    }
}

// Handle GRPC errors
export async function handleGRPCError(err: RpcError) {
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
