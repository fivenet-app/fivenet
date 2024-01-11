import { type DefineComponent } from 'vue';
import { type RoutesNamedLocations } from '@typed-router';

export type ArrayElement<ArrayType extends readonly unknown[]> = ArrayType extends readonly (infer ElementType)[]
    ? ElementType
    : never;

export type ValueOf<T> = T[keyof T];

export type CardElement = {
    title: string;
    description?: string;
    href?: RoutesNamedLocations;
    permission?: string;
    icon?: DefineComponent;
    iconForeground?: string;
    iconBackground?: string;
};

export type CardElements = CardElement[];