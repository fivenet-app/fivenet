<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { HelpIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { useLivemapStore } from '~/stores/livemap';
import type { Coordinate } from '~/types/livemap';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/marker_marker';

const props = defineProps<{
    location?: Coordinate;
    marker?: MarkerMarker;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);
const { addOrUpdateMarkerMarker } = livemapStore;

const markerTypes = [{ type: MarkerType.CIRCLE }, { type: MarkerType.DOT }, { type: MarkerType.ICON }];

const defaultExpiresAt = ref<Date>(new Date());
defaultExpiresAt.value.setTime(defaultExpiresAt.value.getTime() + 1 * 60 * 60 * 1000);

const schema = z.object({
    name: z.string().min(1).max(255),
    description: z.union([z.string().min(3).max(1024), z.string().length(0).optional()]),
    expiresAt: z.date().optional(),
    color: z.string().length(7),
    markerType: z.nativeEnum(MarkerType),
    circleRadius: z.coerce.number().gte(5).lte(250),
    circleOpacity: z.coerce.number().gte(1).lte(75).optional(),
    icon: z.string().max(128).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: props.marker?.name ?? '',
    description: props.marker?.description,
    expiresAt: props.marker?.expiresAt ? toDate(props.marker?.expiresAt) : defaultExpiresAt.value,
    color: props.marker?.color ?? '#ee4b2b',
    markerType: props.marker?.type ?? MarkerType.CIRCLE,
    circleRadius:
        props.marker?.data?.data.oneofKind === 'circle' && props.marker?.data?.data.circle.radius
            ? props.marker?.data?.data.circle.radius
            : 50,
    circleOpacity:
        props.marker?.data?.data.oneofKind === 'circle' && props.marker?.data?.data.circle.opacity
            ? props.marker?.data?.data.circle.opacity
            : 15,
    icon:
        props.marker?.data?.data.oneofKind === 'icon' && props.marker?.data?.data.icon.icon
            ? props.marker?.data?.data.icon.icon
            : HelpIcon.name,
});

async function createOrUpdateMarker(values: Schema): Promise<void> {
    const expiresAt = values.expiresAt ? toTimestamp(values.expiresAt) : undefined;

    try {
        const marker: MarkerMarker = {
            id: props.marker?.id ?? 0,
            job: '',
            jobLabel: '',
            name: values.name,
            description: values.description,
            x: props.marker?.x ?? props.location?.x ?? storeLocation.value?.x ?? 0,
            y: props.marker?.y ?? props.location?.y ?? storeLocation.value?.y ?? 0,
            color: values.color,
            expiresAt: expiresAt,
            type: values.markerType,
        };

        if (values.markerType === MarkerType.CIRCLE) {
            marker.data = {
                data: {
                    oneofKind: 'circle',
                    circle: {
                        radius: values.circleRadius,
                        opacity: values.circleOpacity ?? 3,
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

        const call = $grpc.livemap.livemap.createOrUpdateMarker({
            marker,
        });
        const { response } = await call;

        if (response.marker !== undefined) {
            addOrUpdateMarkerMarker(response.marker);
        }

        emit('close');
        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateMarker(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover :ui="{ width: 'w-screen max-w-xl' }" :overlay="false">
        <UForm class="flex flex-1" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                        padding: 'px-1 py-2 sm:p-2',
                    },
                    ring: '',
                    divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{
                                !marker
                                    ? $t('components.livemap.create_marker.title')
                                    : $t('components.livemap.update_marker.title')
                            }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="gray"
                            variant="ghost"
                            icon="i-mdi-window-close"
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
                                <label class="block text-sm font-medium leading-6" for="name">
                                    {{ $t('common.name') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="name">
                                    <UInput v-model="state.name" type="text" name="name" :placeholder="$t('common.name')" />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label class="block text-sm font-medium leading-6" for="description">
                                    {{ $t('common.description') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="description">
                                    <UInput
                                        v-model="state.description"
                                        type="text"
                                        name="description"
                                        :placeholder="$t('common.description')"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label class="block text-sm font-medium leading-6" for="expiresAt">
                                    {{ $t('common.expires_at') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="expiresAt">
                                    <DatePickerPopoverClient
                                        v-model="state.expiresAt"
                                        date-format="dd.MM.yyyy HH:mm"
                                        :popover="{ popper: { placement: 'bottom-start' } }"
                                        :date-picker="{ mode: 'dateTime', is24hr: true, clearable: true }"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label class="block text-sm font-medium leading-6" for="color">
                                    {{ $t('common.color') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="color">
                                    <ColorPickerClient v-model="state.color" />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label class="block text-sm font-medium leading-6" for="markerType">
                                    {{ $t('common.marker') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="markerType">
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.markerType"
                                            name="markerType"
                                            :options="markerTypes"
                                            value-attribute="type"
                                            :searchable-placeholder="$t('common.search_field')"
                                        >
                                            <template #label>
                                                <span class="truncate">{{
                                                    $t(`enums.livemap.MarkerType.${MarkerType[state.markerType ?? 0]}`)
                                                }}</span>
                                            </template>

                                            <template #option="{ option }">
                                                <span class="truncate">{{
                                                    $t(`enums.livemap.MarkerType.${MarkerType[option.type ?? 0]}`)
                                                }}</span>
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>
                            </dd>
                        </div>
                        <div
                            v-if="state.markerType === MarkerType.CIRCLE"
                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                        >
                            <dt class="text-sm font-medium leading-6">
                                <label class="block text-sm font-medium leading-6" for="circleRadius">
                                    {{ $t('common.radius') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="circleRadius">
                                    <UInput
                                        v-model="state.circleRadius"
                                        type="number"
                                        name="circleRadius"
                                        :min="5"
                                        :max="250"
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
                                <label class="block text-sm font-medium leading-6" for="icon">
                                    {{ $t('common.icon') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="icon">
                                    <IconSelectMenu v-model="state.icon" :color="state.color" />
                                </UFormGroup>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ !marker ? $t('common.create') : $t('common.save') }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            color="black"
                            block
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
