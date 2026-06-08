<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { HelpIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/markers/marker_marker';
import InputDatePicker from '../../partials/InputDatePicker.vue';
import type { Coords } from '~~/gen/ts/resources/livemap/coords.js';

const props = defineProps<{
    location?: Coords;
    marker?: MarkerMarker;
}>();

const emit = defineEmits<{
    close: [boolean];
}>();

const livemapStore = useLivemapStore();
const { location: storeLocation, showLocationMarker, markerCoordPickerActive, markersMarkers } = storeToRefs(livemapStore);
const { addOrUpdateMarkerMarker } = livemapStore;

const livemapLivemapClient = await getLivemapLivemapClient();

const markerTypes = [
    { icon: 'i-mdi-emoticon', value: MarkerType.ICON },
    { icon: 'i-mdi-dot', value: MarkerType.DOT },
    { icon: 'i-mdi-vector-circle', value: MarkerType.CIRCLE },
    { icon: 'i-mdi-vector-polyline', value: MarkerType.POLYLINE },
    { icon: 'i-mdi-vector-rectangle', value: MarkerType.RECTANGLE },
    { icon: 'i-mdi-vector-polygon', value: MarkerType.POLYGON },
];

type ShapePointsKind = 'polygon' | 'polyline';

function clonePoints(points?: Coords[]): Coords[] {
    if (!points) return [];
    return points.map((point) => ({ x: point.x, y: point.y }));
}

function cloneDataPlain<T>(value: T): T {
    return JSON.parse(JSON.stringify(toRaw(value))) as T;
}

function resolveInitialMarkerType(marker?: MarkerMarker): MarkerType {
    switch (marker?.data?.data.oneofKind) {
        case 'icon':
            return MarkerType.ICON;
        case 'circle':
            return MarkerType.CIRCLE;
        case 'polyline':
            return MarkerType.POLYLINE;
        case 'rectangle':
            return MarkerType.RECTANGLE;
        case 'polygon':
            return MarkerType.POLYGON;
    }

    if (marker?.type !== undefined && marker.type !== MarkerType.UNSPECIFIED) {
        return marker.type;
    }

    return MarkerType.ICON;
}

function getInitialCircleRadius(marker?: MarkerMarker): number {
    if (marker?.data?.data.oneofKind === 'circle' && marker.data.data.circle.radius) return marker.data.data.circle.radius;
    return 50;
}

function getInitialCircleOpacity(marker?: MarkerMarker): number {
    if (marker?.data?.data.oneofKind === 'circle' && marker.data.data.circle.opacity) return marker.data.data.circle.opacity;
    return 15;
}

function getInitialRectangleEndX(marker?: MarkerMarker, initialX = 0): number {
    if (marker?.data?.data.oneofKind === 'rectangle' && marker.data.data.rectangle.endX !== undefined) {
        return marker.data.data.rectangle.endX;
    }
    return initialX + 75;
}

function getInitialRectangleEndY(marker?: MarkerMarker, initialY = 0): number {
    if (marker?.data?.data.oneofKind === 'rectangle' && marker.data.data.rectangle.endY !== undefined) {
        return marker.data.data.rectangle.endY;
    }
    return initialY + 75;
}

function getInitialShapeOpacity(marker?: MarkerMarker): number {
    if (marker?.data?.data.oneofKind === 'rectangle') return marker.data.data.rectangle.opacity ?? 15;
    if (marker?.data?.data.oneofKind === 'polygon') return marker.data.data.polygon.opacity ?? 15;
    return 15;
}

function getInitialPolygonPoints(marker?: MarkerMarker, fallback: Coords[] = []): Coords[] {
    if (marker?.data?.data.oneofKind === 'polygon' && marker.data.data.polygon.points.length >= 2) {
        return clonePoints(marker.data.data.polygon.points);
    }
    return clonePoints(fallback);
}

function getInitialPolylinePoints(marker?: MarkerMarker, fallback: Coords[] = []): Coords[] {
    if (marker?.data?.data.oneofKind === 'polyline' && marker.data.data.polyline.points.length >= 2) {
        return clonePoints(marker.data.data.polyline.points);
    }
    return clonePoints(fallback);
}

function getInitialIcon(marker?: MarkerMarker): string {
    if (marker?.data?.data.oneofKind === 'icon' && marker.data.data.icon.icon) return marker.data.data.icon.icon;
    return HelpIcon.name ?? 'i-mdi-help';
}

const defaultExpiresAt = ref<Date>(new Date());
defaultExpiresAt.value.setTime(defaultExpiresAt.value.getTime() + 1 * 60 * 60 * 1000);

const initialX = props.marker?.x ?? props.location?.x ?? storeLocation.value?.x ?? 0;
const initialY = props.marker?.y ?? props.location?.y ?? storeLocation.value?.y ?? 0;

const defaultPolygonPoints: Coords[] = [
    { x: initialX + 75, y: initialY },
    { x: initialX + 35, y: initialY + 75 },
];

const defaultPolylinePoints: Coords[] = [
    { x: initialX + 75, y: initialY },
    { x: initialX + 35, y: initialY + 75 },
    { x: initialX + 70, y: initialY + 150 },
];

const schema = z.object({
    name: z.coerce.string().min(1).max(255),
    description: z.union([z.string().min(3).max(1024), z.string().length(0).optional()]),
    expiresAt: z.date().optional(),
    color: z.coerce.string().length(7),
    x: z.coerce.number(),
    y: z.coerce.number(),
    markerType: z.enum(MarkerType),
    circleRadius: z.coerce.number().gte(5).lte(250),
    circleOpacity: z.coerce.number().gte(1).lte(75).optional(),
    rectangleEndX: z.coerce.number(),
    rectangleEndY: z.coerce.number(),
    shapeOpacity: z.coerce.number().gte(1).lte(75).optional(),
    polygonPoints: z
        .array(z.object({ x: z.coerce.number(), y: z.coerce.number() }))
        .min(2)
        .max(18),
    polylinePoints: z
        .array(z.object({ x: z.coerce.number(), y: z.coerce.number() }))
        .min(2)
        .max(18),
    icon: z.string().max(128).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: props.marker?.name ?? '',
    description: props.marker?.description,
    expiresAt: props.marker?.expiresAt ? toDate(props.marker?.expiresAt) : defaultExpiresAt.value,
    color: props.marker?.color ?? '#ee4b2b',
    x: initialX,
    y: initialY,
    markerType: resolveInitialMarkerType(props.marker),
    circleRadius: getInitialCircleRadius(props.marker),
    circleOpacity: getInitialCircleOpacity(props.marker),
    rectangleEndX: getInitialRectangleEndX(props.marker, initialX),
    rectangleEndY: getInitialRectangleEndY(props.marker, initialY),
    shapeOpacity: getInitialShapeOpacity(props.marker),
    polygonPoints: getInitialPolygonPoints(props.marker, defaultPolygonPoints),
    polylinePoints: getInitialPolylinePoints(props.marker, defaultPolylinePoints),
    icon: getInitialIcon(props.marker),
});

const isPointShape = computed<boolean>(
    () => state.markerType === MarkerType.POLYGON || state.markerType === MarkerType.POLYLINE,
);
const currentShapePoints = computed<Coords[]>(() => getShapePoints());
const currentShapeFieldName = computed<string>(() =>
    state.markerType === MarkerType.POLYLINE ? 'polylinePoints' : 'polygonPoints',
);

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            name: value.name,
            description: value.description ?? '',
            expiresAt: value.expiresAt ? value.expiresAt.toISOString() : '',
            color: value.color,
            x: value.x,
            y: value.y,
            markerType: value.markerType,
            circleRadius: value.circleRadius,
            circleOpacity: value.circleOpacity ?? null,
            rectangleEndX: value.rectangleEndX,
            rectangleEndY: value.rectangleEndY,
            shapeOpacity: value.shapeOpacity ?? null,
            polygonPoints: value.polygonPoints.map((point) => ({ x: point.x, y: point.y })),
            polylinePoints: value.polylinePoints.map((point) => ({ x: point.x, y: point.y })),
            icon: value.icon ?? '',
        }),
});

