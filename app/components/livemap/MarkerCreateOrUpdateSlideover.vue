<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { HelpIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { useLivemapStore } from '~/stores/livemap';
import type { Coordinate } from '~/types/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/marker_marker';
import InputDatePicker from '../partials/InputDatePicker.vue';

const props = defineProps<{
    location?: Coordinate;
    marker?: MarkerMarker;
}>();

const emit = defineEmits<{
    close: [boolean];
}>();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);
const { addOrUpdateMarkerMarker } = livemapStore;

const livemapLivemapClient = await getLivemapLivemapClient();

const markerTypes = [{ value: MarkerType.CIRCLE }, { value: MarkerType.DOT }, { value: MarkerType.ICON }];

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

        const call = livemapLivemapClient.createOrUpdateMarker({
            marker,
        });
        const { response } = await call;

        if (response.marker !== undefined) {
            addOrUpdateMarkerMarker(response.marker);
        }

        emit('close', false);
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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <USlideover
        :title="!marker ? $t('components.livemap.create_marker.title') : $t('components.livemap.update_marker.title')"
        :overlay="false"
    >
        <template #body>
            <UForm ref="formRef" class="flex flex-1" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="name">
                                {{ $t('common.name') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="name">
                                <UInput v-model="state.name" type="text" name="name" :placeholder="$t('common.name')" />
                            </UFormField>
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="description">
                                {{ $t('common.description') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="description">
                                <UInput
                                    v-model="state.description"
                                    type="text"
                                    name="description"
                                    :placeholder="$t('common.description')"
                                />
                            </UFormField>
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="expiresAt">
                                {{ $t('common.expires_at') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="expiresAt">
                                <InputDatePicker v-model="state.expiresAt" date-format="dd.MM.yyyy HH:mm" clearable time />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="color">
                                {{ $t('common.color') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="color">
                                <ColorPicker v-model="state.color" />
                            </UFormField>
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="markerType">
                                {{ $t('common.marker') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="markerType">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.markerType"
                                        name="markerType"
                                        :items="markerTypes"
                                        value-key="value"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #default>
                                            {{ $t(`enums.livemap.MarkerType.${MarkerType[state.markerType ?? 0]}`) }}
                                        </template>

                                        <template #item="{ item }">
                                            {{ $t(`enums.livemap.MarkerType.${MarkerType[item.value ?? 0]}`) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="state.markerType === MarkerType.CIRCLE"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="circleRadius">
                                {{ $t('common.radius') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="circleRadius">
                                <UInputNumber
                                    v-model="state.circleRadius"
                                    name="circleRadius"
                                    :min="5"
                                    :max="250"
                                    :placeholder="$t('common.radius')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-else-if="state.markerType === MarkerType.ICON"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="icon">
                                {{ $t('common.icon') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="icon">
                                <IconSelectMenu v-model="state.icon" :color="state.color" />
                            </UFormField>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="!marker ? $t('common.create') : $t('common.save')"
                    @click="formRef?.submit()"
                />

                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </USlideover>
</template>
