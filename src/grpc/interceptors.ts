import { UnaryInterceptor, UnaryResponse } from 'grpc-web';
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

export class UnaryLoaderInterceptor implements UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any): Promise<UnaryResponse<any, any>> {
        store.dispatch('loader/show');

        return invoker(request)
            .finally(() => {
                store.dispatch('loader/hide');
            });
    }
}
