import { defineNuxtModule, addPluginTemplate, createResolver, findPath } from '@nuxt/kit';
import fs from 'fs';
import type { Nuxt } from 'nuxt/schema';
import path from 'path';

// Define a Nuxt module that sets up a custom plugin template for the SPA loading screen.
export default defineNuxtModule({
    // The setup function is called with the Nuxt instance and module context.
    defaults: {
        delay: 25 // Default delay `ms
    },

    setup(options, nuxt) {
        // Create a resolver for the current module URL.
        const { resolve } = createResolver(import.meta.url);

        // Register a hook that runs after all modules are processed.
        nuxt.hook('modules:done', async () => {
            async function spaLoadingTemplatePath(nuxt: Nuxt) {
                if (typeof nuxt.options.spaLoadingTemplate === 'string') {
                    return resolve(nuxt.options.srcDir, nuxt.options.spaLoadingTemplate);
                }

                const possiblePaths = nuxt.options._layers.map((layer) =>
                    resolve(layer.config.srcDir, layer.config.dir?.app || 'app', 'spa-loading-template.html')
                );

                return (
                    (await findPath(possiblePaths)) ??
                    resolve(nuxt.options.srcDir, nuxt.options.dir?.app || 'app', 'spa-loading-template.html')
                );
            }

            // Read the contents of the SPA loading template HTML file
            const loadingTemplateContent = fs.readFileSync(
                path.resolve(__dirname, await spaLoadingTemplatePath(nuxt)),
                'utf-8'
            );

            // Add the custom plugin template for the SPA loading screen.
            addPluginTemplate({
                getContents({ options }) {
                    const contents = fs.readFileSync(resolve('./runtime/plugin.ts'), 'utf-8');
                    type SpaLoaderOption = 'delay' | 'html';
                    return contents.replace(
                        /<%= options\.([^ ]+) %>/g,
                        (_, option: SpaLoaderOption) => options[option]
                    );
                },
                filename: './runtime/plugin.ts',
                mode: 'client',
                options: {
                    html: loadingTemplateContent,
                    delay: options?.delay || 0
                }
            });
        });
    }
});
