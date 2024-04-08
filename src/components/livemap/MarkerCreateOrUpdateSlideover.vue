<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { useLivemapStore } from '~/store/livemap';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';
import DatePicker from '../partials/DatePicker.vue';
import { format } from 'date-fns';

const props = defineProps<{
    location?: Coordinate;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);
const { addOrpdateMarkerMarker } = livemapStore;

const markerTypes = ref<MarkerType[]>([MarkerType.CIRCLE, MarkerType.DOT, MarkerType.ICON]);

const defaultExpiresAt = ref<Date>(new Date());
defaultExpiresAt.value.setTime(defaultExpiresAt.value.getTime() + 1 * 60 * 60 * 1000);

const schema = z.object({
    name: z.string().min(3).max(255),
    description: z.string().min(6).max(512).optional(),
    expiresAt: z.date().optional(),
    color: z.string().length(7),
    markerType: z.number(),
    circleRadius: z.number().gte(5).lte(250),
    circleOpacity: z.number().gte(1).lte(75).optional(),
    icon: z.string().max(64).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    name: '',
    description: undefined,
    expiresAt: defaultExpiresAt.value,
    color: '#ee4b2b',
    markerType: MarkerType.CIRCLE,
    circleRadius: 50,
    circleOpacity: 15,
    icon: 'i-mdi-help',
});

async function createMarker(values: Schema): Promise<void> {
    const expiresAt = values.expiresAt ? toTimestamp(values.expiresAt) : undefined;

    try {
        const marker: MarkerMarker = {
            info: {
                id: '0',
                job: '',
                jobLabel: '',
                name: values.name,
                description: values.description,
                x: props.location ? props.location.x : storeLocation.value?.x ?? 0,
                y: props.location ? props.location.y : storeLocation.value?.y ?? 0,
                color: values.color,
            },
            expiresAt,
            type: values.markerType,
        };

        if (values.markerType === MarkerType.CIRCLE) {
            marker.data = {
                data: {
                    oneofKind: 'circle',
                    circle: {
                        radius: values.circleRadius,
                        oapcity: values.circleOpacity ?? 3,
                    },
                },
            };
        } else if (values.markerType === MarkerType.ICON) {
            marker.data = {
                data: {
                    oneofKind: 'icon',
                    icon: {
                        icon: values.icon ?? 'i-mdi-help',
                    },
                },
            };
        }

        const call = $grpc.getLivemapperClient().createOrUpdateMarker({
            marker,
        });
        const { response } = await call;

        if (response.marker !== undefined) {
            addOrpdateMarkerMarker(response.marker);
        }

        emit('close');
        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createMarker(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover :overlay="false">
        <UForm :schema="schema" :state="state" class="flex flex-1" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 min-h-[calc(100vh-(2*var(--header-height)))] max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
                        padding: 'px-1 py-2 sm:p-2',
                    },
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.livemap.create_marker.title') }}
                        </h3>

                        <UButton
                            color="gray"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            class="-my-1"
                            @click="
                                $emit('close');
                                isOpen = false;
                            "
                        />
                    </div>
                </template>

                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="name" class="block text-sm font-medium leading-6">
                                    {{ $t('common.name') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="name">
                                    <UInput
                                        v-model="state.name"
                                        type="text"
                                        name="name"
                                        :placeholder="$t('common.name')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="description" class="block text-sm font-medium leading-6">
                                    {{ $t('common.description') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="description">
                                    <UInput
                                        type="text"
                                        name="description"
                                        :placeholder="$t('common.description')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="expiresAt" class="block text-sm font-medium leading-6">
                                    {{ $t('common.expires_at') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="expiresAt">
                                    <UPopover :popper="{ placement: 'bottom-start' }">
                                        <UButton
                                            variant="outline"
                                            color="gray"
                                            block
                                            icon="i-mdi-calendar-month"
                                            :label="state.expiresAt ? format(state.expiresAt, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                                        />

                                        <template #panel="{ close }">
                                            <DatePicker v-model="state.expiresAt" @close="close" />
                                        </template>
                                    </UPopover>
                                </UFormGroup>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="color" class="block text-sm font-medium leading-6">
                                    {{ $t('common.color') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="color">
                                    <ColorInput v-model="state.color" disable-alpha format="hex" position="top" />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="markerType" class="block text-sm font-medium leading-6">
                                    {{ $t('common.marker') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="markerType">
                                    <USelectMenu
                                        v-model="state.markerType"
                                        name="markerType"
                                        :options="markerTypes"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    >
                                        <template #label>
                                            <span class="truncate">{{
                                                $t(`enums.livemap.MarkerType.${MarkerType[state.markerType ?? 0]}`)
                                            }}</span>
                                        </template>
                                        <template #option="{ option }">
                                            <span class="truncate">{{
                                                $t(`enums.livemap.MarkerType.${MarkerType[option ?? 0]}`)
                                            }}</span>
                                        </template>
                                    </USelectMenu>
                                </UFormGroup>
                            </dd>
                        </div>
                        <div
                            v-if="state.markerType === MarkerType.CIRCLE"
                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                        >
                            <dt class="text-sm font-medium leading-6">
                                <label for="circleRadius" class="block text-sm font-medium leading-6">
                                    {{ $t('common.radius') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="circleRadius">
                                    <UInput
                                        v-model="state.circleRadius"
                                        type="number"
                                        name="circleRadius"
                                        min="5"
                                        max="250"
                                        :placeholder="$t('common.radius')"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div
                            v-else-if="state.markerType === MarkerType.ICON"
                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                        >
                            <dt class="text-sm font-medium leading-6">
                                <label for="icon" class="block text-sm font-medium leading-6">
                                    {{ $t('common.icon') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="icon">
                                    <USelectMenu v-model="state.icon" :options="markerIcons">
                                        <template #label>
                                            <UIcon :name="state.icon" class="size-5" :style="{ color: state.color }" />
                                            <span class="truncate">{{ state.icon }}</span>
                                        </template>
                                        <template #option="{ option }">
                                            <UIcon :name="option" class="size-5" :style="{ color: state.color }" />
                                            <span class="truncate">{{ option }}</span>
                                        </template>
                                    </USelectMenu>
                                </UFormGroup>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.create') }}
                        </UButton>

                        <UButton
                            color="black"
                            block
                            class="flex-1"
                            @click="
                                $emit('close');
                                isOpen = false;
                            "
                        >
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </USlideover>
</template>
