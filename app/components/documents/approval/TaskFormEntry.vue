<script lang="ts" setup>
import { z } from 'zod';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { ApprovalAssigneeKind, type ApprovalPolicy } from '~~/gen/ts/resources/documents/approval';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { UserShort } from '~~/gen/ts/resources/users/users';

withDefaults(
    defineProps<{
        policy?: ApprovalPolicy;
        jobs: Job[];
        hideJobs?: string[];
    }>(),
    {
        policy: undefined,
        hideJobs: () => [],
    },
);

const task = defineModel<Task>({ required: true });

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const { game } = useAppConfig();

const completorStore = useCompletorStore();

const _schema = z.union([
    z.object({
        ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.USER),
        userId: z.coerce.number(),
        user: z.custom<UserShort>().optional(),
        job: z.coerce.string().optional(),
        minimumGrade: z.coerce.number().min(game.startJobGrade).optional(),
        label: z.string().max(120).default(''),
        signatureRequired: z.coerce.boolean().default(false),
        slots: z.coerce.number().min(1).max(5).optional().default(1),
        dueAt: z.date().optional(),
        comment: z.coerce.string().max(255).optional(),
    }),
    z.object({
        ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.JOB_GRADE),
        userId: z.coerce.number().optional().default(0),
        user: z.custom<UserShort>().optional(),
        job: z.coerce.string().optional(),
        minimumGrade: z.coerce.number().min(game.startJobGrade).optional(),
        label: z.string().max(120).default(''),
        signatureRequired: z.coerce.boolean().default(false),
        slots: z.coerce.number().min(1).max(5).optional().default(1),
        dueAt: z.date().optional(),
        comment: z.coerce.string().max(255).optional(),
    }),
]);

export type Task = z.output<typeof _schema>;

watch(
    task,
    () => {
        if (task.value.ruleKind === ApprovalAssigneeKind.USER) {
            task.value.job = undefined;
            task.value.minimumGrade = undefined;
            task.value.slots = 1;

            task.value.userId = task.value.user?.userId ?? task.value.userId;
        } else if (task.value.ruleKind === ApprovalAssigneeKind.JOB_GRADE) {
            task.value.userId = 0;
            task.value.user = undefined;
        }
    },
    { deep: true },
);
</script>

<template>
    <div class="flex flex-col gap-1">
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
                />
            </UFormField>

            <UFormField
                v-if="task.ruleKind === ApprovalAssigneeKind.USER"
                name="ruleKind"
                class="flex-1"
                :label="$t('common.target')"
            >
                <SelectMenu
                    v-model="task.user"
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
                >
                    <template v-if="task.user" #default>
                        {{ userToLabel(task.user) }}
                    </template>

                    <template #item-label="{ item }">
                        {{ `${item?.firstname} ${item?.lastname} (${item?.dateofbirth})` }}
                    </template>

                    <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                </SelectMenu>
            </UFormField>

            <template v-else-if="task.ruleKind === ApprovalAssigneeKind.JOB_GRADE">
                <UFormField name="job" class="flex-1" :label="$t('common.job')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="task.job"
                            :items="jobs?.filter((j) => hideJobs.length === 0 || !hideJobs.includes(j.name)) ?? []"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['label', 'name']"
                            value-key="name"
                            class="w-full"
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
                        >
                            <template v-if="task.minimumGrade" #default>
                                {{
                                    jobs.find((j) => j.name === task.job)?.grades.find((g) => g.grade === task.minimumGrade)
                                        ?.label
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

            <UFormField name="dueAt" class="flex-1" :label="$t('common.due_at')">
                <InputDatePicker v-model="task.dueAt" class="w-full" />
            </UFormField>
        </div>

        <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
            <UFormField name="label" class="flex-1" :label="$t('common.name')">
                <UInput v-model="task.label" class="w-full" />
            </UFormField>

            <UFormField
                name="signatureRequired"
                class="h-full flex-initial"
                :label="$t('components.documents.approval.signature_required')"
            >
                <div class="flex flex-1 items-center justify-center">
                    <USwitch v-model="task.signatureRequired" :disabled="policy?.signatureRequired" />
                </div>
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
                    :max="5"
                />
            </UFormField>
        </div>

        <div class="grid grid-cols-2 gap-2 md:flex md:flex-1">
            <UFormField name="comment" class="flex-1" :label="$t('common.comment')">
                <UInput v-model="task.comment" type="text" name="comment" class="w-full" />
            </UFormField>
        </div>
    </div>
</template>
