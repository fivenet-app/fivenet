<script lang="ts" setup>
import { HomeIcon } from '@heroicons/vue/20/solid';

const pages = [
    { name: "Projects", href: "#", current: false },
    { name: "Project Nero", href: "#", current: true },
];
</script>

<template>
    <header>
        <div class="mx-auto max-w-7xl py-0 px-0 px-2">
            <nav class="flex border-b py-0 border-gray-200 bg-white rounded-md" aria-label="Breadcrumb">
                <ol role="list" class="mx-auto flex w-full max-w-screen-xl space-x-4 px-4 sm:px-6 lg:px-8">
                    <li class="flex">
                        <div class="flex items-center">
                            <a href="#" class="text-gray-400 hover:text-gray-500">
                                <HomeIcon class="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                                <span class="sr-only">Home</span>
                            </a>
                        </div>
                    </li>
                    <li v-for="page in pages" :key="page.name" class="flex">
                        <div class="flex items-center">
                            <svg class="h-full w-6 flex-shrink-0 text-gray-200" viewBox="0 0 24 44"
                                preserveAspectRatio="none" fill="currentColor" aria-hidden="true">
                                <path d="M.293 0l22 22-22 22h1.414l22-22-22-22H.293z" />
                            </svg>
                            <router-link :to="page.href" class="ml-4 text-sm font-medium text-gray-500 hover:text-gray-700"
                                :aria-current="$route.name == page.href ? 'page' : undefined">{{ page.name
                                }}</router-link>
                        </div>
                    </li>
                </ol>
            </nav>
        </div>
    </header>
</template>
