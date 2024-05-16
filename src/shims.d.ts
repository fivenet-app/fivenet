import { HookResult } from '@nuxt/schema';
import type { Perms } from '~~/gen/ts/perms';

declare module '#app' {
    interface PageMeta {
        title?: string;
        requiresAuth?: boolean;
        redirectIfAuthed?: boolean;
        permission?: Perms | Perms[];
        authOnlyToken?: boolean;
        showCookieOptions?: boolean;
    }
}

declare module 'vue-router' {
    interface RouteMeta {
        title?: string;
        requiresAuth?: boolean;
        redirectIfAuthed?: boolean;
        permission?: Perms | Perms[];
        authOnlyToken?: boolean;
        showCookieOptions?: boolean;
    }
}

export type ProviderConfig = {
    name: string;
    label: string;
    icon?: string;
};

export type LoginConfig = {
    signupEnabled: boolean;
    lastCharLock: boolean;
    providers: ProviderConfig[];
};

export type DiscordConfig = {
    botInviteURL?: string;
};

export type Links = {
    imprint?: string;
    privacyPolicy?: string;
};

export type FeatureGates = {};

export type AppConfig = {
    version: string;
    login: LoginConfig;
    discord: DiscordConfig;
    links: Links;
    featureGates: FeatureGates;
};

// Custom hooks for custom loading bar logic
declare module '#app' {
    interface RuntimeNuxtHooks {
        'data:loading:start': () => HookResult;
        'data:loading:finish': () => HookResult;
        'data:loading:finish_error': () => HookResult;
    }
    interface NuxtHooks {
        'data:loading:start': () => HookResult;
        'data:loading:finish': () => HookResult;
        'data:loading:finish_error': () => HookResult;
    }
}

// It is always important to ensure you import/export something when augmenting a type
export {};
