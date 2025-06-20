import 'vue-router';
import type { Perms } from '~~/gen/ts/perms';

declare module 'vue-router' {
    interface RouteMeta {
        title?: string;
        requiresAuth?: boolean;
        redirectIfAuthed?: boolean;
        permission?: Perms | Perms[];
        authTokenOnly?: boolean;
        showCookieOptions?: boolean;
    }
}

declare module '@nuxtjs/i18n' {
    interface LocaleObject {
        icon: string;
    }
}

export type ToggleItem = { id: number; label: string; value?: boolean };

export type ClassProp = undefined | string | Record<string, boolean> | (string | Record<string, boolean>)[];

// It is always important to ensure you import/export something when augmenting a type
export {};
