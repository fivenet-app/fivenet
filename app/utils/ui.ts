import { nuxtUiColors, type NuxtUIColor } from '~/utils/color';

function isNuxtUIColor(color: string): color is NuxtUIColor {
    return nuxtUiColors.includes(color as NuxtUIColor);
}

export function stringToButtonColor(s: string): NuxtUIColor {
    const color = s.trim().toLowerCase();

    return isNuxtUIColor(color) ? color : 'primary';
}
