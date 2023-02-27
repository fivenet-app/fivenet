import { store } from './store';

class AuthInterceptor {
	intercept(request: any, invoker: any) {
		if (store.state.accessToken) {
			const metadata = request.getMetadata();
			metadata.Authorization = 'Bearer ' + store.state.accessToken;
		}
		return invoker(request);
	}
}

const authInterceptor = new AuthInterceptor();

export default authInterceptor;