function getShapePoints(): Coords[] {
    return state.markerType === MarkerType.POLYLINE ? state.polylinePoints : state.polygonPoints;
}

const isPickingCoordinates = computed<boolean>(() => markerCoordPickerActive.value === true);
const rectangleSelectionAwaitingEnd = ref<boolean>(false);
const shapeSelectionIndex = ref<number | undefined>(undefined);
const shapeVertexEditIndex = ref<number | undefined>(undefined);

const saved = ref<boolean>(false);
const previewTarget = computed<MarkerMarker | undefined>(() =>
    props.marker?.id ? markersMarkers.value.get(props.marker.id) : undefined,
);
const originalPreviewState =
    previewTarget.value === undefined
        ? undefined
        : {
              x: previewTarget.value.x,
              y: previewTarget.value.y,
              type: previewTarget.value.type,
              data: previewTarget.value.data ? cloneDataPlain(previewTarget.value.data) : undefined,
          };

const coordinatePickerHint = computed<string | undefined>(() => {
    if (!isPickingCoordinates.value) return undefined;
    if (state.markerType === MarkerType.RECTANGLE) {
        return rectangleSelectionAwaitingEnd.value ? '2/2' : '1/2';
    }

    if (isPointShape.value) {
        if (shapeVertexEditIndex.value !== undefined) {
            return `V${shapeVertexEditIndex.value + 1}`;
        }
        const index = shapeSelectionIndex.value ?? 0;
        const total = getShapePoints().length + 1;
        return `${Math.min(index + 1, total)}/${total}`;
    }

    return undefined;
});

