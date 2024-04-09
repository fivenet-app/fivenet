<script setup lang="ts">
import type { NavItem } from '@nuxt/content/dist/runtime/types';

const navigation = inject<NavItem[]>('navigation', []);
</script>

<template>
    <UHeader>
        <template #logo>
            <UDropdown v-slot="{ open }" class="w-full" :ui="{ width: 'w-full' }" :popper="{ strategy: 'absolute' }">
                <UButton
                    color="gray"
                    variant="ghost"
                    :class="[open && 'bg-gray-50 dark:bg-gray-800']"
                    class="w-full"
                    to="/help"
                >
                    <img src="/images/logo-200x200.png" class="size-4" />

                    <span class="truncate font-semibold text-gray-900 dark:text-white">FiveNet - {{ $t('common.help') }}</span>
                </UButton>
            </UDropdown>
        </template>

        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" to="/">
                {{ $t('common.back') }}
            </UButton>

            <UColorModeButton />
        </template>

        <template #panel>
            <UNavigationTree :links="mapContentNavigation(navigation)" />
        </template>
    </UHeader>
</template>
