<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import DocEditor from '../partials/DocEditor.vue';
import { canAccess } from './helpers';

withDefaults(
    defineProps<{
        disabled?: boolean;
    }>(),
    {
        disabled: false,
    },
);

const { isOpen } = useModal();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const schema = z.object({
    signature: z.string().max(1024),
    emails: z.string().array().max(25),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    signature: '',
    emails: [],
});

const canManage = computed(() => canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.MANAGE));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!selectedEmail.value?.id) {
        return;
    }
    canSubmit.value = false;

    const values = event.data;
    await mailerStore
        .setEmailSettings({
            settings: {
                emailId: selectedEmail.value?.id,
                signature: values.signature,
                blockedEmails: values.emails.map((e) => e.trim()),
            },
        })
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));

    isOpen.value = false;
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.settings') }} - {{ selectedEmail?.email }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div class="flex flex-col gap-2">
                    <UFormGroup name="signature" class="flex-1" :label="$t('common.signature')">
                        <ClientOnly>
                            <DocEditor v-model="state.signature" :disabled="disabled || !canManage" :min-height="200" />
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup name="emails" class="flex-1" :label="$t('common.blocklist')">
                        <div class="flex flex-col gap-1">
                            <div v-for="(_, idx) in state.emails" :key="idx" class="flex items-center gap-1">
                                <UFormGroup :name="`emails.${idx}`" class="flex-1">
                                    <UInput
                                        v-model="state.emails[idx]"
                                        type="text"
                                        :placeholder="$t('common.mail')"
                                        :disabled="disabled || !canManage"
                                    />
                                </UFormGroup>

                                <UButton
                                    :ui="{ rounded: 'rounded-full' }"
                                    icon="i-mdi-close"
                                    :disabled="disabled || !canSubmit"
                                    @click="state.emails.splice(idx, 1)"
                                />
                            </div>
                        </div>

                        <UButton
                            v-if="!disabled || canManage"
                            :ui="{ rounded: 'rounded-full' }"
                            icon="i-mdi-plus"
                            :disabled="disabled || !canSubmit || state.emails.length >= 25"
                            :class="state.emails.length ? 'mt-2' : ''"
                            @click="state.emails.push('')"
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            v-if="!disabled || canManage"
                            type="submit"
                            :label="$t('common.save')"
                            block
                            class="flex-1"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                        />
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
