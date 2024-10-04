import { defineNuxtPlugin } from '#app';

// This Nuxt plugin performs the following actions:
// 1. Creates and mounts a component in the DOM to display a loading state.
// 2. Removes the loading component after the page loading is completed.

export default defineNuxtPlugin((nuxtApp) => {
    // Check if the 'document' object and its 'body' property exist.
    if (document && document.body) {
        // Create a 'div' element to display the loading component.
        const loaderHtmlComponent = document.createElement('div');
        // Insert the HTML code for the loading component, specified in the plugin options.
        loaderHtmlComponent.innerHTML = `<%= options.html %>`;
        // Insert the loading component at the beginning of the document's body.
        document.body.insertBefore(loaderHtmlComponent, document.body.firstChild);

        function removeHTMLComponent() {
            if (loaderHtmlComponent) {
                // Remove the loading component.
                loaderHtmlComponent.remove();
            }
        }

        // Register the 'page:finish' hook to perform actions after page loading is complete.
        nuxtApp.hook('page:finish', () => {
            setTimeout(() => {
                // Remove the loading component.
                removeHTMLComponent();
            }, Number(`<%= options.delay %>`));
        });

        // Remove element after fatal error occures
        nuxtApp.hook('app:error', () => {
            // Remove the loading component.
            removeHTMLComponent();
        });
    }
});
