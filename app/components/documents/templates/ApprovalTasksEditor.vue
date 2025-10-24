<script lang="ts" setup>
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { ApprovalAssigneeKind } from '~~/gen/ts/resources/documents/approval';
import type { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    disabled?: boolean;
    signatureRequired?: boolean;
}>();

const tasks = defineModel<
    {
        ruleKind: ApprovalAssigneeKind;
        userId: number;
        user?: UserShort;
        job?: string;
        minimumGrade?: number;
        dueInDays?: number;
        label?: string;
        signatureRequired: boolean;
        slots: number;
        comment?: string;
    }[]
>({ required: true });

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

function addNewTask(): void {
    tasks.value.push({
        ruleKind: ApprovalAssigneeKind.JOB_GRADE,
        userId: 0,
        job: undefined,
        minimumGrade: undefined,
        label: '',
        signatureRequired: props.signatureRequired ?? false,
        slots: 1,
    });
}

function removeTask(idx: number): void {
    tasks.value.splice(idx, 1);
}

onBeforeMount(async () => listJobs());
</script>

<template>
    <div class="flex flex-col gap-1 divide-y divide-default md:divide-y-0">
        <div v-for="(task, idx) in tasks" :key="idx" class="flex flex-1 flex-col gap-1 pb-2 md:flex-row md:pb-0">
            <div class="flex flex-1 flex-col gap-1">
                <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
                    <UFormField
                        name="ruleKind"
                        class="min-w-40 flex-initial"
                        :label="$t('components.documents.approval.policy_form.rule_kind')"
                    >
                        <USelectMenu
                            v-model="task.ruleKind"
                            :items="[
                                { label: $t('common.user'), value: ApprovalAssigneeKind.USER },
                                { label: $t('common.job'), value: ApprovalAssigneeKind.JOB_GRADE },
                            ]"
                            value-key="value"
                            class="w-full"
                            :disabled="disabled"
                        />
                    </UFormField>

                    <UFormField
                        v-if="task.ruleKind === ApprovalAssigneeKind.USER"
                        name="ruleKind"
                        class="flex-1"
                        :label="$t('common.target')"
                    >
                        <SelectMenu
                            v-model="task.userId"
                            :searchable="
                                async (q: string) => {
                                    const users = await completorStore.completeCitizens({
                                        search: q,
                                        userIds: task.userId ? [task.userId] : [],
                                    });
                                    return users.filter((u) => u.userId !== activeChar?.userId);
                                }
                            "
                            searchable-key="completor-citizens"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.target')"
                            trailing
                            class="w-full"
                            value-key="userId"
                            :disabled="disabled"
                        >
                            <template v-if="task.user" #default>
                                {{ userToLabel(task.user) }}
                            </template>

                            <template #item-label="{ item }">
                                {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.citizen', 2)]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <template v-else-if="task.ruleKind === ApprovalAssigneeKind.JOB_GRADE">
                        <UFormField name="job" class="flex-1" :label="$t('common.job')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="task.job"
                                    :items="jobs ?? []"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['label', 'name']"
                                    value-key="name"
                                    class="w-full"
                                    :disabled="disabled"
                                >
                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.job')]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField name="minimumGrade" class="flex-1" :label="$t('common.job_grade')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="task.minimumGrade"
                                    :items="jobs.find((j) => j.name === task.job)?.grades ?? []"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    value-key="grade"
                                    class="w-full"
                                    :disabled="disabled"
                                >
                                    <template v-if="task.minimumGrade" #default>
                                        {{
                                            jobs
                                                .find((j) => j.name === task.job)
                                                ?.grades.find((g) => g.grade === task.minimumGrade)?.label
                                        }}
                                        ({{ task.minimumGrade }})
                                    </template>

                                    <template #item-label="{ item }"> {{ item.label }} ({{ item.grade }}) </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.job_grade')]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>
                    </template>

                    <UFormField name="dueInDays" class="flex-1" :label="$t('common.time_ago.day', 2)">
                        <UButtonGroup>
                            <UInputNumber v-model="task.dueInDays" class="w-full" :min="1" :max="30" :disabled="disabled" />
                            <UButton
                                icon="i-mdi-clear"
                                :disabled="disabled"
                                variant="outline"
                                @click="() => (task.dueInDays = undefined)"
                            />
                        </UButtonGroup>
                    </UFormField>
                </div>

                <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
                    <UFormField name="label" class="flex-1" :label="$t('common.label')">
                        <UInput v-model="task.label" class="w-full" :disabled="disabled" />
                    </UFormField>

                    <UFormField
                        name="signatureRequired"
                        class="h-full flex-initial"
                        :label="$t('components.documents.approval.signature_required')"
                    >
                        <USwitch v-model="task.signatureRequired" :disabled="signatureRequired || disabled" />
                    </UFormField>

                    <UFormField
                        v-if="task.ruleKind === ApprovalAssigneeKind.JOB_GRADE"
                        name="slots"
                        class="flex-initial"
                        :label="$t('components.documents.approval.slots')"
                    >
                        <UInputNumber
                            v-model="task.slots"
                            name="slots"
                            class="w-full"
                            :placeholder="$t('components.documents.approval.slots')"
                            :min="1"
                            :step="1"
                            :max="5"
                            :disabled="disabled"
                        />
                    </UFormField>
                </div>

                <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
                    <UFormField name="comment" class="flex-1" :label="$t('common.comment')">
                        <UInput v-model="task.comment" type="text" name="comment" class="w-full" :disabled="disabled" />
                    </UFormField>
                </div>
            </div>

            <UFormField class="md:mt-1" :ui="{ container: 'flex justify-end-safe md:inline' }">
                <UTooltip :text="$t('components.access.remove_entry')">
                    <UButton
                        color="red"
                        class="flex-initial"
                        icon="i-mdi-close"
                        :label="$t('components.access.remove_entry')"
                        :disabled="disabled"
                        :ui="{ label: 'md:hidden' }"
                        @click="() => removeTask(idx)"
                    />
                </UTooltip>
            </UFormField>
        </div>

        <div>
            <UTooltip :text="$t('common.add')">
                <UButton
                    class="w-full justify-center md:w-auto"
                    icon="i-mdi-plus"
                    :label="$t('common.add')"
                    :disabled="disabled"
                    :ui="{ label: 'md:hidden' }"
                    @click="addNewTask"
                />
            </UTooltip>
        </div>
    </div>
</template>
