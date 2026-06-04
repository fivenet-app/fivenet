<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getSettingsLawsClient } from '~~/gen/ts/clients';
import type { Law, LawBook } from '~~/gen/ts/resources/laws/laws';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    law: Law;
    lawBooks: LawBook[];
}>();

const emit = defineEmits<{
    (e: 'update:law', update: { id: number; law: Law }): void;
    (e: 'close'): void;
}>();

const { display } = useAppConfig();

const notifications = useNotificationsStore();

const settingsLawsClient = await getSettingsLawsClient();

const schema = z.object({
    lawbookId: z.coerce.number().int().positive(),
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(3).max(1024), z.string().length(0).optional()]),
    hint: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    fine: z.coerce.number().nonnegative().max(999_999_999),
    detentionTime: z.coerce.number().nonnegative().max(999_999_999),
    stvoPoints: z.coerce.number().nonnegative().max(999_999_999),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    lawbookId: props.law.lawbookId,
    name: props.law.name,
    description: props.law.description,
    hint: props.law.hint,
    fine: props.law.fine ?? 0,
    detentionTime: props.law.detentionTime ?? 0,
    stvoPoints: props.law.stvoPoints ?? 0,
});

const availableLawBooks = computed(() => props.lawBooks.filter((book) => book.deletedAt === undefined && book.id > 0));
const currentLawBook = computed(() => props.lawBooks.find((book) => book.id === state.lawbookId));

function resetForm(): void {
    state.lawbookId = props.law.lawbookId;
    state.name = props.law.name;
    state.description = props.law.description;
    state.hint = props.law.hint;
    state.fine = props.law.fine ?? 0;
    state.detentionTime = props.law.detentionTime ?? 0;
    state.stvoPoints = props.law.stvoPoints ?? 0;
}

onBeforeMount(() => resetForm());
watch(
    () => props.law,
    () => resetForm(),
    { deep: true },
);

async function saveLaw(id: number, values: Schema): Promise<void> {
    try {
        const call = settingsLawsClient.createOrUpdateLaw({
            law: {
                id: id < 0 ? 0 : id,
                lawbookId: values.lawbookId,
                name: values.name,
                sortOrder: 0,
                description: values.description,
                hint: values.hint,
                fine: values.fine,
                detentionTime: values.detentionTime,
                stvoPoints: values.stvoPoints,
            },
        });
        const { response } = await call;

        if (response.law) {
            state.lawbookId = response.law.lawbookId;
        }

        emit('update:law', { id: id, law: response.law! });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSaveLaw = computed(() => state.lawbookId > 0 && availableLawBooks.value.some((book) => book.id === state.lawbookId));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!canSaveLaw.value) return;

    canSubmit.value = false;
    await saveLaw(props.law.id, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm class="my-2 flex flex-1 flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <div class="flex flex-1 flex-row gap-2">
            <UFormField class="text-sm font-medium">
                <UFieldGroup class="inline-flex w-full" orientation="vertical">
                    <UTooltip :text="$t('common.save')">
                        <UButton type="submit" variant="link" icon="i-mdi-content-save" :disabled="!canSaveLaw || !canSubmit" />
                    </UTooltip>

                    <UTooltip :text="$t('common.cancel')">
                        <UButton variant="link" icon="i-mdi-cancel" @click="$emit('close')" />
                    </UTooltip>
                </UFieldGroup>
            </UFormField>

            <UFormField class="flex-1 text-sm font-medium" :label="$t('common.law')" name="name">
                <UInput v-model="state.name" class="w-full" name="name" type="text" :placeholder="$t('common.law')" />
            </UFormField>
        </div>

        <UFormField class="flex-1 text-sm font-medium" :label="$t('common.law_book')" name="lawbookId">
            <ClientOnly>
                <USelectMenu
                    v-model="state.lawbookId"
                    class="w-full"
                    :items="availableLawBooks"
                    label-key="name"
                    value-key="id"
                    :search-input="{ placeholder: $t('common.search_field') }"
                    :disabled="availableLawBooks.length === 0"
                >
                    <template #default>
                        {{ currentLawBook?.name ?? $t('common.none_selected', [$t('common.law_book')]) }}
                    </template>

                    <template #item-label="{ item }">
                        {{ item.name }}
                    </template>
                </USelectMenu>
            </ClientOnly>
        </UFormField>

        <div class="flex flex-1 justify-between gap-2">
            <UFormField :label="$t('common.fine')" name="fine">
                <UInputNumber
                    v-model="state.fine"
                    name="fine"
                    :min="0"
                    :step="1000"
                    :format-options="{
                        style: 'currency',
                        currency: display.currencyName,
                        currencyDisplay: 'code',
                        currencySign: 'accounting',
                        maximumFractionDigits: 0,
                    }"
                    :placeholder="$t('common.fine')"
                />
            </UFormField>

            <UFormField :label="$t('common.detention_time')" name="detentionTime">
                <UInputNumber
                    v-model="state.detentionTime"
                    name="detentionTime"
                    :min="0"
                    :step="1"
                    :placeholder="$t('common.detention_time')"
                />
            </UFormField>

            <UFormField :label="$t('common.traffic_infraction_points')" name="stvoPoints">
                <UInputNumber
                    v-model="state.stvoPoints"
                    name="stvoPoints"
                    :min="0"
                    :step="1"
                    :placeholder="$t('common.traffic_infraction_points')"
                />
            </UFormField>
        </div>

        <UFormField :label="$t('common.description')" name="description">
            <UTextarea
                v-model="state.description"
                class="w-full"
                name="description"
                type="text"
                :placeholder="$t('common.description')"
            />
        </UFormField>

        <UFormField :label="$t('common.hint')" name="hint">
            <UTextarea v-model="state.hint" class="w-full" name="hint" type="text" :placeholder="$t('common.hint')" />
        </UFormField>
    </UForm>
</template>
