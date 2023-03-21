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

// manual extension of route types
declare module 'vue-router/auto/routes' {
    import type {
        RouteRecordInfo,
        ParamValue,
        ParamValueOneOrMore,
        ParamValueZeroOrMore,
        ParamValueZeroOrOne,
    } from 'unplugin-vue-router';

    export interface RouteNamedMap {
        'custom-dynamic-name': RouteRecordInfo<
            'custom-dynamic-name',
            '/added-during-runtime/[...path]',
            { id: ParamValue<true> },
            { path: ParamValue<false> }
        >;
    }
}
