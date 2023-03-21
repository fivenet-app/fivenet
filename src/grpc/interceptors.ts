import { StreamInterceptor, UnaryInterceptor } from 'grpc-web';
import { store } from '../store/store';

class AuthInterceptor implements UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        if (store.state.auth?.accessToken) {
            const metadata = request.getMetadata();
            metadata.Authorization = 'Bearer ' + store.state.auth?.accessToken;
        }
        return invoker(request);
    }
}

export const authInterceptor = new AuthInterceptor();

class UnaryErrorHandlerInterceptor implements UnaryInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        // TODO
        return invoker(request);
    }
}

export const unaryErrorHandlerInterceptor = new UnaryErrorHandlerInterceptor();

class StreamErrorHandlerInterceptor implements StreamInterceptor<any, any> {
    intercept(request: any, invoker: any) {
        // TODO
        return invoker(request);
    }
}

export const streamErrorHandlerInterceptor = new StreamErrorHandlerInterceptor();
