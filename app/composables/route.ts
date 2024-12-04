import type { UnwrapNestedRefs, WritableComputedRef } from 'vue';

const _useRouteQueryObject = <T extends object>(key: string, def: T): WritableComputedRef<T> => {
    const query = useRouteQuery(key, JSON.stringify(def));

    const value = computed({
        get() {
            try {
                const parsed = JSON.parse(query.value as string) as T;
                return {
                    ...def,
                    ...parsed,
                };
            } catch (_) {
                return def;
            }
        },
        set(val) {
            query.value = JSON.stringify(val);
        },
    });

    return value;
};

export const useRouteQueryObject = createSharedComposable(_useRouteQueryObject);

const _useRouteQueryReactive = <T extends object>(key: string, def: T): UnwrapNestedRefs<T> => {
    const query = useRouteQuery(key, JSON.stringify(def));

    const value = reactiveComputed(() => {
        try {
            const parsed = JSON.parse(query.value as string) as T;
            return {
                ...def,
                ...parsed,
            };
        } catch (_) {
            return def;
        }
    });

    watch(value, () => (query.value = JSON.stringify(value)), { deep: true, flush: 'sync' });

    return value;
};

export const useRouteQueryReactive = createSharedComposable(_useRouteQueryReactive);
