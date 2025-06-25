// Extended Vue ↔ Yjs array & object helpers – **complete file**
// ---------------------------------------------------------------------------
//  This module complements `useYPrimitives.ts` by adding higher‑level helpers
//  that bind Vue arrays / nested objects to Yjs containers.  All helpers use
//  the unified fallback‑to‑local initialisation controlled through
//  `YjsSyncOptions` (authoritative client + provider with `synced` event).
// ---------------------------------------------------------------------------

import { getCurrentInstance, onUnmounted, ref, toRef, watch, type Ref } from 'vue';
import * as Y from 'yjs';
import { LOCAL_ORIGIN, useYBoolean, useYNumber, useYText, type YjsSyncOptions } from './useYPrimitives';

// ---------------------------------------------------------------------------
//  Primitive helpers & types
// ---------------------------------------------------------------------------

/** Only string, boolean or number are supported as "primitive" Y.Map values. */
export type Primitive = string | boolean | number;

/** A flat record or nested structure that can be mirrored into Yjs. */
export interface YStateMap {
    [key: string]: Primitive | YStateMap | Array<Primitive | YStateMap>;
}

/** Union for arrays that may contain primitives or nested Y.Map objects. */
export type MaybeMap = Primitive | Y.Map<unknown>;

// ---------------------------------------------------------------------------
//  useYArray  ⇄  Y.Array ⇄ Vue Ref<T[]>
// ---------------------------------------------------------------------------

/**
 * Bind a Vue `Ref<T[]>` to a Yjs `Y.Array` that contains primitives or
 * objects.  Deep objects are handled recursively via `useYSyncStructure`.
 *
 * Fallback logic:
 *   • If remote already contains items → remote wins.
 *   • Otherwise, if `opts.authoritative === true` → local array is serialised
 *     and pushed to the Yjs array.
 */
export function useYArray<T extends object | Primitive>(
    yArray: Y.Array<MaybeMap>,
    itemsRef?: Ref<T[]>,
    opts: YjsSyncOptions = {},
): Ref<T[]> {
    const { provider } = opts;
    const items = (itemsRef ?? ref([])) as Ref<T[]>;
    let remoteApplying = false;

    const isPrimitiveMode = () => {
        const first = yArray.length ? yArray.get(0) : undefined;
        return typeof first === 'string' || typeof first === 'number' || typeof first === 'boolean';
    };

    const syncFromY = () => {
        remoteApplying = true;
        if (isPrimitiveMode()) {
            items.value = yArray.toArray() as T[];
        } else {
            while (items.value.length > yArray.length) items.value.pop();
            while (items.value.length < yArray.length) items.value.push({} as T);
            yArray.forEach((m, i) => {
                if (m instanceof Y.Map) {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    useYSyncStructure(m, items.value[i]! as any);
                }
            });
        }
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocalToY = () => {
        yArray.doc?.transact(() => {
            yArray.delete(0, yArray.length);
            if (isPrimitiveMode()) {
                yArray.insert(0, items.value as unknown as MaybeMap[]);
            } else {
                const maps = (items.value as unknown as object[]).map((o) => {
                    const m = new Y.Map(Object.entries(o));
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    useYSyncStructure(m, o as any);
                    return m;
                });
                yArray.insert(0, maps as MaybeMap[]);
            }
        }, LOCAL_ORIGIN);
    };

    const onSync = (s: boolean) => s && init();

    const init = () => {
        remoteApplying = true;

        if (yArray.length > 0) {
            syncFromY();
        } else if (provider && provider.isAuthoritative) {
            writeLocalToY();
        }

        nextTick(() => {
            remoteApplying = false;
        });

        yArray.observe(handleYArray);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                yArray.unobserve(handleYArray);
            });

        watch(
            items,
            () => {
                if (remoteApplying) return;
                writeLocalToY();
            },
            { deep: true, flush: 'post' },
        );
    };

    if (provider) {
        provider.on('sync', onSync);
    } else {
        init();
    }

    const handleYArray = (_evt: Y.YArrayEvent<MaybeMap>, tr: Y.Transaction) => {
        if (tr.origin === LOCAL_ORIGIN) return;
        syncFromY();
    };

    return items;
}

