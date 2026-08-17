import type { Perms } from '~~/gen/ts/perms';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

export type CardElement = {
    label: string;
    description?: string;
    to?: string;
    permission?: Perms;
    icon?: string;
    color?: string;
    deletedAt?: Timestamp;
};

export type ToggleItem<V> = { id: number; label: string; value: V };
