import type { Perms } from '~~/gen/ts/perms';
import 'vue-router';
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

export type WebsiteConfig = {
    links: Links;
    statsPage: boolean = false;
};

export type Links = {
    imprint?: string;
    privacyPolicy?: string;
};

export type FeatureGates = {};

export type GameConfig = {
    unemployedJobName: string;
};

export type OpenClose = { id: number; label: string; closed?: boolean };

// It is always important to ensure you import/export something when augmenting a type
export {};
