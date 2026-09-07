<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getSettingsSystemClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    disabled?: boolean;
}>();

const notifications = useNotificationsStore();

const settingsSystemClient = await getSettingsSystemClient();

function parseEntries(value: string): { userId: number[]; identifiers: string[]; error?: string } {
    const userId: number[] = [];
    const identifiers: string[] = [];
    const entries = value
        .split(/\r?\n/)
        .map((entry) => entry.trim())
        .filter(Boolean);

    for (const entry of entries) {
        if (entry.startsWith('id:')) {
            const rawId = entry.slice(3).trim();
            const id = Number(rawId);
            if (!/^\d+$/.test(rawId) || !Number.isSafeInteger(id) || id <= 0 || id > 2_147_483_647) {
                return { userId: [], identifiers: [], error: 'invalid_id' };
            }
            if (!userId.includes(id)) userId.push(id);
        } else if (entry.startsWith('identifier:')) {
            const identifier = entry.slice('identifier:'.length).trim();
            if (!identifier || identifier.length > 64) {
                return { userId: [], identifiers: [], error: 'invalid_identifier' };
            }
            if (!identifiers.includes(identifier)) identifiers.push(identifier);
        } else {
            return { userId: [], identifiers: [], error: 'invalid_format' };
        }
    }

    if (userId.length + identifiers.length === 0) return { userId, identifiers, error: 'required' };
    if (userId.length + identifiers.length > 50) return { userId, identifiers, error: 'too_many' };

    return { userId, identifiers };
}

const schema = z.object({
    entries: z.string().superRefine((value, ctx) => {
        const result = parseEntries(value);
        if (result.error) {
            ctx.addIssue({
                code: 'custom',
                message: `components.settings.system_status.trigger_user_sync.validation.${result.error}`,
            });
        }
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    entries: '',
});

const isOpen = ref(false);
const canSubmit = ref(true);

async function triggerUserSync(values: Schema): Promise<void> {
    const parsed = parseEntries(values.entries);
    if (parsed.error) return;

    canSubmit.value = false;
    try {
        await settingsSystemClient.triggerUserSync({ userId: parsed.userId, identifiers: parsed.identifiers });
        notifications.add({
            title: { key: 'components.settings.system_status.trigger_user_sync.success.title', parameters: {} },
            description: { key: 'components.settings.system_status.trigger_user_sync.success.description', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        state.entries = '';
        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        useTimeoutFn(() => (canSubmit.value = true), 400);
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => triggerUserSync(event.data), 1000);
</script>

<template>
    <UDrawer
        v-model:open="isOpen"
        :title="$t('components.settings.system_status.trigger_user_sync.title')"
        :close="false"
        :dismissible="canSubmit"
        :ui="{ body: 'sm:mx-auto sm:max-w-2xl sm:w-full' }"
    >
        <UButton
            icon="i-mdi-account-sync"
            :label="$t('components.settings.system_status.trigger_user_sync.open')"
            variant="soft"
            :disabled="props.disabled"
        />

        <template #title>
            <div class="flex flex-1 items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.settings.system_status.trigger_user_sync.title') }}
                </h3>
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="isOpen = false"
                />
            </div>
        </template>

        <template #body>
            <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmitThrottle">
                <UFormField name="entries" :label="$t('components.settings.system_status.trigger_user_sync.entries')" required>
                    <UTextarea
                        v-model="state.entries"
                        class="w-full"
                        :rows="8"
                        :placeholder="$t('components.settings.system_status.trigger_user_sync.placeholder')"
                        :disabled="!canSubmit"
                    />
                    <template #description>
                        {{ $t('components.settings.system_status.trigger_user_sync.description') }}
                    </template>
                </UFormField>
                <div class="flex justify-end gap-2">
                    <UButton color="neutral" variant="ghost" :disabled="!canSubmit" @click="isOpen = false">
                        {{ $t('common.cancel') }}
                    </UButton>
                    <UButton type="submit" icon="i-mdi-account-sync" :loading="!canSubmit">
                        {{ $t('components.settings.system_status.trigger_user_sync.submit') }}
                    </UButton>
                </div>
            </UForm>
        </template>
    </UDrawer>
</template>
