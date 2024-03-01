import type { Perms } from '~~/gen/ts/perms';

declare module '#app' {
    interface PageMeta {
        title?: string;
        requiresAuth?: boolean;
        permission?: Perms | Perms[];
        authOnlyToken?: boolean;
        showQuickButtons?: boolean;
        showCookieOptions?: boolean;
        breadcrumbTitle?: string;
    }
}

declare module 'vue-router' {
    interface RouteMeta {
        title?: string;
        requiresAuth?: boolean;
        permission?: Perms;
        authOnlyToken?: boolean;
        showQuickButtons?: boolean;
        showCookieOptions?: boolean;
        breadcrumbTitle?: string;
    }
}

type ProviderConfig = {
    name: string;
    label: string;
};

type LoginConfig = {
    signupEnabled: boolean;
    providers: ProviderConfig[];
};

type DiscordConfig = {
    botInviteURL?: string;
};

type Links = {
    imprint?: string;
    privacyPolicy?: string;
};

type FeatureGates = {};

export type AppConfig = {
    version: string;
    login: LoginConfig;
    discord: DiscordConfig;
    links: Links;
    featureGates: FeatureGates;
};

// It is always important to ensure you import/export something when augmenting a type
export {};
