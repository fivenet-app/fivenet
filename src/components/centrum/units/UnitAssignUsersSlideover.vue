<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCompletorStore } from '~/store/completor';
import { Unit } from '~~/gen/ts/resources/centrum/units';
import { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    unit: Unit;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const completorStore = useCompletorStore();

const entriesCitizens = ref<UserShort[]>([]);
const selectedCitizens = ref<UserShort[]>(props.unit.users.filter((u) => u !== undefined).map((u) => u.user!));
const queryCitizens = ref('');

async function assignUnit(): Promise<void> {
    try {
        const toAdd: number[] = [];
        const toRemove: number[] = [];
        selectedCitizens.value?.forEach((u) => {
            toAdd.push(u.userId);
        });
        props.unit.users?.forEach((u) => {
            const idx = selectedCitizens.value.findIndex((su) => su.userId === u.userId);
            if (idx === -1) {
                toRemove.push(u.userId);
            }
        });

        const call = $grpc.getCentrumClient().assignUnit({
            unitId: props.unit.id,
            toAdd,
            toRemove,
        });
        await call;

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function findCitizens(): Promise<void> {
    entriesCitizens.value = await completorStore.completeCitizens({
        search: queryCitizens.value,
        currentJob: true,
        onDuty: true,
    });
    entriesCitizens.value.unshift(...selectedCitizens.value);
}

function charsGetDisplayValue(chars: UserShort[]): string {
    const cs: string[] = [];
    chars.forEach((c) => cs.push(`${c?.firstname} ${c?.lastname}`));

    return cs.join(', ');
}

watchDebounced(queryCitizens, async () => await findCitizens(), {
    debounce: 500,
    maxWait: 1250,
});

onMounted(async () => {
    findCitizens();
});

const canSubmit = ref(true);

const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await assignUnit().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.assign_unit.title') }}: {{ unit.name }} ({{ unit.initials }})
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <div class="flex flex-1 flex-col justify-between">
                    <div class="divide-y divide-gray-200 px-2 sm:px-6">
                        <div class="mt-1">
                            <div class="my-2 space-y-24">
                                <div class="flex-1">
                                    <Combobox v-model="selectedCitizens" as="div" multiple nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (chars: any) => (chars ? charsGetDisplayValue(chars) : $t('common.na'))
                                                    "
                                                    :placeholder="$t('common.citizen', 2)"
                                                    @change="queryCitizens = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="entriesCitizens.length > 0"
                                                class="absolute z-30 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="user in entriesCitizens"
                                                    v-slot="{ active, selected }"
                                                    :key="user?.userId"
                                                    :value="user"
                                                    as="char"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ user?.firstname }}
                                                            {{ user?.lastname }}
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="size-5" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>

                                    <div class="mt-4 overflow-hidden rounded-md bg-base-800">
                                        <ul role="list" class="divide-y divide-gray-200 text-sm font-medium text-gray-100">
                                            <li
                                                v-for="user in selectedCitizens"
                                                :key="user.userId"
                                                class="inline-flex items-center px-6 py-4"
                                            >
                                                <CitizenInfoPopover :user="user" />
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <div>
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton :disabled="!canSubmit" :loading="!canSubmit" @click="onSubmitThrottle">
                        {{ $t('common.update') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </USlideover>
</template>
