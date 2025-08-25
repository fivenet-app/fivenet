/* eslint-disable @typescript-eslint/no-explicit-any */
import { useRoute, useRouter } from 'vue-router';
import { type z, type ZodDefault, ZodFirstPartyTypeKind, type ZodObject, type ZodRawShape, type ZodTypeAny } from 'zod';
import { useSearchesStore } from '~/stores/searches';

// helper to unwrap a default‐wrapped schema
function getCoreSchema(s: ZodTypeAny): ZodTypeAny {
    if (s._def.typeName === ZodFirstPartyTypeKind.ZodDefault) {
        return (s as ZodDefault<any>)._def.innerType;
    }
    return s;
}

// format a Date for ?date.start=YYYY-MM-DD
function fmtDate(d: Date) {
    return d.toISOString().slice(0, 10);
}

/**
 * key:    unique identifier for this form (only for store lookup)
 * schema: your ZodObject describing the form state (with any .default()s you want)
 */
export function useSearchForm<T extends ZodRawShape, S extends ZodObject<T>>(key: string, schema: S) {
    type State = z.infer<S>;
    const route = useRoute();
    const router = useRouter();
    const store = useSearchesStore();

    // Step 1: compute your baseline defaults
    const defaults: State = schema.parse({});

    // Step 2: load any existing state from the central store
    //           or fall back to defaults if none yet
    const base = store.getSearch<State>(key) ?? defaults;

    // Step 3: build a "raw" override object from URL params
    const raw = {} as Partial<State>;
    for (const field of Object.keys(schema.shape) as Array<keyof State>) {
        const core = getCoreSchema(schema.shape[field] as ZodTypeAny);
        const type = core._def.typeName as ZodFirstPartyTypeKind;
        const param = route.query[field as string];

        // Skip empty string
        if (param === '') continue;

        // booleans
        if (type === ZodFirstPartyTypeKind.ZodBoolean) {
            if (param === 'true') raw[field] = true as any;
            if (param === 'false') raw[field] = false as any;
            continue;
        }

        // Dates
        if (type === ZodFirstPartyTypeKind.ZodDate && typeof param === 'string') {
            raw[field] = new Date(param) as any;
            continue;
        }

        // Everything else (arrays, objects, nested, primitives): try JSON.parse
        if (typeof param === 'string') {
            try {
                const parsed = JSON.parse(param);
                // Skip empty arrays/objects?
                if (typeof parsed === 'number' && isNaN(parsed)) {
                    continue;
                }
                raw[field] = parsed as any;
            } catch {
                // Fallback to string
                raw[field] = param as any;
            }
        }
    }

    // Step 4: Merge defaults ← store ← raw, then validate
    const merged: State = schema.parse({ ...defaults, ...base, ...raw });

    // Step 5: Put or patch the reactive in the central store
    let state = store.getSearch<State>(key);
    if (state) {
        Object.assign(state, merged);
    } else {
        state = reactive<State>(merged) as State;
        store.setSearch(key, state);
    }

    // Step 6: Watch it → sync back to URL & store
    watch(
        state,
        (s) => {
            const q: Record<string, string> = {};
            if (route.query.tab) {
                // TODO this is a temporary fix
                q.tab = String(route.query.tab);
            }

            for (const [k, v] of Object.entries(s) as [keyof State, any][]) {
                // Skip null/undefined
                if (v === undefined || v === null) continue;
                // Skip empty string
                if (typeof v === 'string' && v === '') continue;
                // Skip NaN
                if (typeof v === 'number' && isNaN(v)) continue;

                const core = getCoreSchema(schema.shape[k] as ZodTypeAny);
                const type = core._def.typeName as ZodFirstPartyTypeKind;

                if (type === ZodFirstPartyTypeKind.ZodDate) {
                    q[k as string] = fmtDate(v as Date);
                } else if (type === ZodFirstPartyTypeKind.ZodBoolean) {
                    q[k as string] = String(v);
                } else if (typeof v === 'string' || typeof v === 'number') {
                    q[k as string] = String(v);
                } else {
                    // Everything else (arrays & objects)
                    q[k as string] = JSON.stringify(v);
                }
            }

            router.replace({ query: q });
            // update the central store
            store.setSearch(key, state);
        },
        { deep: true, immediate: true },
    );

    return state;
}
