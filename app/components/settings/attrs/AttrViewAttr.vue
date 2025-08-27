<script lang="ts" setup>
import { useCompletorStore } from '~/stores/completor';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import type { AttributeValues, RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';

const props = defineProps<{
    modelValue: RoleAttribute;
    disabled?: boolean;
    permission: Permission;
    defaultOpen?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: AttributeValues): void;
    (e: 'changed'): void;
    (e: 'opened', value: boolean): void;
}>();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const attribute = useVModel(props, 'modelValue', emit);

if (attribute.value?.validValues === undefined) {
    switch (lowercaseFirstLetter(attribute.value.type)) {
        case 'stringList': {
            attribute.value.validValues = {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [],
                    },
                },
            };
            break;
        }

        case 'jobList': {
            attribute.value.validValues = {
                validValues: {
                    oneofKind: 'jobList',
                    jobList: {
                        strings: [],
                    },
                },
            };
            break;
        }

        case 'jobGradeList': {
            attribute.value.validValues = {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: {
                        jobs: {},
                        fineGrained: false,
                        grades: {},
                    },
                },
            };
            break;
        }
    }
}

if (attribute.value?.maxValues === undefined || attribute.value?.maxValues.validValues.oneofKind === undefined) {
    switch (lowercaseFirstLetter(attribute.value.type)) {
        case 'stringList': {
            attribute.value.maxValues = {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [],
                    },
                },
            };
            break;
        }

        case 'jobList': {
            attribute.value.maxValues = {
                validValues: {
                    oneofKind: 'jobList',
                    jobList: {
                        strings: [],
                    },
                },
            };
            break;
        }

        case 'jobGradeList': {
            attribute.value.maxValues = {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: {
                        jobs: {},
                        fineGrained: false,
                        grades: {},
                    },
                },
            };
            break;
        }
    }
}

const attrValues = ref<AttributeValues>(attribute.value.maxValues!);

const validValues = computed<AttributeValues | undefined>(() => attribute.value.validValues);

watchOnce(attrValues, () => emit('changed'), { deep: true });

async function toggleStringListValue(value: string): Promise<void> {
    if (attrValues.value.validValues.oneofKind !== 'stringList') {
        return;
    }

    const idx = attrValues.value.validValues.stringList.strings.findIndex((v) => v === value);
    if (idx === -1) {
        attrValues.value.validValues.stringList.strings.push(value);
    } else {
        attrValues.value.validValues.stringList.strings.splice(idx, 1);
    }
}

async function toggleStringListAll(): Promise<void> {
    if (attrValues.value.validValues.oneofKind !== 'stringList') {
        return;
    }
    if (validValues.value?.validValues.oneofKind !== 'stringList') {
        return;
    }

    if (attrValues.value.validValues.stringList.strings.length === validValues.value?.validValues.stringList.strings.length) {
        attrValues.value.validValues.stringList.strings = [];
    } else {
        attrValues.value.validValues.stringList.strings = [...(validValues.value?.validValues.stringList.strings ?? [])];
    }
}

async function toggleJobListValue(value: string): Promise<void> {
    if (attrValues.value.validValues.oneofKind !== 'jobList') {
        return;
    }

    const idx = attrValues.value.validValues.jobList.strings.findIndex((v) => v === value);
    if (idx === -1) {
        attrValues.value.validValues.jobList.strings.push(value);
    } else {
        attrValues.value.validValues.jobList.strings.splice(idx, 1);
    }
}

async function toggleJobListAll(): Promise<void> {
    if (attrValues.value.validValues.oneofKind !== 'jobList') {
        return;
    }
    if (validValues.value?.validValues.oneofKind !== 'jobList') {
        return;
    }

    if (attrValues.value.validValues.jobList.strings.length === validValues.value?.validValues.jobList.strings.length) {
        attrValues.value.validValues.jobList.strings = [];
    } else {
        attrValues.value.validValues.jobList.strings = [...(validValues.value?.validValues.jobList.strings ?? [])];
    }
}

