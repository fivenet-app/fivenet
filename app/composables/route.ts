const _useRouteQueryObject = <T extends object>(key: string, def: T): Ref<T> => {
    const query = useRouteQuery<string, T>(key, JSON.stringify(def), {
        transform: {
            get(val) {
                try {
                    const parsed = JSON.parse(val) as T;
                    return {
                        ...def,
                        ...parsed,
                    };
                } catch (_) {
                    return def;
                }
            },
            set(val) {
                return JSON.stringify(val);
            },
        },
    });

    return query;
};

export const useRouteQueryObject = _useRouteQueryObject;
