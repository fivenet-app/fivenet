import { createApp } from 'vue';
import App from './App.vue';
import * as Sentry from '@sentry/vue';
import config from './config';
import { LoadingPlugin } from 'vue-loading-overlay';
import { router } from './router';
import { store, key } from './store/store';
import slug from './utils/slugify';

// Load styles and Inter font (all weights)
import './style.css';
import 'vue-loading-overlay/dist/css/index.css';
import '@fontsource/inter/100.css';
import '@fontsource/inter/200.css';
import '@fontsource/inter/300.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/inter/700.css';
import '@fontsource/inter/800.css';
import '@fontsource/inter/900.css';

const app = createApp(App);

Sentry.init({
    app,
    dsn: config.sentryDSN,
    tracesSampleRate: 0.0,
    logErrors: true,
    trackComponents: false,
});

app.use(LoadingPlugin);
app.use(router);
app.use(store, key);

// Add `v-can` directive for easy permission checking
app.directive('can', (el, binding, vnode) => {
    const permissions = store.state.auth?.permissions;
    const val = slug(binding.value as string);
    if (permissions && (permissions.includes(val) || val === '')) {
        return (vnode.el.hidden = false);
    } else {
        return (vnode.el.hidden = true);
    }
});

app.mount('#app');
