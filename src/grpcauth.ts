import store from './store';
import { RpcError, StatusCode, StreamInterceptor, UnaryInterceptor } from 'grpc-web';
import { _RouteLocationBase } from 'vue-router/auto';
import router from './main';
import { dispatchNotification } from './components/Notification';

class AuthInterceptor implements StreamInterceptor<any, any>, UnaryInterceptor<any, any> {
	intercept(request: any, invoker: any) {
		if (store.state.accessToken) {
			const metadata = request.getMetadata();
			metadata.Authorization = 'Bearer ' + store.state.accessToken;
		}
		return invoker(request);
	}
	handleError(err: RpcError, route: _RouteLocationBase): boolean {
		if (err.code == StatusCode.UNAUTHENTICATED) {
			store.commit('doLogout');

            router.push({ name: 'login', query: { redirect: route.fullPath } });

            dispatchNotification({ title: 'Please login again!', content: 'You are not signed in anymore', type: 'warning' });
			return true;
		} else if (err.code == StatusCode.PERMISSION_DENIED) {
            dispatchNotification({ title: 'Error!', content: err.message, type: 'error' });
			return true;
		}

		return false;
	}
}

const authInterceptor = new AuthInterceptor();

export default authInterceptor;
