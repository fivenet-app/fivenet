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

// It is always important to ensure you import/export something when augmenting a type
export {};
