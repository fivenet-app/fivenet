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

export class LoadingManager {
    public counted: Ref<number>;

    constructor() {
        this.counted = ref(0);
    }

    async start(): Promise<void> {
        this.counted.value++;
        if (this.counted.value === 1) {
            useNuxtApp().callHook('data:loading:start');
        }
    }

    async finish(): Promise<void> {
        if (this.counted.value > 0) {
            this.counted.value--;
            useNuxtApp().callHook('data:loading:finish');
        }
    }

    async errored(): Promise<void> {
        this.counted.value = 0;
        useNuxtApp().callHook('data:loading:finish_error');
    }
}
