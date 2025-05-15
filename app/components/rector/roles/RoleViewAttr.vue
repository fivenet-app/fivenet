<script lang="ts" setup>
import { useCompletorStore } from '~/stores/completor';
import type { AttributeValues, Permission, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import type { Job } from '~~/gen/ts/resources/users/jobs';

const props = defineProps<{
    modelValue: RoleAttribute;
    disabled?: boolean;
    permission: Permission;
    defaultOpen?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: AttributeValues): void;
    (e: 'changed'): void;
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
if (attribute.value?.value === undefined) {
    switch (lowercaseFirstLetter(attribute.value.type)) {
        case 'stringList': {
            attribute.value.value = {
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
            attribute.value.value = {
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
            attribute.value.value = {
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

const attrValue = ref<AttributeValues>(attribute.value.value!);

let maxValues = attribute.value.maxValues;
if (maxValues === undefined) {
    switch (lowercaseFirstLetter(attribute.value.type)) {
        case 'stringList':
            maxValues = {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [],
                    },
                },
            };
            break;

        case 'jobList':
            maxValues = {
                validValues: {
                    oneofKind: 'jobList',
                    jobList: {
                        strings: [],
                    },
                },
            };
            break;

        case 'jobGradeList':
            maxValues = {
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

const validValues = computed<AttributeValues | undefined>(() => attribute.value.validValues);

watchOnce(attrValue, () => emit('changed'), { deep: true });

async function toggleStringListValue(value: string): Promise<void> {
    if (attrValue.value.validValues.oneofKind !== 'stringList') {
        return;
    }

    const idx = attrValue.value.validValues.stringList.strings.findIndex((v) => v === value);
    if (idx === -1) {
        attrValue.value.validValues.stringList.strings.push(value);
    } else {
        attrValue.value.validValues.stringList.strings.splice(idx, 1);
    }
}

async function toggleJobListValue(value: string): Promise<void> {
    if (attrValue.value.validValues.oneofKind !== 'jobList') {
        return;
    }

    const idx = attrValue.value.validValues.jobList.strings.findIndex((v) => v === value);
    if (idx === -1) {
        attrValue.value.validValues.jobList.strings.push(value);
    } else {
        attrValue.value.validValues.jobList.strings.splice(idx, 1);
    }
}

async function toggleJobGradeValue(job: Job, checked: boolean): Promise<void> {
    if (attrValue.value.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    if (!job.grades[0]) {
        return;
    }

    if (!attrValue.value.validValues.jobGradeList.fineGrained) {
        if (checked && !attrValue.value.validValues.jobGradeList.jobs[job.name]) {
            attrValue.value.validValues.jobGradeList.jobs[job.name] = job.grades[0].grade;
        } else if (!checked && attrValue.value.validValues.jobGradeList.jobs[job.name]) {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete attrValue.value.validValues.jobGradeList.jobs[job.name];
        }
    } else {
        if (checked && !attrValue.value.validValues.jobGradeList.grades[job.name]) {
            attrValue.value.validValues.jobGradeList.grades[job.name] = {
                grades: [job.grades[0].grade],
            };
        } else if (!checked && attrValue.value.validValues.jobGradeList.grades[job.name]) {
            // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
            delete attrValue.value.validValues.jobGradeList.grades[job.name];
        }
    }
}

async function toggleJobGradeListFineGrained(checked: boolean): Promise<void> {
    if (attrValue.value.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    if (!attrValue.value.validValues.jobGradeList.grades) {
        attrValue.value.validValues.jobGradeList.grades = {};
    }

    attrValue.value.validValues.jobGradeList.fineGrained = checked;
}

onBeforeMount(async () => {
    if (attrValue.value.validValues.oneofKind === 'jobList' || attrValue.value.validValues.oneofKind === 'jobGradeList')
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
                            attrValue.validValues.oneofKind === 'stringList' &&
                            maxValues?.validValues.oneofKind === 'stringList'
                        "
                        class="flex flex-row flex-wrap gap-2"
                    >
                        <span v-if="maxValues.validValues.stringList.strings.length === 0">
                            {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                        </span>
                        <template v-else>
                            <div
                                v-for="value in maxValues.validValues.stringList.strings"
                                :key="value"
                                class="flex flex-initial flex-row flex-nowrap gap-1"
                            >
                                <UToggle
                                    :name="value"
                                    :model-value="!!attrValue.validValues.stringList.strings.find((v) => v === value)"
                                    :disabled="disabled"
                                    @update:model-value="toggleStringListValue(value)"
                                />
                                <span>{{
                                    $t(`perms.${permission.category}.${permission.name}.attrs.${value.replaceAll('.', '_')}`)
                                }}</span>
                            </div>
                        </template>
                    </div>
                    <div
                        v-else-if="
                            attrValue.validValues.oneofKind === 'jobList' && maxValues?.validValues.oneofKind === 'jobList'
                        "
                        class="flex flex-row flex-wrap gap-2"
                    >
                        <span v-if="maxValues.validValues.jobList.strings.length === 0">
                            {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                        </span>
                        <template v-else>
                            <div
                                v-for="job in jobs.filter(
                                    (j) =>
                                        maxValues?.validValues.oneofKind === 'jobList' &&
                                        (!maxValues?.validValues.jobList?.strings.length ||
                                            maxValues.validValues?.jobList?.strings.includes(j.name)),
                                )"
                                :key="job.name"
                                class="flex flex-initial flex-row flex-nowrap gap-1"
                            >
                                <UToggle
                                    :name="job.name"
                                    :model-value="!!attrValue.validValues.jobList?.strings.find((v) => v === job.name)"
                                    :disabled="disabled"
                                    @update:model-value="toggleJobListValue(job.name)"
                                />
                                <span>{{ job.label }}</span>
                            </div>
                        </template>
                    </div>

                    <div
                        v-else-if="
                            attrValue.validValues.oneofKind === 'jobGradeList' &&
                            maxValues?.validValues.oneofKind === 'jobGradeList'
                        "
                        class="flex flex-col flex-wrap gap-2"
                    >
                        <span v-if="Object.keys(maxValues.validValues.jobGradeList.jobs).length === 0">
                            {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                        </span>
                        <template v-else>
                            <div
                                v-for="job in jobs.filter(
                                    (j) =>
                                        maxValues &&
                                        maxValues.validValues.oneofKind === 'jobGradeList' &&
                                        maxValues.validValues.jobGradeList.jobs[j.name],
                                )"
                                :key="job.name"
                                class="flex flex-initial flex-row flex-nowrap items-center gap-1"
                            >
                                <UToggle
                                    v-if="!attrValue.validValues.jobGradeList.fineGrained"
                                    :name="job.name"
                                    :model-value="!!attrValue.validValues?.jobGradeList.jobs[job.name]"
                                    :disabled="disabled"
                                    @update:model-value="toggleJobGradeValue(job, $event)"
                                />
                                <UToggle
                                    v-else
                                    :name="job.name"
                                    :model-value="!!attrValue.validValues?.jobGradeList.grades[job.name]"
                                    :disabled="disabled"
                                    @update:model-value="toggleJobGradeValue(job, $event)"
                                />

                                <span class="flex-1">{{ job.label }}</span>

                                <ClientOnly>
                                    <USelectMenu
                                        v-if="!attrValue.validValues.jobGradeList.fineGrained"
                                        v-model="attrValue.validValues.jobGradeList.jobs[job.name]"
                                        class="flex-1"
                                        :disabled="disabled || !attrValue.validValues?.jobGradeList.jobs[job.name]"
                                        :options="
                                            job.grades.filter(
                                                (g) =>
                                                    maxValues &&
                                                    maxValues.validValues.oneofKind === 'jobGradeList' &&
                                                    (maxValues.validValues.jobGradeList.jobs[job.name] ?? game.startJobGrade) >=
                                                        g.grade,
                                            )
                                        "
                                        :search-attributes="['label']"
                                        :searchable-placeholder="$t('common.search_field')"
                                        :placeholder="$t('common.rank')"
                                        value-attribute="grade"
                                    >
                                        <template #label>
                                            <template v-if="job.grades && attrValue.validValues.jobGradeList.jobs[job.name]">
                                                <span class="truncate text-gray-900 dark:text-white"
                                                    >{{
                                                        job.grades.find(
                                                            (g) =>
                                                                attrValue.validValues.oneofKind === 'jobGradeList' &&
                                                                g.grade ===
                                                                    (attrValue.validValues.jobGradeList.jobs[job.name] ??
                                                                        game.startJobGrade),
                                                        )?.label ?? $t('common.na')
                                                    }}
                                                    ({{ attrValue.validValues.jobGradeList.jobs[job.name] }})</span
                                                >
                                            </template>
                                        </template>

                                        <template #option="{ option: grade }">
                                            {{ grade?.label
                                            }}<span v-if="grade.grade >= game.startJobGrade"> ({{ grade?.grade }})</span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.rank')]) }} </template>
                                    </USelectMenu>
                                    <USelectMenu
                                        v-else-if="attrValue.validValues.jobGradeList.grades[job.name]?.grades"
                                        v-model="attrValue.validValues.jobGradeList.grades[job.name]!.grades"
                                        class="flex-1"
                                        multiple
                                        :disabled="disabled || !attrValue.validValues?.jobGradeList.grades[job.name]"
                                        :options="
                                            job.grades.filter(
                                                (g) =>
                                                    maxValues &&
                                                    maxValues.validValues.oneofKind === 'jobGradeList' &&
                                                    maxValues.validValues.jobGradeList.jobs &&
                                                    (maxValues.validValues.jobGradeList.jobs[job.name] ?? game.startJobGrade) >=
                                                        g.grade,
                                            )
                                        "
                                        :search-attributes="['label']"
                                        :searchable-placeholder="$t('common.search_field')"
                                        :placeholder="$t('common.rank')"
                                        value-attribute="grade"
                                    >
                                        <template #label>
                                            {{
                                                $t(
                                                    'common.selected',
                                                    attrValue.validValues.jobGradeList.grades[job.name] === undefined
                                                        ? 0
                                                        : attrValue.validValues.jobGradeList.grades[job.name]!.grades.length,
                                                )
                                            }}
                                        </template>

                                        <template #option="{ option: grade }">
                                            {{ grade?.label
                                            }}<span v-if="grade.grade >= game.startJobGrade"> ({{ grade?.grade }})</span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.rank')]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </div>

                            <UDivider />

                            <div class="flex flex-row items-center gap-2">
                                <UToggle
                                    :model-value="attrValue.validValues.jobGradeList.fineGrained"
                                    :disabled="disabled"
                                    @update:model-value="toggleJobGradeListFineGrained($event)"
                                />

                                <UFormGroup
                                    :label="$t('components.rector.role_view.fine_grained_toggle.title')"
                                    :description="$t('components.rector.role_view.fine_grained_toggle.description')"
                                />
                            </div>
                        </template>
                    </div>

                    <div v-else>{{ attrValue.validValues.oneofKind }} {{ validValues }}</div>
                </div>
            </template>
        </UAccordion>
    </div>
</template>
