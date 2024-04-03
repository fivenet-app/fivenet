<script setup lang="ts">
import { useAuthStore } from '~/store/auth';
import FiveNetLogo from '../partials/logos/FiveNetLogo.vue';
import LanguageSwitcherModal from '../partials/LanguageSwitcherModal.vue';

const { t } = useI18n();

const authStore = useAuthStore();
const { accessToken } = storeToRefs(authStore);

const links = [
    {
        label: t('common.home'),
        to: '/',
    },
    {
        label: t('common.about'),
        to: '/about',
    },
];

const modal = useModal();
</script>

<template>
    <UHeader :links="links">
        <template #logo> <FiveNetLogo class="h-10 w-auto" /> </template>

        <template #right>
            <UButton
                :label="$t('common.language')"
                icon="i-mdi-translate"
                color="gray"
                @click="modal.open(LanguageSwitcherModal, {})"
            />

            <template v-if="!accessToken">
                <UButton :label="$t('components.auth.login.title')" icon="i-mdi-login" color="gray" to="/auth/login" />
                <UButton
                    :label="$t('components.auth.registration_form.title')"
                    icon="i-mdi-account-plus"
                    trailing
                    color="black"
                    to="/auth/registration"
                    class="hidden lg:flex"
                />
            </template>
            <template v-else>
                <UButton :label="$t('common.overview')" icon="i-mdi-home" to="/overview" />
            </template>
        </template>
    </UHeader>
</template>
