<script lang="ts" setup>
import { ChevronRightIcon } from 'mdi-vue3';
import { type CardElement } from '~/utils/types';

withDefaults(
    defineProps<{
        items: CardElement[];
        showIcon?: boolean;
    }>(),
    {
        showIcon: true,
    },
);

defineEmits<{
    (e: 'selected', idx: number): void;
}>();
</script>

<template>
    <div class="sm:px-2">
        <div
            class="w-full overflow-hidden divide-y-4 rounded-lg bg-base-900 sm:grid sm:gap-1 divide-accent-900 sm:divide-y-0"
            :class="[items.length === 1 ? '' : 'sm:grid-cols-2 sm:max-w-6xl']"
        >
            <template v-for="(item, itemIdx) in items">
                <div
                    v-if="can(item.permission ?? '')"
                    :key="item.title"
                    :class="[
                        itemIdx === 0 ? 'rounded-tl-lg rounded-tr-lg sm:rounded-tr-none' : '',
                        itemIdx === 1 ? 'sm:rounded-tr-lg' : '',
                        itemIdx === items.length - 2 && itemIdx % 2 === 1 ? 'sm:rounded-br-lg' : '',
                        itemIdx === items.length - 1 && itemIdx % 2 === 0 ? 'rounded-br-lg' : '',
                        itemIdx === items.length - 1 ? 'rounded-bl-lg sm:rounded-bl-none' : '',
                        'group relative bg-base-700 p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-neutral',
                    ]"
                >
                    <div v-if="item.icon">
                        <span :class="[item.iconBackground, item.iconForeground, 'inline-flex rounded-lg p-3']">
                            <component :is="item.icon" class="h-auto w-7" aria-hidden="true" />
                        </span>
                    </div>
                    <div class="mt-4" @click="$emit('selected', itemIdx)">
                        <h3 class="text-base font-semibold leading-6 text-neutral">
                            <template v-if="item.href">
                                <NuxtLink :to="item.href" class="focus:outline-none">
                                    <!-- Extend touch target to entire panel -->
                                    <span class="absolute inset-0" aria-hidden="true" />
                                    {{ item.title }}
                                </NuxtLink>
                            </template>
                            <template v-else>
                                <!-- Extend touch target to entire panel -->
                                <span class="absolute inset-0" aria-hidden="true" />
                                {{ item.title }}
                            </template>
                        </h3>
                        <p class="mt-2 text-sm text-base-200">
                            {{ item.description }}
                        </p>
                    </div>
                    <span
                        v-if="showIcon"
                        class="absolute pointer-events-none top-6 right-6 text-base-300 group-hover:text-base-200"
                        aria-hidden="true"
                    >
                        <ChevronRightIcon class="w-6 h-6" />
                    </span>
                </div>
            </template>
        </div>
    </div>
</template>
