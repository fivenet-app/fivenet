import { watchDebounced } from '@vueuse/shared';

type MaybePromise<T> = T | Promise<T>;

type RefreshFn = () => MaybePromise<unknown>;

export function useDebouncedRefresh(
    source: Parameters<typeof watchDebounced>[0],
    refresh: RefreshFn,
    options: {
        debounce?: number;
        maxWait?: number;
    } = {},
): void {
    const { debounce = 175, maxWait = 1150 } = options;

    watchDebounced(
        source,
        async () => {
            await refresh();
        },
        {
            debounce,
            maxWait,
            deep: true,
        },
    );
}
