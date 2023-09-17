// Add `v-can` directive for easy client-side permission checking
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.directive('can', (el, binding, vnode) => {
        // Ignore undefined/ empty permissions
        if (!binding.value || binding.value === '') {
            vnode.el.hidden = false;
            return;
        }

        vnode.el.hidden = !can(binding.value);
        return;
    });
});
