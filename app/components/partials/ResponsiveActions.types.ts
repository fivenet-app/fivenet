import type { ButtonProps } from '@nuxt/ui';
import type { RouteLocationRaw } from 'vue-router';

export type ResponsiveActionItem = {
    kind?: 'action';
    label: string;
    tooltip?: string;
    kbds?: string[];
    icon?: string;
    class?: string;
    color?: ButtonProps['color'];
    variant?: ButtonProps['variant'];
    disabled?: boolean;
    to?: RouteLocationRaw;
    onClick?: () => void | Promise<void>;
    children?: ResponsiveActionItem[];
};

export type ResponsiveActionSeparator = {
    kind: 'separator';
};

export type ResponsiveActionEntry = ResponsiveActionItem | ResponsiveActionSeparator;

export function separator(): ResponsiveActionSeparator {
    return { kind: 'separator' };
}