// ---------------------------------------------------------------------------
//  useYArrayFiltered  ⇄  object array with key filtering
// ---------------------------------------------------------------------------

export type OptsKeyFilter = { omit?: string[]; only?: string[] };

export function useYArrayFiltered<T extends object>(
    yArr: Y.Array<unknown>,
    items: Ref<T[]>,
    filter: OptsKeyFilter,
    opts: YjsSyncOptions = {},
) {
    const { provider } = opts;
    let remoteApplying = false;

    const serialize = (o: T): Record<string, unknown> => {
        if (filter.only) return Object.fromEntries(Object.entries(o).filter(([k]) => filter.only!.includes(k)));
        if (filter.omit) return Object.fromEntries(Object.entries(o).filter(([k]) => !filter.omit!.includes(k)));
        return o as Record<string, unknown>;
    };

    const hydrate = (target: T, map: Y.Map<unknown>) => {
        Object.entries(map.toJSON() as Record<string, unknown>).forEach(([k, v]) => {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            (target as any)[k] = v;
        });
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        useYSyncStructure(map, target as any);
    };

    const syncFromY = () => {
        remoteApplying = true;
        while (items.value.length > yArr.length) items.value.pop();
        while (items.value.length < yArr.length) items.value.push({} as T);
        (yArr as Y.Array<Y.Map<unknown>>).forEach((raw, i) => {
            if (raw instanceof Y.Map) hydrate(items.value[i]!, raw);
        });
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocalToY = () => {
        yArr.doc?.transact(() => {
            yArr.delete(0, yArr.length);
            const maps = items.value.map((o) => new Y.Map(Object.entries(serialize(o))));
            yArr.insert(0, maps);
        }, LOCAL_ORIGIN);
    };

    const handleYArr = (_evt: Y.YArrayEvent<Y.Map<unknown>>, tr: Y.Transaction) => {
        if (tr.origin === LOCAL_ORIGIN) return;
        syncFromY();
    };

    const onSync = (s: boolean) => s && init();

    const init = () => {
        remoteApplying = true;

        if (provider && provider.isAuthoritative) {
            writeLocalToY();
        } else {
            syncFromY();
        }

        nextTick(() => {
            remoteApplying = false;
        });

        (yArr as Y.Array<Y.Map<unknown>>).observe(handleYArr);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                (yArr as Y.Array<Y.Map<unknown>>).unobserve(handleYArr);
            });

        watch(
            items,
            () => {
                if (remoteApplying) return;
                writeLocalToY();
            },
            { deep: true, flush: 'post' },
        );
    };

    if (provider) {
        provider.on('sync', onSync);
    } else {
        init();
    }
}

// ---------------------------------------------------------------------------
//  useYMap  ⇄  flat object ⇄ Y.Map
// ---------------------------------------------------------------------------

/**
 * Binds a subset of fields in a Vue reactive object to a Yjs Y.Map<Primitive>, synchronizing only primitive values.
 *
 * @param yMap    - The Y.Map<Primitive> instance to bind to.
 * @param state   - A Vue‐reactive object whose primitive fields should be synced.
 * @param filter  - Optional {@link OptsKeyFilter} specifying which keys to include or exclude.
 *                  If neither `only` nor `omit` is provided, all keys of `state` are synced.
 * @param opts    - {@link YjsSyncOptions} controlling authoritative seeding and provider handling.
 *                  `opts.authoritative` (boolean) means "seed this map on empty."
 *                  `opts.provider` should be a Yjs provider that emits a `"sync"` event.
 * @returns       The same `state` object, now kept in sync with `yMap`.
 */
