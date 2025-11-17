<script lang="ts" setup>
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import ChangeUsernameModal from '~/components/auth/account/ChangeUsernameModal.vue';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import type { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';

defineProps<{
    account: GetAccountInfoResponse;
}>();

const overlay = useOverlay();

const changeUsernameModal = overlay.create(ChangeUsernameModal);
const changePasswordModal = overlay.create(ChangePasswordModal);
</script>

<template>
    <UPageCard :title="$t('components.auth.AccountInfo.title')" :description="$t('components.auth.AccountInfo.subtitle')">
        <UFormField class="grid grid-cols-2 items-center gap-2" name="username" :label="$t('common.username')">
            <div class="inline-flex w-full justify-between gap-2">
                <span class="truncate">
                    {{ account.account?.username }}
                </span>
                <CopyToClipboardButton v-if="account.account?.username" :value="account.account?.username" />
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="license"
            :label="$t('components.auth.AccountInfo.license')"
        >
            <div class="inline-flex w-full justify-between gap-2">
                <span class="truncate">
                    {{ account.account?.license }}
                </span>

                <CopyToClipboardButton v-if="account.account?.license" :value="account.account?.license" />
            </div>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="change_username"
            :label="$t('components.auth.AccountInfo.change_username')"
        >
            <UButton @click="changeUsernameModal.open()">
                {{ $t('components.auth.AccountInfo.change_username_button') }}
            </UButton>
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="change_password"
            :label="$t('components.auth.AccountInfo.change_password')"
        >
            <UButton @click="changePasswordModal.open()">
                {{ $t('components.auth.AccountInfo.change_password_button') }}
            </UButton>
        </UFormField>
    </UPageCard>
</template>
