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

// Manual extension of route types
declare module 'vue-router/auto/routes' {
    import type { RouteRecordInfo, ParamValue } from 'unplugin-vue-router';

    export interface RouteNamedMap {
        'Citizens: Info': RouteRecordInfo<
            'Citizens: Info',
            '/citizens/:id',
            { id: ParamValue<true> },
            { id: ParamValue<false> }
        >;
        'Documents: Info': RouteRecordInfo<
            'Documents: Info',
            '/documents/:id',
            { id: ParamValue<true> },
            { id: ParamValue<false> }
        >;
    }
}
