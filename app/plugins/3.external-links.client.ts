export default defineNuxtPlugin({
    name: 'external-links',
    parallel: true,

    setup() {
        const clickListener = (event: MouseEvent): void => {
            if (!event.target || event.defaultPrevented) return;

            let element: HTMLElement | null = event.target as HTMLElement;
            for (; element && element !== document.body; element = element.parentElement as HTMLElement) {
                if (element.tagName.toLowerCase() === 'a' || element.hasAttribute('href')) break;
            }
            if (!element) return;

            const href = element.getAttribute('href');
            if (!href || href.startsWith('/') || href.startsWith('#') || href === '') return;

            event.preventDefault();
            navigateTo({
                name: 'dereferer',
                query: {
                    target: href,
                },
            });
        };

        useEventListener(window, 'click', clickListener);
    },
});
