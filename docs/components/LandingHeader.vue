<script setup lang="ts">
import LanguageSwitcherModal from './LanguageSwitcherModal.vue';

const { t } = useI18n();

const links = [
    {
        label: t('common.home'),
        to: '/',
    },
    {
        label: t('common.docs'),
        to: '/docs',
    },
];

const modal = useModal();
</script>

<template>
    <UHeader :links="links">
        <template #logo>
            <ULink class="inline-flex w-full items-center gap-2" to="/">
                <img src="/images/logo-200x200.png" class="size-10" />

                <span class="truncate font-semibold text-gray-900 dark:text-white">FiveNet</span>
            </ULink>
        </template>

        <template #right>
            <UButton
                :label="$t('common.language')"
                icon="i-mdi-translate"
                color="gray"
                @click="modal.open(LanguageSwitcherModal, {})"
            />

            <UButton :label="$t('common.docs')" icon="i-mdi-book-open-variant-outline" to="/docs" />
        </template>
    </UHeader>
</template>
