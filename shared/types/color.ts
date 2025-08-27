import type { ChipProps } from '@nuxt/ui';

export type Color = {
    label: string;
    chip?: ChipProps;
    class: string;
};

export type RGB = { r: number; g: number; b: number };
