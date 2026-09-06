import type { AsyncDataOptions } from 'nuxt/app';
import { toValue, type MaybeRefOrGetter } from 'vue';

type AuthDataScope = keyof ReturnType<typeof useAuth>['keys'];
type KeysOf<T> = Array<T extends T ? (keyof T extends string ? keyof T : never) : never>;

type AuthedAsyncDataHandler<T> = (context: { signal: AbortSignal }) => Promise<T>;

/**
 * Lazy async data scoped to the last route-validated auth context.
 *
 * The handler only receives the abort signal because Nuxt's app instance is
 * not useful to callers and would otherwise be repeated at every call site.
 */
export function useAuthedLazyAsyncData<
    ResT,
    DataT = ResT,
    PickKeys extends KeysOf<DataT> = KeysOf<DataT>,
    DefaultT = undefined,
>(
    scope: AuthDataScope,
    resourceKey: MaybeRefOrGetter<string>,
    handler: AuthedAsyncDataHandler<ResT>,
    options: AsyncDataOptions<ResT, DataT, PickKeys, DefaultT> = {},
) {
    const auth = useAuth();
    const enabled = computed(() => !auth.isQueryTransitioning.value && toValue(options.enabled ?? true));

    return useLazyAsyncData<ResT, unknown, DataT, PickKeys, DefaultT>(
        computed(() => `${auth.keys[scope].value}:${toValue(resourceKey)}`),
        (_nuxtApp, { signal }) => handler({ signal }),
        {
            ...options,
            dedupe: options.dedupe ?? 'cancel',
            enabled,
        },
    );
}
