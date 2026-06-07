import { watchDebounced } from '@vueuse/shared';

type MaybePromise<T> = T | Promise<T>;

type ValidateFn = () => MaybePromise<unknown>;
type RefreshFn = () => MaybePromise<unknown>;

export function useFormValidatedDebouncedRefresh(
    source: Parameters<typeof watchDebounced>[0],
    validate: ValidateFn,
    refresh: RefreshFn,
    options: {
        debounce?: number;
        maxWait?: number;
    } = {},
): void {
    const { debounce = 200, maxWait = 1250 } = options;

    watchDebounced(
        source,
        async () => {
            const valid = await validate();
            if (!valid) return;

            await refresh();
        },
        {
            debounce,
            maxWait,
            deep: true,
        },
    );
}
