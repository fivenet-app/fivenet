declare module '#app' {
    interface PageMeta {
        title?: string;
        requiresAuth?: boolean;
        permission?: string | string[];
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
        permission?: string;
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

export type AppConfig = {
    version: string;
    sentryDSN?: string;
    sentryEnv?: string;
    login: LoginConfig;
    discord: DiscordConfig;
    links: Links;
};

// It is always important to ensure you import/export something when augmenting a type
export {};
