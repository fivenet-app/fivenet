import type { RoutesNamedLocations } from '@typed-router';
import type { Perms } from '~~/gen/ts/perms';

export type ArrayElement<ArrayType extends readonly unknown[]> = ArrayType extends readonly (infer ElementType)[]
    ? ElementType
    : never;

export type ValueOf<T> = T[keyof T];

export type CardElement = {
    title: string;
    description?: string;
    to?: RoutesNamedLocations;
    permission?: Perms;
    icon?: string;
};

export type CardElements = CardElement[];

export function isNumber(value?: string | number): boolean {
    return value !== undefined && value !== null && value !== '' && !isNaN(Number(value.toString()));
}
