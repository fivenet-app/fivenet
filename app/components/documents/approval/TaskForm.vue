<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';

const props = defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const approvalClient = await getDocumentsApprovalClient();

const schema = z.object({
    tasks: z.object({}).array().min(1),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    tasks: [],
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await upsertPolicy(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function upsertPolicy(values: Schema): Promise<void> {
    const call = approvalClient.startApprovalRound({
        documentId: props.documentId,
    });
    await call;

    // TODO
}

function addNewTask(): void {
    state.tasks.push({});
}

function removeTask(idx: number): void {
    state.tasks.splice(idx, 1);
}

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDrawer
        :title="$t('common.approve')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <template #title>
            <div class="inline-flex flex-1 items-center gap-1">
                <span>{{ $t('common.approve') }}</span>
                <UIcon name="i-mdi-slash-forward" />
                <span>{{ $t('common.policy') }}</span>
                <UIcon name="i-mdi-slash-forward" />
                <span>{{ $t('components.documents.approval.task_form.title') }}</span>
            </div>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <UCard :ui="{ body: 'p-4 sm:p-4', footer: 'p-4 sm:px-4' }">
                    <template #header>
                        <div class="flex items-center justify-between gap-1">
                            <h3 class="flex-1 text-base leading-6 font-semibold">
                                {{ $t('components.documents.approval.task_form.title') }}
                            </h3>
                        </div>
                    </template>

                    <UForm ref="formRef" :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
                        <div class="flex flex-col gap-1 divide-y divide-default md:divide-y-0">
                            <div
                                v-for="(task, idx) in state.tasks"
                                :key="idx"
                                class="flex flex-1 flex-col gap-1 pb-2 md:flex-row md:pb-0"
                            >
                                <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
                                    <UFormField
                                        name="ruleKind"
                                        :label="$t('components.documents.approval.policy_form.rule_kind')"
                                    >
                                        TODO
                                        {{ task }}
                                    </UFormField>
                                </div>

                                <UFormField class="md:mt-1" :ui="{ container: 'flex justify-end-safe md:inline' }">
                                    <UTooltip :text="$t('components.access.remove_entry')">
                                        <UButton
                                            color="red"
                                            class="flex-initial"
                                            icon="i-mdi-close"
                                            :label="$t('components.access.remove_entry')"
                                            :ui="{ label: 'md:hidden' }"
                                            @click="() => removeTask(idx)"
                                        />
                                    </UTooltip>
                                </UFormField>
                            </div>

                            <div>
                                <UTooltip :text="$t('components.access.add_entry')">
                                    <UButton
                                        class="w-full justify-center md:w-auto"
                                        icon="i-mdi-plus"
                                        :label="$t('components.access.add_entry')"
                                        :ui="{ label: 'md:hidden' }"
                                        @click="addNewTask"
                                    />
                                </UTooltip>
                            </div>
                        </div>
                    </UForm>

                    <template #footer>
                        <UButtonGroup class="w-full flex-1">
                            <UButton
                                type="submit"
                                :disabled="!canSubmit"
                                block
                                class="w-full"
                                :label="$t('common.submit')"
                                trailing-icon="i-mdi-arrow-forward"
                                @click="formRef?.submit()"
                            />
                        </UButtonGroup>
                    </template>
                </UCard>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col">
                <UButtonGroup class="w-full flex-1">
                    <UButton
                        color="neutral"
                        variant="subtle"
                        icon="i-mdi-arrow-back"
                        block
                        :label="$t('common.back')"
                        @click="() => $emit('close', false)"
                    />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
