import 'vue-router';
import type { Perms } from '~~/gen/ts/perms';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

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

export type ServerAppConfig = {
    defaultLocale: string;

    login: LoginConfig;
    discord: DiscordConfig;
    website: WebsiteConfig;
    featureGates: FeatureGates;
    game: GameConfig;
    system: SystemConfig;
};

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
    /**
     * @default false
     */
    statsPage: boolean;
};

export type Links = {
    imprint?: string;
    privacyPolicy?: string;
};

export type FeatureGates = {
    imageProxy: boolean;
};

export type GameConfig = {
    unemployedJobName: string;
    startJobGrade: number;
};

export type SystemConfig = {
    bannerMessage?: BannerMessage;
};

export type BannerMessage = {
    id: string;
    title: string;
    icon?: string;
    color?: string;
    createdAt?: Timestamp;
    expiresAt?: Timestamp;
};

export type OpenClose = { id: number; label: string; closed?: boolean };

export type ClassProp = undefined | string | Record<string, boolean> | (string | Record<string, boolean>)[];

// It is always important to ensure you import/export something when augmenting a type
export {};
