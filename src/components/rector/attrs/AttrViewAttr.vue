<script lang="ts" setup>
import { useCompletorStore } from '~/store/completor';
import { AttributeValues, Permission, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

const props = defineProps<{
    attribute: RoleAttribute;
    states: Map<string, AttributeValues | undefined>;
    disabled?: boolean;
    permission: Permission;
}>();

const emit = defineEmits<{
    (e: 'update:states', payload: Map<string, AttributeValues | undefined>): void;
    (e: 'changed'): void;
}>();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const jobGrades = ref<Map<string, JobGrade>>(new Map());

const states = ref<typeof props.states>(props.states);
const id = ref<string>(props.attribute.attrId);

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

const currentValue: AttributeValues = states.value.get(id.value)!;
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

    const map = currentValue.validValues.jobGradeList.jobs;
    if (checked && !map[job.name]) {
        map[job.name] = 1;
        jobGrades.value.set(job.name, job.grades[0]);
    } else if (!checked && map[job.name]) {
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
    jobGrades.value.set(job.name, job.grades[grade.grade - 1]);

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

            jobGrades.value.set(job.name, job.grades[(currentValue.validValues?.jobGradeList.jobs[job.name] ?? 1) - 1]);
        });
    }
});
</script>

<template>
    <div v-if="attribute">
        <UAccordion :items="[{ label: $t(`perms.${attribute.category}.${attribute.name}.attrs_types.${attribute.key}`) }]">
            <template #item>
                <div class="flex flex-col gap-2">
                    <div
                        v-if="
                            currentValue.validValues.oneofKind === 'stringList' &&
                            validValues?.validValues &&
                            validValues?.validValues.oneofKind === 'stringList'
                        "
                        class="flex flex-row flex-wrap gap-4"
                    >
                        <div
                            v-for="value in validValues.validValues.stringList.strings"
                            :key="value"
                            class="flex flex-initial flex-row flex-nowrap"
                        >
                            <UCheckbox
                                :name="value"
                                :model-value="!!currentValue.validValues.stringList.strings.find((v) => v === value)"
                                class="text-primary-500 focus:ring-primary-500 my-auto size-4 rounded border-base-300"
                                @click="toggleStringListValue(value)"
                            />
                            <span class="ml-1">{{
                                $t(`perms.${permission.category}.${permission.name}.attrs.${value.replaceAll('.', '_')}`)
                            }}</span>
                        </div>
                    </div>
                    <div
                        v-else-if="
                            currentValue.validValues.oneofKind === 'jobList' &&
                            validValues?.validValues &&
                            validValues?.validValues.oneofKind === 'jobList' &&
                            jobs !== undefined
                        "
                        class="flex flex-row flex-wrap gap-4"
                    >
                        <div v-for="job in jobs" :key="job.name" class="flex flex-initial flex-row flex-nowrap">
                            <UCheckbox
                                :name="job.name"
                                :model-value="!!currentValue.validValues.jobList?.strings.find((v) => v === job.name)"
                                class="text-primary-500 focus:ring-primary-500 my-auto size-4 rounded border-base-300"
                                @click="toggleJobListValue(job.name)"
                            />
                            <span class="ml-1">{{ job.label }}</span>
                        </div>
                    </div>
                    <div
                        v-else-if="
                            currentValue.validValues.oneofKind === 'jobGradeList' &&
                            validValues?.validValues &&
                            validValues.validValues.oneofKind === 'jobGradeList'
                        "
                        class="flex flex-col gap-2"
                    >
                        <div v-for="job in jobs" :key="job.name" class="flex flex-initial flex-row flex-nowrap gap-2">
                            <UCheckbox
                                :name="job.name"
                                :model-value="!!currentValue.validValues?.jobGradeList.jobs[job.name]"
                                class="text-primary-500 focus:ring-primary-500 my-auto size-4 rounded border-base-300"
                                @change="toggleJobGradeValue(job, $event)"
                            />
                            <span class="my-auto flex-1">{{ job.label }}</span>

                            <USelectMenu
                                @update:model-value="updateJobGradeValue(job, $event)"
                                nullable
                                :options="job.grades"
                                :search-attributes="['label']"
                                by="grade"
                                :placeholder="jobGrades.has(job.name) ? jobGrades.get(job.name)?.label : $t('common.na')"
                            >
                                <template #option="{ option: grade }">
                                    {{ grade?.label }}
                                </template>
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>
                                <template #empty> {{ $t('common.not_found', [$t('common.rank')]) }} </template>
                            </USelectMenu>
                        </div>
                    </div>
                    <div v-else>{{ currentValue.validValues.oneofKind }} {{ validValues }}</div>
                </div>
            </template>
        </UAccordion>
    </div>
</template>