/**
 * Bind a plain object (or Ref to an object) whose values are primitives to a Yjs Y.Map.
 * Synchronizes only the keys allowed by the filter.
 *
 * @typeParam T - The shape of the local object.
 * @param yMap        - The Y.Map instance to bind.
 * @param stateOrRef  - The local object or a Ref to it. Its primitive fields will stay in sync.
 * @param filter      - Optional keys filter (whitelist or blacklist).
 * @param opts        - Yjs synchronization options (authoritative seed + provider).
 * @returns The unwrapped state object, kept in sync with Yjs.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function useYMap<T extends Record<string, any>>(
    yMap: Y.Map<Primitive>,
    stateOrRef: T | Ref<T>,
    filter?: OptsKeyFilter,
    opts: YjsSyncOptions = {},
): T {
    const { provider } = opts;

    // Unwrap a Ref<T> or use the raw object
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const isStateRef = (val: any): val is Ref<T> => val && typeof val === 'object' && 'value' in val;
    const state: T = isStateRef(stateOrRef) ? stateOrRef.value : stateOrRef;

    // Compute which keys to synchronize
    const allKeys = Object.keys(state) as (keyof T)[];
    const keys: (keyof T)[] = filter?.only
        ? (filter.only as (keyof T)[])
        : filter?.omit
          ? allKeys.filter((k) => !filter.omit!.includes(k as string))
          : allKeys;

    const stops: Array<() => void> = [];
    let remoteApplying = false;
    let observerAttached = false;

    /** Pull remote Yjs values into the local state. */
    const pullRemote = () => {
        remoteApplying = true;
        keys.forEach((k) => {
            const v = yMap.get(k as string);
            if (v !== undefined) {
                state[k] = v as T[typeof k];
            }
        });
        nextTick(() => {
            remoteApplying = false;
        });
    };

    /** Attach a Yjs observer to handle remote-to-local updates. */
    const attachObserver = () => {
        if (observerAttached) return;
        const handle = (_evt: Y.YMapEvent<Primitive>, tr: Y.Transaction) => {
            if (tr.origin === LOCAL_ORIGIN) return;
            pullRemote();
        };
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        yMap.observe(handle as any);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                yMap.unobserve(handle as any);
            });
        observerAttached = true;
    };

    const onSync = (s: boolean) => s && init();

    // Initialize: pull existing remote data or seed if authoritative, then set up watchers.
    const init = () => {
        const hasRemote = keys.some((k) => yMap.get(k as string) !== undefined);
        if (hasRemote) {
            pullRemote();
        } else if (provider && provider.isAuthoritative) {
            yMap.doc?.transact(() => {
                keys.forEach((k) => yMap.set(k as string, state[k]));
            }, LOCAL_ORIGIN);
        }

        nextTick(() => {
            attachObserver();
            // Local-to-remote: watch each key and write to Yjs when it changes
            keys.forEach((k) => {
                const stop = watch(
                    () => state[k],
                    (newVal) => {
                        if (remoteApplying) return;
                        yMap.doc?.transact(() => yMap.set(k as string, newVal), LOCAL_ORIGIN);
                    },
                    { flush: 'post' },
                );
                stops.push(stop);
            });
        });
    };

    if (provider) {
        provider.on('sync', onSync);
    } else {
        init();
    }

    return state;
}

// ---------------------------------------------------------------------------
//  Internal utilities: container creators & recursive binder
// ---------------------------------------------------------------------------

const isDoc = (c: Y.Doc | Y.Map<unknown>): c is Y.Doc => (c as Y.Doc).clientID !== undefined;

const ensureText = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Text => {
    // Always return a Y.Text. If a primitive was stored, upgrade in‑place.
    const replace = (txt: Y.Text) => {
        if (isDoc(c))
            c.getMap('root')?.set?.(key, txt); // doc root edge‑case
        else c.set(key, txt);
        return txt;
    };
    if (isDoc(c)) return c.getText(key);
    const cur = c.get(key);
    if (cur instanceof Y.Text) return cur;
    const t = new Y.Text();
    if (typeof cur === 'string' && cur.length) t.insert(0, cur);
    return replace(t);
};

const ensureArray = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Array<unknown> => {
    if (isDoc(c)) return c.getArray(key);
    const cur = c.get(key);
    if (cur instanceof Y.Array) return cur;
    const arr = new Y.Array();
    c.set(key, arr);
    return arr;
};

const ensureMap = (c: Y.Doc | Y.Map<unknown>, key: string): Y.Map<unknown> => {
    if (isDoc(c)) return c.getMap(key);
    const cur = c.get(key);
    if (cur instanceof Y.Map) return cur;
    const m = new Y.Map();
    c.set(key, m);
    return m;
};

