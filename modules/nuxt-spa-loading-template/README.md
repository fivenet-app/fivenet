# Nuxt 3 Custom SPA Loading Module

This Nuxt 3 module provides a temporary solution for the issue identified in [Nuxt.js GitHub Issue #21721](https://github.com/nuxt/nuxt/issues/21721).

## Installation

Please note that this module is not officially part of Nuxt 3 and is intended as a temporary workaround for the
specified issue. It may not be actively maintained, and its usage is at your discretion.

To manually install and use this module, follow these steps:

1. Clone or download this repository to your local machine.

2. Follow the standard Nuxt 3 module installation procedure to add this module to your Nuxt project.

    - In your Nuxt 3 project directory, create a folder named `modules` if it doesn't already exist.

    - Place the downloaded module files into the `modules` directory.

3. Configure your Nuxt project to use the module by adding it to the `nuxt.config.ts` file.

```ts
export default defineNuxtConfig({
  // Other Nuxt configuration options...
  spaLoadingTemplate: 'loader.html', // Start loader SPA  
  modules: [
    // Add the module by specifying its path in your project.
    '@/modules/spa-loading-template/module'
  ],
  // Other Nuxt configuration options...
})
```
