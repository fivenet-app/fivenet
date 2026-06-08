import { isProxy, isReactive, isRef, toRaw } from 'vue';

export function deepToRaw<T>(sourceObj: T): T {
    const objectIterator = (input: unknown): unknown => {
        if (Array.isArray(input)) {
            return input.map((item) => objectIterator(item));
        }

        if (isRef(input)) {
            return objectIterator(input.value);
        }

        if (isReactive(input) || isProxy(input)) {
            return objectIterator(toRaw(input as object));
        }

        if (input && typeof input === 'object') {
            const proto = Object.getPrototypeOf(input);
            if (proto !== Object.prototype && proto !== null) {
                return input;
            }

            return Object.keys(input).reduce(
                (acc, key) => {
                    (acc as Record<string, unknown>)[key] = objectIterator((input as Record<string, unknown>)[key]);
                    return acc;
                },
                {} as Record<string, unknown>,
            );
        }

        return input;
    };

    return objectIterator(sourceObj) as T;
}