async function toggleJobGradeValue(job: Job, checked: boolean): Promise<void> {
    if (attrValues.value.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    attrValues.value.validValues.jobGradeList.jobs[job.name] = game.startJobGrade;
    if (!job.grades[0]) {
        return;
    }

    if (!attrValues.value.validValues.jobGradeList.fineGrained) {
        if (checked && !attrValues.value.validValues.jobGradeList.jobs[job.name]) {
            attrValues.value.validValues.jobGradeList.jobs[job.name] = job.grades[0].grade;
        } else if (!checked && attrValues.value.validValues.jobGradeList.jobs[job.name]) {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete attrValues.value.validValues.jobGradeList.jobs[job.name];
        }
    } else {
        if (checked && !attrValues.value.validValues.jobGradeList.grades[job.name]) {
            attrValues.value.validValues.jobGradeList.grades[job.name] = {
                grades: [job.grades[0].grade],
            };
        } else if (!checked && attrValues.value.validValues.jobGradeList.grades[job.name]) {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete attrValues.value.validValues.jobGradeList.grades[job.name];
        }
    }
}

onBeforeMount(async () => {
    if (attrValues.value.validValues.oneofKind === 'jobList' || attrValues.value.validValues.oneofKind === 'jobGradeList')
        await listJobs();
});

const { game } = useAppConfig();
</script>

<template>
    <div v-if="attribute">
        <UAccordion
            variant="outline"
            :items="[
                {
                    label: $t(`perms.${attribute.category}.${attribute.name}.attrs_types.${attribute.key}`),
                    disabled: defaultOpen,
                },
            ]"
            :unmount="true"
            :default-open="defaultOpen"
            :ui="{ default: { class: 'mb-0.5' } }"
        >
            <template #item>
                <div class="flex flex-col gap-2">
                    <div
                        v-if="
                            attrValues.validValues.oneofKind === 'stringList' &&
                            validValues?.validValues.oneofKind === 'stringList'
                        "
                        class="flex flex-col flex-wrap gap-2"
                    >
                        <span v-if="validValues.validValues.stringList.strings.length === 0">
                            {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                        </span>
                        <template v-else>
                            <div class="flex flex-row flex-wrap gap-2">
                                <div
                                    v-for="value in validValues.validValues.stringList.strings"
                                    :key="value"
                                    class="flex flex-initial flex-row flex-nowrap gap-1"
                                >
                                    <USwitch
                                        :name="value"
                                        :model-value="!!attrValues.validValues.stringList.strings.find((v) => v === value)"
                                        :disabled="disabled"
                                        @update:model-value="toggleStringListValue(value)"
                                    />
                                    <span>{{
                                        $t(
                                            `perms.${permission.category}.${permission.name}.attrs.${value.replaceAll('.', '_')}`,
                                        )
                                    }}</span>
                                </div>
                            </div>

                            <UButton
                                v-if="!disabled"
                                class="self-end"
                                size="xs"
                                color="neutral"
                                :icon="
                                    attrValues.validValues.stringList.strings.length !==
                                    validValues.validValues.stringList.strings.length
                                        ? 'i-mdi-check-all'
                                        : 'i-mdi-close'
                                "
                                :label="
                                    attrValues.validValues.stringList.strings.length !==
                                    validValues.validValues.stringList.strings.length
                                        ? $t('common.check_all')
                                        : $t('common.uncheck_all')
                                "
                                @click="toggleStringListAll()"
                            />
                        </template>
                    </div>
                    <div
                        v-else-if="
                            attrValues.validValues.oneofKind === 'jobList' && validValues?.validValues.oneofKind === 'jobList'
                        "
                        class="flex flex-col flex-wrap gap-2"
                    >
                        <div v-for="job in jobs" :key="job.name" class="flex flex-initial flex-row flex-nowrap gap-1">
                            <div class="flex flex-row flex-wrap gap-2">
                                <USwitch
                                    :name="job.name"
                                    :model-value="!!attrValues.validValues.jobList?.strings.find((v) => v === job.name)"
                                    :disabled="disabled"
                                    @update:model-value="toggleJobListValue(job.name)"
                                />
                                <span>{{ job.label }}</span>
                            </div>
                        </div>

                        <UButton
                            v-if="!disabled"
                            class="self-end"
                            size="xs"
                            color="neutral"
                            :icon="
                                attrValues.validValues.jobList.strings.length !== validValues.validValues.jobList.strings.length
                                    ? 'i-mdi-check-all'
                                    : 'i-mdi-close'
                            "
                            :label="
                                attrValues.validValues.jobList.strings.length !== validValues.validValues.jobList.strings.length
                                    ? $t('common.check_all')
                                    : $t('common.uncheck_all')
                            "
                            @click="toggleJobListAll()"
                        />
                    </div>

                    <div
                        v-else-if="
                            attrValues.validValues.oneofKind === 'jobGradeList' &&
                            validValues?.validValues.oneofKind === 'jobGradeList'
                        "
                        class="flex flex-col flex-wrap gap-2"
                    >
                        <div
                            v-for="job in jobs"
                            :key="job.name"
                            class="flex flex-initial flex-row flex-nowrap items-center gap-1"
                        >
                            <USwitch
                                :name="job.name"
                                :model-value="attrValues.validValues?.jobGradeList.jobs[job.name] !== undefined"
                                :disabled="disabled"
                                @update:model-value="toggleJobGradeValue(job, $event)"
                            />

                            <span class="flex-1">{{ job.label }}</span>

                            <ClientOnly>
                                <USelectMenu
                                    v-if="job.grades.length > 0"
                                    v-model="attrValues.validValues.jobGradeList.jobs[job.name]"
                                    class="flex-1"
                                    :disabled="disabled || attrValues.validValues?.jobGradeList.jobs[job.name] === undefined"
                                    :items="job.grades"
                                    :filter-fields="['label']"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :placeholder="$t('common.rank')"
                                    value-key="grade"
                                >
                                    <template #item-label>
                                        <template
                                            v-if="
                                                job.grades && attrValues.validValues.jobGradeList.jobs[job.name] !== undefined
                                            "
                                        >
                                            <span class="truncate text-highlighted"
                                                >{{
                                                    job.grades.find(
                                                        (g) =>
                                                            attrValues.validValues.oneofKind === 'jobGradeList' &&
                                                            g.grade ===
                                                                (attrValues.validValues.jobGradeList.jobs[job.name] ??
                                                                    game.startJobGrade),
                                                    )?.label ?? $t('common.na')
                                                }}
                                                ({{ attrValues.validValues.jobGradeList.jobs[job.name] }})</span
                                            >
                                        </template>
                                    </template>

                                    <template #item="{ item }">
                                        {{ item?.label
                                        }}<span v-if="item.grade >= game.startJobGrade"> ({{ item?.grade }})</span>
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.rank')]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </div>
                    </div>

                    <div v-else>{{ attrValues.validValues.oneofKind }} {{ validValues }}</div>
                </div>
            </template>
        </UAccordion>
    </div>
</template>