export function useYSyncStructure<T extends YStateMap>(
    ycontainer: Y.Doc | Y.Map<unknown>,
    state: T,
    fields?: Array<keyof T>,
    opts: YjsSyncOptions = {},
): void {
    const keys = (fields ?? Object.keys(state)) as Array<keyof T>;
    keys.forEach((k) => {
        const name = k as string;
        const val = state[k];

        if (typeof val === 'string') {
            useYText(ensureText(ycontainer, name), toRef(state, name) as Ref<string>, opts);
            return;
        }
        if (typeof val === 'boolean') {
            const target = isDoc(ycontainer) ? ensureMap(ycontainer, 'primitives') : (ycontainer as Y.Map<unknown>);
            useYBoolean(target, name, toRef(state, name) as Ref<boolean>, opts);
            return;
        }
        if (typeof val === 'number') {
            const target = isDoc(ycontainer) ? ensureMap(ycontainer, 'primitives') : (ycontainer as Y.Map<unknown>);
            useYNumber(target, name, toRef(state, name) as Ref<number>, opts);
            return;
        }

        if (Array.isArray(val)) {
            const yarr = ensureArray(ycontainer, name);
            const isObjs = val.every((v) => v && typeof v === 'object' && !Array.isArray(v));
            if (isObjs) {
                useYArrayOfObjects(yarr as Y.Array<Y.Map<unknown>>, toRef(state, name) as Ref<YStateMap[]>, opts);
            } else {
                useYArray<Primitive>(yarr as Y.Array<MaybeMap>, toRef(state, name) as Ref<Primitive[]>, opts);
            }
            return;
        }

        if (val && typeof val === 'object') {
            const nested = ensureMap(ycontainer, name);
            const isFlat = Object.values(val).every((v) => ['string', 'number', 'boolean'].includes(typeof v));
            if (isFlat) {
                useYMap<Record<string, Primitive>>(
                    nested as Y.Map<Primitive>,
                    val as Record<string, Primitive>,
                    undefined,
                    opts,
                );
            } else {
                useYSyncStructure(nested, val as YStateMap);
            }
        }
    });
}

function useYArrayOfObjects(yarr: Y.Array<Y.Map<unknown>>, list: Ref<YStateMap[]>, opts: YjsSyncOptions = {}) {
    const { provider } = opts;
    let remoteApplying = false;

    const syncFromY = () => {
        remoteApplying = true;
        while (list.value.length > yarr.length) list.value.pop();
        while (list.value.length < yarr.length) list.value.push({});
        list.value.forEach((obj, i) => useYSyncStructure(yarr.get(i)!, obj, undefined, opts));
        nextTick(() => {
            remoteApplying = false;
        });
    };

    const writeLocal = () => {
        yarr.doc?.transact(() => {
            yarr.delete(0, yarr.length);
            const maps = list.value.map((o) => {
                const m = new Y.Map(Object.entries(o));
                useYSyncStructure(m, o, undefined, opts);
                return m;
            });
            yarr.insert(0, maps);
        }, LOCAL_ORIGIN);
    };

    const onSync = (s: boolean) => s && init();

    const init = () => {
        remoteApplying = true;

        syncFromY();
        watch(
            list,
            () => {
                if (remoteApplying) return;
                writeLocal();
            },
            { deep: true, flush: 'post' },
        );

        nextTick(() => {
            remoteApplying = false;
        });

        yarr.observe(syncFromY);
        if (getCurrentInstance())
            onUnmounted(() => {
                provider?.off('sync', onSync);
                yarr.unobserve(syncFromY);
            });
    };

    if (provider) {
        provider.on('sync', onSync);
    } else {
        init();
    }
}

