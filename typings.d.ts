// typings.d.ts or router.ts
import 'vue-router/auto';

declare module 'vue-router/auto' {
    interface RouteMeta {
        requiresAuth?: boolean
    }
}