function resetCoordinatePickingProgress(): void {
    rectangleSelectionAwaitingEnd.value = false;
    shapeSelectionIndex.value = undefined;
    shapeVertexEditIndex.value = undefined;
}

function stopCoordinatePicking(): void {
    markerCoordPickerActive.value = false;
    resetCoordinatePickingProgress();
    showLocationMarker.value = false;
}

watch([() => state.x, () => state.y], ([x, y]) => {
    if (!previewTarget.value) return;
    previewTarget.value.x = x;
    previewTarget.value.y = y;
});

watch([() => state.circleRadius, () => state.circleOpacity], ([radius, opacity]) => {
    if (!previewTarget.value || state.markerType !== MarkerType.CIRCLE || previewTarget.value.data?.data.oneofKind !== 'circle')
        return;
    previewTarget.value.data.data.circle.radius = radius;
    previewTarget.value.data.data.circle.opacity = opacity;
});

watch([() => state.rectangleEndX, () => state.rectangleEndY, () => state.shapeOpacity], ([endX, endY, opacity]) => {
    if (
        !previewTarget.value ||
        state.markerType !== MarkerType.RECTANGLE ||
        previewTarget.value.data?.data.oneofKind !== 'rectangle'
    )
        return;
    previewTarget.value.data.data.rectangle.endX = endX;
    previewTarget.value.data.data.rectangle.endY = endY;
    previewTarget.value.data.data.rectangle.opacity = opacity;
});

function syncPointShapePreview(kind: ShapePointsKind, points: Coords[]): void {
    const data = previewTarget.value?.data?.data;
    if (!data) return;

    if (kind === 'polygon') {
        if (data.oneofKind !== 'polygon') return;
        data.polygon.points = clonePoints(toRaw(points));
    } else {
        if (data.oneofKind !== 'polyline') return;
        data.polyline.points = clonePoints(toRaw(points));
    }
}

watch(
    () => state.polygonPoints,
    (points) => syncPointShapePreview('polygon', points),
    { deep: true },
);

watch(
    () => state.polylinePoints,
    (points) => syncPointShapePreview('polyline', points),
    { deep: true },
);

watch(
    () => state.shapeOpacity,
    (opacity) => {
        if (!previewTarget.value || previewTarget.value.data?.data.oneofKind === 'polyline') return;
        if (previewTarget.value.data?.data.oneofKind === 'rectangle') {
            previewTarget.value.data.data.rectangle.opacity = opacity;
        } else if (previewTarget.value.data?.data.oneofKind === 'polygon') {
            previewTarget.value.data.data.polygon.opacity = opacity;
        }
    },
);