/**
 * Composable to sync a Vue Ref whose entire object is replaced (not updated field-by-field)
 * against a Yjs Y.Map<Primitive>. On remote changes, the local ref is replaced with the new object.
 * On local replacement of the object (ref.value = newObj), Yjs map is cleared and re-seeded.
 *
 * @param yMap       - The Y.Map<Primitive> instance to bind.
 * @param objRef     - Ref to the object; assigning a new object to `objRef.value` triggers a full replace.
 * @param filter     - Optional {@link OptsKeyFilter} to whitelist/blacklist keys to sync.
 * @param opts       - {@link YjsSyncOptions} controlling authoritative seed + provider.
 * @returns The same object ref, kept in sync.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function useYObject<T extends Record<string, any>>(
    yMap: Y.Map<Primitive>,
    objRef: Ref<T>,
    filter?: OptsKeyFilter,
    opts: YjsSyncOptions = {},
): Ref<T> {
    const { provider } = opts;
    let remoteApplying = false;
    let observerAttached = false;

    /**
     * Apply filter to a plain object, returning a new object with allowed keys only.
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const applyFilter = (obj: Record<string, any>): Record<string, any> => {
        const allKeys = Object.keys(obj);
        const keysToSync = filter?.only
            ? filter.only
            : filter?.omit
              ? allKeys.filter((k) => !filter.omit!.includes(k))
              : allKeys;
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const result: Record<string, any> = {};
        keysToSync.forEach((k) => {
            if (obj[k] !== undefined) result[k] = obj[k];
        });
        return result;
    };

    /**
     * Pull entire Y.Map into local object ref, replacing objRef.value.
     */
    const pullRemote = () => {
        remoteApplying = true;
        // Get all primitive key-values from Yjs
        const remoteObj = yMap.toJSON() as Record<string, Primitive>;
        // Merge primitives into the existing local object, removing any primitive keys not present remotely
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const updated: Record<string, any> = { ...objRef.value };
        // First, clear out all keys that are in the filter set
        const allKeys = filter?.only
            ? filter.only
            : filter?.omit
              ? Object.keys(objRef.value).filter((k) => !filter.omit!.includes(k))
              : Object.keys(objRef.value);
        allKeys.forEach((k) => {
            // eslint-disable-next-line no-prototype-builtins
            if (remoteObj.hasOwnProperty(k)) {
                updated[k] = remoteObj[k];
            } else {
                // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
                delete updated[k];
            }
        });
        objRef.value = updated as T;
        nextTick(() => {
            remoteApplying = false;
        });
    };

    /**
     * Attach a Yjs observer so remote changes call pullRemote().
     */
    const attachObserver = () => {
        if (observerAttached) return;
        const handle = (_evt: Y.YMapEvent<Primitive>, tr: Y.Transaction) => {
            if (tr.origin === LOCAL_ORIGIN) return;
            pullRemote();
        };
        yMap.observe(handle);
        if (getCurrentInstance()) {
            onUnmounted(() => {
                yMap.unobserve(handle);
                provider?.off('sync', onSync);
            });
        }
        observerAttached = true;
    };

    const onSync = (s: boolean) => s && init();

    // Initialize: if remote has any keys, pull them; otherwise if authoritative, seed from objRef.value.
    const init = () => {
        if (provider && provider.isAuthoritative) {
            remoteApplying = true;
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const filtered = applyFilter(objRef.value as Record<string, any>);
            yMap.doc?.transact(() => {
                Object.entries(filtered).forEach(([k, v]) => {
                    if (typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean') {
                        yMap.set(k, v as Primitive);
                    }
                });
            }, LOCAL_ORIGIN);
        } else {
            pullRemote();
        }
        nextTick(() => {
            remoteApplying = false;
        });

        attachObserver();

        // Watch local ref replacement: when objRef.value is replaced, push full object to Yjs
        watch(
            objRef,
            (newObj) => {
                if (remoteApplying) return;
                // Clear Yjs map and re-seed
                yMap.doc?.transact(() => {
                    // Delete all existing keys
                    yMap.clear();
                    // Insert filtered new object entries
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    const filtered = applyFilter(newObj as Record<string, any>);
                    Object.entries(filtered).forEach(([k, v]) => {
                        if (typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean') {
                            yMap.set(k, v as Primitive);
                        }
                    });
                }, LOCAL_ORIGIN);
            },
            { flush: 'post' },
        );
    };

    // Wait for provider "sync" event before init, or init immediately if no provider
    if (provider) {
        provider.on('sync', onSync);
    } else {
        init();
    }

    return objRef;
}
