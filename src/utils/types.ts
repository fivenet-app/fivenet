import { FunctionalComponent } from 'vue';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';

export type ArrayElement<ArrayType extends readonly unknown[]> = ArrayType extends readonly (infer ElementType)[]
    ? ElementType
    : never;

export type ValueOf<T> = T[keyof T];

export type CardElement = {
    title: string;
    description?: string;
    href?: RoutesNamedLocations;
    permission?: string;
    icon?: FunctionalComponent;
    iconForeground?: string;
    iconBackground?: string;
};

export type CardElements = CardElement[];
