export default defineNuxtPlugin(() => {
    if (!import.meta.client) {
        return;
    }

    /* FiveM NUI uses Chromium Embedded Framework version 103.x - As of now requires following polyfills:
     *
     * toSorted polyfill
     */
    if (!Array.prototype.toSorted) {
        console.info('Adding toSorted Polyfill');
        Object.defineProperty(Array.prototype, 'toSorted', {
            value: function <T>(this: T[], compareFn?: (a: T, b: T) => number): T[] {
                return [...this].sort(compareFn);
            },
            writable: true,
            configurable: true,
        });
    }
});
