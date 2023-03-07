<script lang="ts">
import { defineComponent } from 'vue';

import { ChevronRightIcon, HomeIcon } from '@heroicons/vue/20/solid'

export default defineComponent({
    components: {
        ChevronRightIcon,
        HomeIcon,
    },
    props: {
        title: {
            required: true,
            type: String,
        },
    },
    data() {
        return {
            homepage: { name: 'Overview', href: '/overview', },
        };
    },
    computed: {
        breadCrumbs() {
            if (this.$route.meta.breadCrumbs === null) {
                return null;
            }

            return this.$route.meta.breadCrumbs;
        },
    },
});
</script>

<template>
    <header class="bg-gray-700 shadow-sm">
        <div class="mx-auto max-w-7xl py-4 px-4 sm:px-6 lg:px-8">
            <div class="grid gap-4 grid-cols-2 grid-rows-1">
                <h1 class="text-lg font-semibold leading-6 text-white">{{ title }}</h1>
                <nav v-if="breadCrumbs" class="flex" aria-label="Breadcrumb">
                    <ol role="list" class="flex items-center space-x-4">
                        <li>
                            <div>
                                <router-link :to="homepage.href"
                                    class="text-gray-200 hover:text-gray-400"
                                    active-class="underline decoration-solid font-bold"
                                    :aria-current="$route.path == homepage.href ? 'page' : undefined">
                                    <HomeIcon class="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                                    <span class="sr-only">{{ homepage.name }}</span>
                                </router-link>
                            </div>
                        </li>
                        <li v-for="page in breadCrumbs" :key="page.name">
                            <div class="flex items-center">
                                <ChevronRightIcon class="h-5 w-5 flex-shrink-0 text-gray-200" aria-hidden="true" />
                                <router-link :to="page.href ? page.href : '#'"
                                    class="ml-4 text-sm font-medium text-gray-200 hover:text-gray-400"
                                    active-class="underline decoration-solid font-bold"
                                    :aria-current="$route.path == page.href ? 'page' : undefined">{{ page.name
                                    }}</router-link>
                            </div>
                        </li>
                    </ol>
                </nav>
            </div>
        </div>
    </header>
</template>