watch(
    () => storeLocation.value,
    (val) => {
        if (!isPickingCoordinates.value || !val) return;

        if (state.markerType === MarkerType.RECTANGLE) {
            if (!rectangleSelectionAwaitingEnd.value) {
                state.x = val.x;
                state.y = val.y;
                rectangleSelectionAwaitingEnd.value = true;
                showLocationMarker.value = true;
                return;
            }

            state.rectangleEndX = val.x;
            state.rectangleEndY = val.y;
            stopCoordinatePicking();
            return;
        }

        if (isPointShape.value) {
            const points = getShapePoints();

            if (shapeVertexEditIndex.value !== undefined) {
                if (shapeVertexEditIndex.value === 0) {
                    state.x = val.x;
                    state.y = val.y;
                } else if (shapeVertexEditIndex.value - 1 < points.length) {
                    points[shapeVertexEditIndex.value - 1] = {
                        x: val.x,
                        y: val.y,
                    };
                }

                stopCoordinatePicking();
                return;
            }

            const currentIndex = shapeSelectionIndex.value ?? 0;
            const totalPoints = points.length + 1;
            if (currentIndex === 0) {
                state.x = val.x;
                state.y = val.y;
            } else if (currentIndex - 1 < points.length) {
                points[currentIndex - 1] = {
                    x: val.x,
                    y: val.y,
                };
            }

            const nextIndex = currentIndex + 1;
            if (nextIndex >= totalPoints) {
                stopCoordinatePicking();
            } else {
                shapeSelectionIndex.value = nextIndex;
            }
            return;
        }

        state.x = val.x;
        state.y = val.y;
    },
);

watch(
    () => state.markerType,
    () => {
        resetCoordinatePickingProgress();
    },
);

function toggleCoordinatePicker(): void {
    if (markerCoordPickerActive.value) {
        stopCoordinatePicking();
        return;
    }

    markerCoordPickerActive.value = true;
    resetCoordinatePickingProgress();
    if (state.markerType === MarkerType.RECTANGLE) {
        showLocationMarker.value = false;
    } else if (isPointShape.value) {
        shapeSelectionIndex.value = 0;
        showLocationMarker.value = false;
    } else {
        storeLocation.value = { x: state.x, y: state.y };
        showLocationMarker.value = true;
    }
}

function beginShapeVertexEdit(index: number): void {
    if (index < 0 || index > getShapePoints().length) return;
    resetCoordinatePickingProgress();
    shapeVertexEditIndex.value = index;
    markerCoordPickerActive.value = true;
    showLocationMarker.value = false;
}

function addShapePoint(): void {
    const points = getShapePoints();
    if (points.length >= 18) return;

    const source = points[points.length - 1] ?? { x: state.x, y: state.y };
    points.push({
        x: source.x + 25,
        y: source.y + 25,
    });
}

function removeShapePoint(index: number): void {
    const points = getShapePoints();
    if (points.length <= 2) return;
    shapeVertexEditIndex.value = undefined;
    points.splice(index, 1);
}

function createMarkerData(values: Schema): MarkerMarker['data'] | undefined {
    switch (values.markerType) {
        case MarkerType.CIRCLE:
            return {
                data: {
                    oneofKind: 'circle',
                    circle: {
                        radius: values.circleRadius,
                        opacity: values.circleOpacity ?? 3,
                    },
                },
            };
        case MarkerType.ICON:
            return {
                data: {
                    oneofKind: 'icon',
                    icon: {
                        icon: values.icon ?? 'i-mdi-help',
                    },
                },
            };
        case MarkerType.RECTANGLE:
            return {
                data: {
                    oneofKind: 'rectangle',
                    rectangle: {
                        endX: values.rectangleEndX,
                        endY: values.rectangleEndY,
                        opacity: values.shapeOpacity ?? 15,
                    },
                },
            };
        case MarkerType.POLYGON:
            return {
                data: {
                    oneofKind: 'polygon',
                    polygon: {
                        points: clonePoints(values.polygonPoints),
                        opacity: values.shapeOpacity ?? 15,
                    },
                },
            };
        case MarkerType.POLYLINE:
            return {
                data: {
                    oneofKind: 'polyline',
                    polyline: {
                        points: clonePoints(values.polylinePoints),
                    },
                },
            };
        default:
            return undefined;
    }
}

