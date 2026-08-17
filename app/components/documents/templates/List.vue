<script lang="ts" setup>
import CardsList from '~/components/partials/CardsList.vue';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { CardElement } from '~/utils/types';
import { resolveNeighborMovePayload } from '~/utils/reorder';
import { getDocumentsTemplatesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { TemplateShort } from '~~/gen/ts/resources/documents/templates/templates';
import type { ComponentPublicInstance } from 'vue';
import { useDraggable } from 'vue-draggable-plus';

const props = withDefaults(
    defineProps<{
        link?: boolean;
        reorderable?: boolean;
        sortMode?: boolean;
        searchTitle?: string;
    }>(),
    {
        link: false,
        reorderable: false,
        sortMode: false,
        searchTitle: undefined,
    },
);

defineEmits<{
    (e: 'selected', t: TemplateShort | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const { data: templates, status, refresh, error } = useLazyAsyncData('documents-templates', () => listTemplates());

defineExpose({
    status,
    refresh,
});

const documentsTemplatesClient = await getDocumentsTemplatesClient();
const notifications = useNotificationsStore();
const movingTemplateId = ref<number | undefined>(undefined);

async function listTemplates(): Promise<TemplateShort[]> {
    try {
        const call = documentsTemplatesClient.listTemplates({});
        const { response } = await call;

        return response.templates;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const visibleTemplates = computed<TemplateShort[]>(
    () => templates.value?.filter((v) => v.title.toLowerCase().includes(props.searchTitle?.toLowerCase() ?? '')) ?? [],
);

const items = computed<CardElement[]>(() =>
    visibleTemplates.value.map((v) => ({
        label: v.title,
        description: v.description,
        icon: v.icon ?? 'i-mdi-file-outline',
        color: v.color ?? 'primary',
        to: props.link ? `/documents/templates/${v.id}` : undefined,
        deletedAt: v.deletedAt,
    })),
);

function selected(idx: number): TemplateShort | undefined {
    return visibleTemplates.value.at(idx);
}

const sortEnabled = computed(() => props.reorderable && props.sortMode);
const sortableTemplates = computed<TemplateShort[]>({
    get: () => templates.value?.filter((template) => template.deletedAt === undefined) ?? [],
    set: (value) => {
        const deletedTemplates = templates.value?.filter((template) => template.deletedAt !== undefined) ?? [];
        templates.value = [...value, ...deletedTemplates];
    },
});

type TemplateMovePayload = {
    templateId: number;
    beforeId?: number;
    afterId?: number;
};

function resolveTemplateMovePayload(
    entries: TemplateShort[],
    oldIndex: number | undefined,
    newIndex: number | undefined,
): TemplateMovePayload | undefined {
    const payload = resolveNeighborMovePayload(entries, oldIndex, newIndex);
    if (!payload) return undefined;

    return {
        templateId: payload.id,
        beforeId: payload.beforeId,
        afterId: payload.afterId,
    };
}

async function moveTemplate(payload: TemplateMovePayload): Promise<void> {
    movingTemplateId.value = payload.templateId;

    try {
        await documentsTemplatesClient.moveTemplate(payload);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await refresh();
    } catch (e) {
        await refresh().catch(() => undefined);
        handleGRPCError(e as RpcError);
        throw e;
    } finally {
        movingTemplateId.value = undefined;
    }
}

async function moveTemplateByIndex(oldIndex: number, newIndex: number): Promise<void> {
    if (oldIndex === newIndex) return;
    if (oldIndex < 0 || oldIndex >= sortableTemplates.value.length) return;
    if (newIndex < 0 || newIndex >= sortableTemplates.value.length) return;

    const updatedTemplates = [...sortableTemplates.value];
    const moved = updatedTemplates.splice(oldIndex, 1)[0];
    if (!moved) return;
    updatedTemplates.splice(newIndex, 0, moved);
    sortableTemplates.value = updatedTemplates;

    const payload = resolveTemplateMovePayload(sortableTemplates.value, oldIndex, newIndex);
    if (!payload) {
        await refresh();
        return;
    }

    await moveTemplate(payload);
}

async function onDragEnd(event: { oldIndex?: number; newIndex?: number }): Promise<void> {
    const payload = resolveTemplateMovePayload(sortableTemplates.value, event.oldIndex, event.newIndex);
    if (!payload) return;

    await moveTemplate(payload);
}

const sortGridRef = ref<ComponentPublicInstance | HTMLElement | null>(null);
const sortListRef = computed<HTMLElement | null>(() => {
    const value = sortGridRef.value;
    if (value instanceof HTMLElement) return value;

    return value?.$el ?? null;
});
const draggable = useDraggable(sortListRef, sortableTemplates, {
    immediate: false,
    animation: 150,
    handle: '.handle',
    draggable: '> .template-sort-card',
    onEnd: onDragEnd,
});
const activeDraggableRoot = shallowRef<HTMLElement | null>(null);

watch(
    [sortEnabled, sortListRef],
    ([enabled, root]) => {
        if (!enabled || !root) {
            draggable.pause();
            activeDraggableRoot.value = null;
            return;
        }

        if (activeDraggableRoot.value !== root) {
            draggable.start(root);
            activeDraggableRoot.value = root;
        }

        if (movingTemplateId.value !== undefined) {
            draggable.pause();
            return;
        }

        draggable.resume();
    },
    { immediate: true, flush: 'post' },
);

watch(movingTemplateId, (movingId) => {
    if (!sortEnabled.value || !activeDraggableRoot.value) return;

    if (movingId !== undefined) {
        draggable.pause();
        return;
    }

    draggable.resume();
});
</script>

<template>
    <div v-if="isRequestPending(status)" class="flex justify-center">
        <UPageGrid>
            <UPageCard v-for="idx in 6" :key="idx">
                <template #title>
                    <USkeleton class="h-6 w-[275px]" />
                </template>
                <template #description>
                    <USkeleton class="h-11 w-[350px]" />
                </template>
            </UPageCard>
        </UPageGrid>
    </div>

    <DataErrorBlock
        v-else-if="error"
        :title="$t('common.unable_to_load', [$t('common.template', 2)])"
        :error="error"
        :retry="refresh"
    />
    <DataNoDataBlock v-else-if="!templates || templates.length === 0" :type="$t('common.template', 2)" />

    <DataNoDataBlock v-else-if="sortEnabled && sortableTemplates.length === 0" :type="$t('common.template', 2)" />

    <div v-else-if="sortEnabled" class="flex justify-center">
        <UPageGrid ref="sortGridRef" :class="$attrs.class">
            <UPageCard
                v-for="(template, index) in sortableTemplates"
                :key="template.id"
                class="template-sort-card"
                :title="template.title"
                :icon="template.icon?.startsWith('i-') ? template.icon : undefined"
                :ui="{ title: 'w-full flex flex-row gap-2', leading: 'gap-2' }"
            >
                <template #title>
                    <span>{{ template.title }}</span>

                    <UBadge v-if="template.deletedAt" icon="i-mdi-delete" :label="$t('common.deleted')" color="warning" />
                </template>

                <template #leading>
                    <DraggableHandle />
                    <ReorderButtons
                        :idx="index"
                        :move-up="(idx: number) => moveTemplateByIndex(idx, idx - 1)"
                        :move-down="(idx: number) => moveTemplateByIndex(idx, idx + 1)"
                        :disable-up="index === 0 || movingTemplateId !== undefined"
                        :disable-down="index >= sortableTemplates.length - 1 || movingTemplateId !== undefined"
                    />

                    <UIcon
                        v-if="template.icon"
                        class="h-10 w-10 shrink-0"
                        :class="`text-${template.color ?? 'primary'}`"
                        :name="convertComponentIconNameToDynamic(template.icon)"
                    />
                </template>

                <template #description>
                    <span class="line-clamp-2">{{ template.description }}</span>
                </template>
            </UPageCard>
        </UPageGrid>
    </div>

    <div v-else class="flex justify-center">
        <CardsList :class="$attrs.class" :items="items" @selected="$emit('selected', selected($event))" />
    </div>
</template>
