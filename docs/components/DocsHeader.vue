<script setup lang="ts">
import type { NavItem } from '@nuxt/content/dist/runtime/types';
import LanguageSwitcherModal from './LanguageSwitcherModal.vue';

const navigation = inject<NavItem[]>('navigation', []);

const modal = useModal();
</script>

<template>
    <UHeader>
        <template #logo>
            <ULink class="inline-flex w-full items-center gap-2" to="/">
                <img src="/images/logo-200x200.png" class="size-10" />

                <span class="truncate font-semibold text-gray-900 dark:text-white">FiveNet</span>
            </ULink>
        </template>

        <template #center>
            <UContentSearchButton class="hidden lg:flex" />
        </template>

        <template #right>
            <UContentSearchButton :label="null" class="lg:hidden" />

            <UButton
                :label="$t('common.language')"
                icon="i-mdi-translate"
                color="gray"
                @click="modal.open(LanguageSwitcherModal, {})"
            />

            <UColorModeButton />
        </template>

        <template #panel>
            <UNavigationTree :links="mapContentNavigation(navigation)" />
        </template>
    </UHeader>
</template>