async function createOrUpdateMarker(values: Schema): Promise<void> {
    const expiresAt = values.expiresAt ? toTimestamp(values.expiresAt) : undefined;

    try {
        stopCoordinatePicking();
        const marker: MarkerMarker = {
            id: props.marker?.id ?? 0,
            job: '',
            jobLabel: '',
            name: values.name,
            description: values.description,
            x: values.x,
            y: values.y,
            color: values.color,
            expiresAt: expiresAt,
            type: values.markerType,
            data: createMarkerData(values),
        };

        const call = livemapLivemapClient.createOrUpdateMarker({
            marker,
        });
        const { response } = await call;

        if (response.marker !== undefined) {
            addOrUpdateMarkerMarker(response.marker);
        }

        saved.value = true;
        emit('close', true);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateMarker(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeSlideover(): Promise<void> {
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}

onBeforeUnmount(() => {
    stopCoordinatePicking();
    if (!saved.value && previewTarget.value && originalPreviewState) {
        previewTarget.value.x = originalPreviewState.x;
        previewTarget.value.y = originalPreviewState.y;
        previewTarget.value.type = originalPreviewState.type;
        previewTarget.value.data = originalPreviewState.data;
    }
});
</script>

<template>
    <USlideover
        :title="!marker ? $t('components.livemap.create_marker.title') : $t('components.livemap.update_marker.title')"
        :overlay="false"
        :modal="false"
        :close="{ onClick: closeSlideover }"
        :dismissible="!isPickingCoordinates && !hasUnsavedChanges"
        :ui="{ content: isPickingCoordinates ? 'max-w-sm' : 'max-w-xl' }"
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="name">
                                {{ $t('common.name') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="name" required>
                                <UInput
                                    v-model="state.name"
                                    class="w-full"
                                    type="text"
                                    name="name"
                                    :placeholder="$t('common.name')"
                                />
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
                                    class="w-full"
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
                                <InputDatePicker v-model="state.expiresAt" clearable time />
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
                                <ColorPicker v-model="state.color" class="w-full" />
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
                                        class="w-full"
                                        name="markerType"
                                        :items="markerTypes"
                                        value-key="value"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #default>
                                            {{ $t(`enums.livemap.MarkerType.${MarkerType[state.markerType ?? 0]}`) }}
                                        </template>

                                        <template #item-label="{ item }">
                                            {{ $t(`enums.livemap.MarkerType.${MarkerType[item.value ?? 0]}`) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="x">
                                {{ $t('common.longitude') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="x">
                                <UInputNumber
                                    v-model="state.x"
                                    class="w-full"
                                    name="x"
                                    :step="0.00001"
                                    :placeholder="$t('common.longitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="y">
                                {{ $t('common.latitude') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="y">
                                <UInputNumber
                                    v-model="state.y"
                                    class="w-full"
                                    name="y"
                                    :step="0.00001"
                                    :placeholder="$t('common.latitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium" />
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UButton
                                :variant="isPickingCoordinates ? 'solid' : 'soft'"
                                :color="isPickingCoordinates ? 'warning' : 'neutral'"
                                icon="i-mdi-crosshairs-gps"
                                :label="isPickingCoordinates ? $t('common.cancel') : $t('common.select')"
                                type="button"
                                @click="toggleCoordinatePicker"
                            />
                            <p v-if="coordinatePickerHint" class="mt-2 text-xs text-dimmed">{{ coordinatePickerHint }}</p>
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
                                    class="w-full"
                                    name="circleRadius"
                                    :min="5"
                                    :max="250"
                                    :placeholder="$t('common.radius')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="state.markerType === MarkerType.RECTANGLE"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="rectangleEndX">
                                {{ $t('common.longitude') }} 2
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="rectangleEndX">
                                <UInputNumber
                                    v-model="state.rectangleEndX"
                                    class="w-full"
                                    name="rectangleEndX"
                                    :step="0.00001"
                                    :placeholder="$t('common.longitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="state.markerType === MarkerType.RECTANGLE"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="rectangleEndY">
                                {{ $t('common.latitude') }} 2
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="rectangleEndY">
                                <UInputNumber
                                    v-model="state.rectangleEndY"
                                    class="w-full"
                                    name="rectangleEndY"
                                    :step="0.00001"
                                    :placeholder="$t('common.latitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="state.markerType === MarkerType.RECTANGLE || state.markerType === MarkerType.POLYGON"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="shapeOpacity">
                                {{ $t('common.opacity') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="shapeOpacity">
                                <UInputNumber
                                    v-model="state.shapeOpacity"
                                    class="w-full"
                                    name="shapeOpacity"
                                    :min="1"
                                    :max="75"
                                    :placeholder="$t('common.opacity')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div v-if="isPointShape" class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium">
                                {{ $t('common.points', 2) }}
                            </label>
                        </dt>

                        <dd class="mt-1 space-y-2 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <div class="rounded-md border border-default p-2">
                                <div class="mb-1 flex items-center justify-between text-xs text-dimmed">
                                    <span>{{ $t('common.point', 1) }} 1</span>
                                    <UButton
                                        color="neutral"
                                        variant="soft"
                                        icon="i-mdi-crosshairs-gps"
                                        size="xs"
                                        :label="$t('common.select')"
                                        type="button"
                                        @click="beginShapeVertexEdit(0)"
                                    />
                                </div>
                                <div class="grid grid-cols-2 gap-2">
                                    <UFormField :name="`${currentShapeFieldName}.0.x`">
                                        <UInputNumber
                                            v-model="state.x"
                                            class="w-full"
                                            :name="`${currentShapeFieldName}.0.x`"
                                            :step="0.00001"
                                            :placeholder="$t('common.longitude')"
                                        />
                                    </UFormField>

                                    <UFormField :name="`${currentShapeFieldName}.0.y`">
                                        <UInputNumber
                                            v-model="state.y"
                                            class="w-full"
                                            :name="`${currentShapeFieldName}.0.y`"
                                            :step="0.00001"
                                            :placeholder="$t('common.latitude')"
                                        />
                                    </UFormField>
                                </div>
                            </div>

                            <div
                                v-for="(point, idx) in currentShapePoints"
                                :key="`${currentShapeFieldName}-${idx}`"
                                class="rounded-md border border-default p-2"
                            >
                                <div class="mb-1 flex items-center justify-between text-xs text-dimmed">
                                    <span>{{ $t('common.point', 1) }} {{ idx + 2 }}</span>
                                    <div class="flex items-center gap-1">
                                        <UButton
                                            color="neutral"
                                            variant="soft"
                                            icon="i-mdi-crosshairs-gps"
                                            size="xs"
                                            :label="$t('common.select')"
                                            type="button"
                                            @click="beginShapeVertexEdit(idx + 1)"
                                        />
                                        <UButton
                                            color="error"
                                            variant="link"
                                            icon="i-mdi-delete"
                                            size="xs"
                                            type="button"
                                            :disabled="currentShapePoints.length <= 2"
                                            @click="removeShapePoint(idx)"
                                        />
                                    </div>
                                </div>

                                <div class="grid grid-cols-2 gap-2">
                                    <UFormField :name="`${currentShapeFieldName}.${idx + 1}.x`">
                                        <UInputNumber
                                            v-model="point.x"
                                            class="w-full"
                                            :name="`${currentShapeFieldName}.${idx + 1}.x`"
                                            :step="0.00001"
                                            :placeholder="$t('common.longitude')"
                                        />
                                    </UFormField>

                                    <UFormField :name="`${currentShapeFieldName}.${idx + 1}.y`">
                                        <UInputNumber
                                            v-model="point.y"
                                            class="w-full"
                                            :name="`${currentShapeFieldName}.${idx + 1}.y`"
                                            :step="0.00001"
                                            :placeholder="$t('common.latitude')"
                                        />
                                    </UFormField>
                                </div>
                            </div>

                            <UButton
                                :label="$t('common.add')"
                                icon="i-mdi-plus"
                                variant="soft"
                                type="button"
                                :disabled="currentShapePoints.length >= 18"
                                @click="addShapePoint"
                            />
                        </dd>
                    </div>

                    <div v-if="state.markerType === MarkerType.ICON" class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="icon">
                                {{ $t('common.icon') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="icon">
                                <IconSelectMenu v-model="state.icon" class="w-full" :hex-color="state.color" />
                            </UFormField>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="!marker ? $t('common.create') : $t('common.save')"
                    @click="formRef?.submit()"
                />

                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="closeSlideover" />
            </UFieldGroup>
        </template>
    </USlideover>
</template>
