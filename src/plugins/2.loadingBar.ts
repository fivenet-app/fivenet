export default defineNuxtPlugin((nuxtApp) => {
    const lm = new LoadingManager();

    nuxtApp.hook('page:start', async () => lm.start());

    nuxtApp.hook('page:finish', async () => lm.finish());

    nuxtApp.hook('app:error', async (err: any) => lm.errored());

    return {
        provide: {
            loading: lm,
        },
    };
});

export class LoadingManager {
    private counted: number;

    constructor() {
        this.counted = 1;
    }

    async start(): Promise<void> {
        console.log('data:loading:start:', this.counted);
        this.counted++;
        if (this.counted == 1) {
            useNuxtApp().callHook('data:loading:start');
        }
    }

    async finish(): Promise<void> {
        console.log('data:loading:finish:', this.counted);
        if (this.counted > 0) {
            this.counted--;
            useNuxtApp().callHook('data:loading:finish');
        }
    }

    async errored(): Promise<void> {
        console.log('data:loading:finish_error:', this.counted);
        this.counted = 0;
        useNuxtApp().callHook('data:loading:finish_error');
    }

    // GRPC interceptor
    intercept(request: any, invoker: any) {
        this.start();
        const ret = invoker(request);
        this.finish();
        return ret;
    }
}
