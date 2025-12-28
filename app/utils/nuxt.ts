import type { NavigationMenuItem } from '@nuxt/ui';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Handler = (e?: any) => void;

export function extractShortcutsFromNavItems(
    items: NavigationMenuItem[],
    separator?: Parameters<typeof extractShortcuts>[1],
): Record<string, Handler> {
    return extractShortcuts(
        items.map((item) => ({
            ...item,
            onSelect: item.to
                ? () =>
                      navigateTo(
                          // @ts-expect-error ignore route type (is string here)
                          item.to,
                      )
                : item.onSelect,
        })),
        separator,
    );
}
