import type { DiscordConfig, FeatureGates, Links, LoginConfig } from '~/shims';

export default defineAppConfig({
    version: '',
    login: {
        signupEnabled: true,
        providers: [],
    } as LoginConfig,
    discord: {
        botInviteURL: '',
    } as DiscordConfig,
    links: {} as Links,
    featureGates: {} as FeatureGates,

    // Nuxt UI app config
    ui: {
        primary: 'blue',
        gray: 'cool',
        tooltip: {
            default: {
                openDelay: 500,
            },
        },
        button: {
            default: {
                loadingIcon: 'i-mdi-loading',
            },
        },
        table: {
            td: {
                padding: 'px-1.5 py-1.5',
            },
        },
        selectMenu: {
            default: {
                selectedIcon: 'i-mdi-check',
            },
        },
        card: {
            footer: {
                padding: 'px-2 py-3 sm:px-4',
            },
        },
        icons: {
            dynamic: true,
            // Nuxt UI Pro Icons
            dark: 'i-mdi-moon-and-stars',
            light: 'i-mdi-weather-sunny',
            system: 'i-mdi-computer',
            search: 'i-mdi-search',
            external: 'i-mdi-external-link',
            chevron: 'i-mdi-chevron-down',
            hash: 'i-mdi-hashtag',
            menu: 'i-mdi-menu',
            close: 'i-mdi-window-close',
            check: 'i-mdi-check-circle',
        },
    },

    filestore: {
        fileSizes: {
            rector: 5 * 1024 * 1024,
            images: 2 * 1024 * 1024,
        },
        types: {
            images: ['image/jpeg', 'image/jpg', 'image/png'],
        },
    },

    seo: {
        siteName: 'Nuxt UI Pro - Docs template',
    },
    header: {
        logo: {
            alt: '',
            light: '',
            dark: '',
        },
        search: true,
        colorMode: true,
        links: [
            {
                icon: 'i-simple-icons-github',
                to: 'https://github.com/nuxt-ui-pro/docs',
                target: '_blank',
                'aria-label': 'Docs template on GitHub',
            },
        ],
    },
    footer: {
        credits: 'Copyright © 2023',
        colorMode: false,
        links: [
            {
                icon: 'i-simple-icons-nuxtdotjs',
                to: 'https://nuxt.com',
                target: '_blank',
                'aria-label': 'Nuxt Website',
            },
            {
                icon: 'i-simple-icons-discord',
                to: 'https://discord.com/invite/ps2h6QT',
                target: '_blank',
                'aria-label': 'Nuxt UI on Discord',
            },
            {
                icon: 'i-simple-icons-x',
                to: 'https://x.com/nuxt_js',
                target: '_blank',
                'aria-label': 'Nuxt on X',
            },
            {
                icon: 'i-simple-icons-github',
                to: 'https://github.com/nuxt/ui',
                target: '_blank',
                'aria-label': 'Nuxt UI on GitHub',
            },
        ],
    },
    toc: {
        title: 'Table of Contents',
        bottom: {
            title: 'Community',
            edit: 'https://github.com/nuxt-ui-pro/docs/edit/main/content',
            links: [
                {
                    icon: 'i-heroicons-star',
                    label: 'Star on GitHub',
                    to: 'https://github.com/nuxt/ui',
                    target: '_blank',
                },
                {
                    icon: 'i-heroicons-book-open',
                    label: 'Nuxt UI Pro docs',
                    to: 'https://ui.nuxt.com/pro/guide',
                    target: '_blank',
                },
                {
                    icon: 'i-simple-icons-nuxtdotjs',
                    label: 'Purchase a license',
                    to: 'https://ui.nuxt.com/pro/purchase',
                    target: '_blank',
                },
            ],
        },
    },
});
