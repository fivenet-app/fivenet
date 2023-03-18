import { createApp } from 'vue';
import App from './App.vue';
import * as Sentry from '@sentry/vue';
import config from './config';
import router from './router';
import store from './store';
import slugify from 'slugify';

// Load styles and Inter font (all weights)
import './style.css';
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

app.use(router);
app.use(store);

// Add `v-can` directive for easy permission checking
app.directive('can', (el, binding, vnode) => {
    const permissions = store.state.permissions;
    const val = slugify(binding.value as string);
    if (permissions.includes(val) || val === '') {
        return (vnode.el.hidden = false);
    } else {
        return (vnode.el.hidden = true);
    }
});

app.mount('#app');
