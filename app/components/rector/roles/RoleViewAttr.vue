<script lang="ts" setup>
import { useCompletorStore } from '~/store/completor';
import type { AttributeValues, Permission, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

const props = defineProps<{
    attribute: RoleAttribute;
    states: Map<number, AttributeValues | undefined>;
    disabled?: boolean;
    permission: Permission;
}>();

const emit = defineEmits<{
    (e: 'update:states', payload: Map<number, AttributeValues | undefined>): void;
    (e: 'changed'): void;
}>();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const jobGrades = ref(new Map<string, JobGrade>());

const states = ref<typeof props.states>(props.states);
const id = ref<number>(props.attribute.attrId);

let maxValues = props.attribute.maxValues;
if (maxValues === undefined) {
    switch (lowercaseFirstLetter(props.attribute.type)) {
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
                    },
                },
            };
            break;
    }
}

if (!states.value.has(id.value) || states.value.get(id.value) === undefined) {
    switch (lowercaseFirstLetter(props.attribute.type)) {
        case 'stringList': {
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [],
                    },
                },
            });
            break;
        }

        case 'jobList': {
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'jobList',
                    jobList: {
                        strings: [],
                    },
                },
            });
            break;
        }

        case 'jobGradeList': {
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: {
                        jobs: {},
                    },
                },
            });
            break;
        }
    }
}

const currentValue = states.value.get(id.value)!;
const validValues = ref<AttributeValues | undefined>(props.attribute.validValues);

async function toggleStringListValue(value: string): Promise<void> {
    if (currentValue.validValues.oneofKind !== 'stringList') {
        return;
    }

    const array = currentValue.validValues.stringList.strings;
    if (!array.includes(value)) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    currentValue.validValues.stringList.strings = array;
    states.value.set(id.value, currentValue);
    emit('update:states', states.value);
    emit('changed');
}

async function toggleJobListValue(value: string): Promise<void> {
    if (currentValue.validValues.oneofKind !== 'jobList') {
        return;
    }

    const array = currentValue.validValues.jobList.strings;
    if (!array.includes(value)) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    currentValue.validValues.jobList.strings = array;
    states.value.set(id.value, currentValue);
    emit('update:states', states.value);
    emit('changed');
}

