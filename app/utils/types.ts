import type { RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

export type CardElement = {
    title: string;
    description?: string;
    to?: RoutesNamedLocations | string;
    permission?: Perms;
    icon?: string;
    color?: string;
};

export type ToggleItem<V> = { id: number; label: string; value: V };
