<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue'
import { RoleAttribute } from '@fivenet/gen/resources/permissions/permissions_pb';
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/solid';
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';

const props = defineProps<{
    attribute: RoleAttribute,
    states: Map<number, (string | number)[]>,
    jobs: Job[],
    disabled?: boolean,
}>();

const emit = defineEmits<{
    (e: 'update:states', payload: Map<number, (string | number)[]>): void,
}>();

const states = ref<typeof props.states>(props.states);
const id = ref<number>(props.attribute.getAttrId());

const validValues = ref<string[] | undefined>(props.attribute.getValidValues() ? JSON.parse(props.attribute.getValidValues()) as string[] : undefined);
const type = ref<string | undefined>(props.attribute.getType());

function getState(): (string | number)[] {
    if (!states.value.has(id.value)) states.value.set(id.value, []);

    return states.value.get(id.value)!;
}

async function toggleListValue(value: string | number): Promise<void> {
    const state = getState();

    if (state.indexOf(value) < 0) {
        state.push(value);
    } else {
        state.splice(state.indexOf(value), 1);
    }

    states.value.set(id.value, state);
    emit('update:states', states.value);
}
</script>

<style scoped>
.upsidedown {
    transform: rotate(180deg);
}
</style>

<template>
    <div v-if="$props.attribute">
        <Disclosure as="div"
            :class="[$props.disabled ? 'border-neutral/10 text-base-300' : 'hover:border-neutral/70 border-neutral/20 text-neutral']"
            v-slot="{ open }">
            <DisclosureButton :disabled="$props.disabled"
                :class="[open ? 'rounded-t-lg border-b-0' : 'rounded-lg', $props.disabled ? 'cursor-not-allowed' : '', ' flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors']">
                <span class="text-base leading-7 transition-colors">
                    Options
                </span>
                <span class="ml-6 flex h-7 items-center">
                    <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                        aria-hidden="true" />
                </span>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2">
                <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                    <div v-if="type === 'StringList' && validValues" class="flex flex-row gap-4 flex-wrap">
                        <div v-for="value in validValues" :key="value" class="flex flex-row flex-initial flex-nowrap">
                            <input :id="value" :name="value" type="checkbox" :checked="!!getState().find(v => v === value)"
                                @click="toggleListValue(value)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="ml-1">{{ value }}</span>
                        </div>
                    </div>
                    <div v-else-if="type === 'JobList'" class="flex flex-row gap-4 flex-wrap">
                        <div v-for="job in props.jobs" :key="job.getName()" class="flex flex-row flex-initial flex-nowrap">
                            <input :id="job.getName()" :name="job.getName()" type="checkbox"
                                :checked="!!getState().find(v => v === job.getName())"
                                @click="toggleListValue(job.getName())"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="ml-1">{{ job.getLabel() }}</span>
                        </div>
                    </div>
                    <div v-else-if="type === 'JobGradeList'">
                        JobGradeList
                    </div>
                    <div v-else>{{ type }} {{ validValues }}</div>
                </div>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
