<script lang="ts" setup>
import { FunctionalComponent } from 'vue';
import { RoutesNamedLocations } from '~~/.nuxt/typed-router/__routes';
import { ArrowUpRightIcon } from '@heroicons/vue/24/solid';

defineProps<{
    items: { title: string, description: string, href?: RoutesNamedLocations, permission?: string, icon?: FunctionalComponent, iconForeground?: string, iconBackground?: string }[],
    showIcon: boolean,
}>();
</script>

<template>
    <div
        class="overflow-hidden divide-y-4 rounded-lg bg-base-900 shadow-float sm:grid sm:grid-cols-2 sm:gap-1 sm:max-w-6xl sm:mx-auto divide-base-900 sm:divide-y-0">
        <div v-for="(item, itemIdx) in items" v-can="item.permission ?? ''" :key="item.title" :class="[
            itemIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '',
            itemIdx === 1 ? 'sm:rounded-tr-lg' : '',
            itemIdx === items.length - 2 && items.length % 2 === 0 ? 'sm:rounded-bl-lg' : '',
            itemIdx === items.length - 1 && items.length % 2 === 0 ? 'rounded-br-lg' : '',
            itemIdx === items.length - 1 ? 'rounded-bl-lg sm:rounded-bl-none' : '',
            'group relative bg-base-700 p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-neutral',
        ]">
            <div v-if="item.icon">
                <span :class="[item.iconBackground, item.iconForeground, 'inline-flex rounded-lg p-3']">
                    <component :is="item.icon" class="h-auto w-7" aria-hidden="true" />
                </span>
            </div>
            <div class="mt-4">
                <h3 class="text-base font-semibold leading-6 text-neutral">
                    <span v-if="item.href">
                        <NuxtLink :to="item.href" class="focus:outline-none">
                            <!-- Extend touch target to entire panel -->
                            <span class="absolute inset-0" aria-hidden="true" />
                            {{ item.title }}
                        </NuxtLink>
                    </span>
                    <span v-else>
                        <!-- Extend touch target to entire panel -->
                        <span class="absolute inset-0" aria-hidden="true" />
                        {{ item.title }}
                    </span>
                </h3>
                <p class="mt-2 text-sm text-base-200">{{ item.description }}</p>
            </div>
            <span v-if="showIcon" class="absolute pointer-events-none top-6 right-6 text-base-300 group-hover:text-base-200"
                aria-hidden="true">
                <ArrowUpRightIcon class="w-6 h-6" />
            </span>
        </div>
    </div>
</template>
