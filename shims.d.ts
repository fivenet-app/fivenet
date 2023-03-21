/// <reference types="vite/client" />

import 'vue-router/auto';

declare module 'vue-router/auto' {
    interface RouteMeta {
        requiresAuth?: boolean;
        permission?: String;
        breadCrumbs?: null | Array<BreadCrumbPart>;
    }

    interface BreadCrumbPart {
        name: string;
        href?: RouteNamedMap;
    }
}
