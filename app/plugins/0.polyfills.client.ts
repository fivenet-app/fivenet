export default defineNuxtPlugin({
    async setup() {
        if (!import.meta.client) return;

        // Make sure we're running in a Chromium based browser
        if (getChromeVersion() === false) return;

        console.info('Loading Polyfills for CEF');
        /* FiveM NUI uses Chromium Embedded Framework version 103.x - As of now requires following polyfills:
         *
         * - toSorted polyfill
         * - CSS Blank Pseudo-element
         * - CSS Has Pseudo-class
         */
        if (!Array.prototype.toSorted) {
            console.debug('Polyfill: Adding toSorted');
            Object.defineProperty(Array.prototype, 'toSorted', {
                value: function <T>(this: T[], compareFn?: (a: T, b: T) => number): T[] {
                    return [...this].sort(compareFn);
                },
                writable: true,
                configurable: true,
            });
        }

        // Load and Activate the CSS polyfills
        const [cssBlankPseudoInit, cssHasPseudo] = await Promise.all([
            // @ts-expect-error No type needed for the CSS polyfill
            import('css-blank-pseudo/browser').then((module) => module.default),
            // @ts-expect-error No type needed for the CSS polyfill
            import('css-has-pseudo/browser').then((module) => module.default),
        ]);

        cssBlankPseudoInit();
        cssHasPseudo(document);
    },
});

function getChromeVersion() {
    const raw = navigator.userAgent.match(/Chrom(e|ium)\/([0-9]+)\./);
    if (!raw || !raw[2]) return false;

    return parseInt(raw[2], 10);
}
