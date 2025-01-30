export default defineNuxtPlugin(() => {
    if (!import.meta.client) {
        return;
    }

    // toSorted polyfill
    if (!Array.prototype.toSorted) {
        console.info('Adding Polyfill');
        Object.defineProperty(Array.prototype, 'toSorted', {
            value: function <T>(this: T[], compareFn?: (a: T, b: T) => number): T[] {
                return [...this].sort(compareFn);
            },
            writable: true,
            configurable: true,
        });
    }
});