async function toggleJobGradeValue(job: Job, checked: boolean): Promise<void> {
    if (currentValue.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    if (!job.grades[0]) {
        return;
    }

    const map = currentValue.validValues.jobGradeList.jobs;
    if (checked && !map[job.name]) {
        map[job.name] = 1;
        jobGrades.value.set(job.name, job.grades[0]);
    } else if (!checked && map[job.name]) {
        // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
        delete map[job.name];
        jobGrades.value.set(job.name, job.grades[0]);
    } else {
        return;
    }

    currentValue.validValues.jobGradeList.jobs = map;
    states.value.set(id.value, currentValue);
    emit('update:states', states.value);
    emit('changed');
}

async function updateJobGradeValue(job: Job, grade: JobGrade): Promise<void> {
    if (currentValue.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    const map = currentValue.validValues.jobGradeList.jobs;

    map[job.name] = grade.grade;
    if (!job.grades[grade.grade - 1]) {
        return;
    }
    jobGrades.value.set(job.name, job.grades[grade.grade - 1]!);

    currentValue.validValues.jobGradeList.jobs = map;
    states.value.set(id.value, currentValue);
    emit('update:states', states.value);
    emit('changed');
}

onBeforeMount(async () => {
    if (currentValue.validValues.oneofKind === 'jobList' || currentValue.validValues.oneofKind === 'jobGradeList') {
        await listJobs();

        jobs.value.forEach((job) => {
            if (currentValue.validValues.oneofKind !== 'jobGradeList') {
                return;
            }

            if (maxValues && maxValues.validValues.oneofKind === 'jobGradeList') {
                if (!maxValues.validValues.jobGradeList.jobs[job.name]) {
                    return;
                }
            }

            const grade = job.grades[(currentValue.validValues?.jobGradeList.jobs[job.name] ?? 1) - 1];
            if (!grade) {
                return;
            }
            jobGrades.value.set(job.name, grade);
        });
    }
});

const { game } = useAppConfig();
</script>

<template>
    <div v-if="attribute">
        <UAccordion
            variant="outline"
            :items="[{ label: $t(`perms.${attribute.category}.${attribute.name}.attrs_types.${attribute.key}`) }]"
            :unmount="true"
        >
            <template #item>
                <div class="flex flex-col gap-2">
                    <div
                        v-if="
                            currentValue.validValues.oneofKind === 'stringList' &&
                            maxValues?.validValues &&
                            maxValues?.validValues.oneofKind === 'stringList'
                        "
                        class="flex flex-row flex-wrap gap-3"
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
                                    :model-value="!!currentValue.validValues.stringList.strings.find((v) => v === value)"
                                    @click="toggleStringListValue(value)"
                                />
                                <span>{{
                                    $t(`perms.${permission.category}.${permission.name}.attrs.${value.replaceAll('.', '_')}`)
                                }}</span>
                            </div>
                        </template>
                    </div>
                    <div
                        v-else-if="
                            currentValue.validValues.oneofKind === 'jobList' &&
                            maxValues?.validValues &&
                            maxValues?.validValues.oneofKind === 'jobList'
                        "
                        class="flex flex-row flex-wrap gap-3"
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
                                    :model-value="!!currentValue.validValues.jobList?.strings.find((v) => v === job.name)"
                                    @click="toggleJobListValue(job.name)"
                                />
                                <span>{{ job.label }}</span>
                            </div>
                        </template>
                    </div>
                    <div
                        v-else-if="
                            currentValue.validValues.oneofKind === 'jobGradeList' &&
                            maxValues?.validValues &&
                            maxValues.validValues.oneofKind === 'jobGradeList'
                        "
                        class="flex flex-col flex-wrap gap-3"
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
                                <UCheckbox
                                    :name="job.name"
                                    :model-value="!!currentValue.validValues?.jobGradeList.jobs[job.name]"
                                    @change="toggleJobGradeValue(job, ($event.target as any)?.checked)"
                                />

                                <span class="flex-1">{{ job.label }}</span>

                                <ClientOnly>
                                    <USelectMenu
                                        class="flex-1"
                                        :disabled="!currentValue.validValues?.jobGradeList.jobs[job.name]"
                                        :options="
                                            job.grades.filter(
                                                (g) =>
                                                    maxValues &&
                                                    maxValues.validValues.oneofKind === 'jobGradeList' &&
                                                    (maxValues.validValues.jobGradeList.jobs[job.name] ?? game.startJobGrade) +
                                                        1 >
                                                        g.grade,
                                            )
                                        "
                                        :search-attributes="['label']"
                                        by="grade"
                                        :placeholder="$t('common.rank')"
                                        :searchable-placeholder="$t('common.search_field')"
                                        @update:model-value="updateJobGradeValue(job, $event)"
                                    >
                                        <template #label>
                                            <template v-if="job.grades && currentValue.validValues.jobGradeList.jobs[job.name]">
                                                <span class="truncate">{{
                                                    job.grades[
                                                        (currentValue.validValues.jobGradeList.jobs[job.name] ??
                                                            game.startJobGrade) - 1
                                                    ]?.label ?? $t('common.na')
                                                }}</span>
                                            </template>
                                        </template>

                                        <template #option="{ option: grade }">
                                            {{ grade?.label }}
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.rank')]) }} </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </div>
                        </template>
                    </div>
                    <div v-else>{{ currentValue.validValues.oneofKind }} {{ validValues }}</div>
                </div>
            </template>
        </UAccordion>
    </div>
</template>
