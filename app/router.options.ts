import type { RouterConfig } from '@nuxt/schema';

function cssAttrEscape(value: string) {
    // good enough for attribute selector values
    return value.replace(/\\/g, '\\\\').replace(/"/g, '\\"');
}

function decodeHash(hash: string) {
    return decodeURIComponent(hash.startsWith('#') ? hash.slice(1) : hash);
}

function findIdAnchor(hash: string): { el: string; top: number } | undefined {
    // Try normal #id selector first (may throw if hash is not a valid selector)
    let el: HTMLElement | null = null;
    try {
        el = document.querySelector(hash) as HTMLElement | null;
    } catch {
        el = null;
    }
    if (!el) return;

    const top = parseFloat(getComputedStyle(el).scrollMarginTop) || 0;
    return { el: hash, top };
}

function findDataIdElement(hash: string): { el: HTMLElement; top: number } | undefined {
    const id = decodeHash(hash);
    if (!id) return;

    const el = document.querySelector<HTMLElement>(`[data-id="${cssAttrEscape(id)}"]`);
    if (!el) return;

    const top = parseFloat(getComputedStyle(el).scrollMarginTop) || 0;
    return { el, top };
}

function manualScrollIntoView(el: HTMLElement, topOffset: number) {
    // 1) Scroll target into view
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });

    // 2) Apply scroll-margin-top offset (router doesn’t do this reliably for manual cases)
    if (topOffset) {
        // Let the initial scrollIntoView settle a moment, then nudge up
        requestAnimationFrame(() => {
            window.scrollBy({ top: -topOffset, behavior: 'smooth' });
        });
    }
}

// https://router.vuejs.org/api/#routeroptions
export default <RouterConfig>{
    scrollBehavior(to, from, savedPosition) {
        const nuxtApp = useNuxtApp();

        if (history.state && history.state.stop) {
            return;
        }
        if (history.state && history.state.smooth) {
            return {
                el: history.state.smooth,
                behavior: 'smooth',
            };
        }

        // If history back
        if (savedPosition) {
            // Handle Suspense resolution
            return new Promise((resolve) => {
                nuxtApp.hooks.hookOnce('page:finish', () => {
                    setTimeout(() => resolve(savedPosition), 50);
                });
            });
        }

        // Scroll to heading on click
        if (to.hash) {
            return new Promise((resolve) => {
                const run = () => {
                    // 1) Prefer native id anchors via router (best integration)
                    const idAnchor = findIdAnchor(to.hash);
                    if (idAnchor) {
                        resolve({ el: idAnchor.el, behavior: 'smooth', top: idAnchor.top });
                        return;
                    }

                    // 2) Fallback: data-id -> manual scrollIntoView
                    const dataId = findDataIdElement(to.hash);
                    if (dataId) {
                        manualScrollIntoView(dataId.el, dataId.top);
                        // Tell router "we handled scrolling"
                        resolve(false);
                        return;
                    }

                    // Nothing found; do nothing
                    resolve(false);
                };

                if (to.path === from.path) {
                    setTimeout(run, 50);
                } else {
                    nuxtApp.hooks.hookOnce('page:finish', () => {
                        setTimeout(run, 50);
                    });
                }
            });
        }

        // Scroll to top of window
        return { top: 0 };
    },
};
