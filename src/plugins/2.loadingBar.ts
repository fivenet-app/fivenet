import { NextUnaryFn, MethodInfo, RpcOptions, RpcInterceptor, UnaryCall } from '@protobuf-ts/runtime-rpc/build/types';

export default defineNuxtPlugin((nuxtApp) => {
    const lm = new LoadingManager();

    nuxtApp.hook('page:start', async () => lm.start());

    nuxtApp.hook('page:finish', async () => lm.finish());

    nuxtApp.hook('app:error', async (_: any) => lm.errored());

    return {
        provide: {
            loading: lm,
        },
    };
});

export class LoadingManager implements RpcInterceptor {
    private counted: number;

    constructor() {
        this.counted = 1;
    }

    async start(): Promise<void> {
        this.counted++;
        if (this.counted === 1) {
            //@ts-ignore TODO need to add types in the future
            useNuxtApp().callHook('data:loading:start');
        }
    }

    async finish(): Promise<void> {
        if (this.counted > 0) {
            this.counted--;
            //@ts-ignore TODO need to add types in the future
            useNuxtApp().callHook('data:loading:finish');
        }
    }

    async errored(): Promise<void> {
        this.counted = 0;
        //@ts-ignore TODO need to add types in the future
        useNuxtApp().callHook('data:loading:finish_error');
    }

    // GRPC unary interceptor
    interceptUnary(next: NextUnaryFn, method: MethodInfo, input: object, options: RpcOptions): UnaryCall {
        this.start();
        const ret = next(method, input, options);
        this.finish();
        return ret;
    }
}
