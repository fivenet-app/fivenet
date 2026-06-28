<script lang="ts" setup>
import type { EmptyProps } from '@nuxt/ui';

type DataNoDataBlockProps = {
    title?: EmptyProps['title'];
    message?: EmptyProps['description'];
    icon?: EmptyProps['icon'];
    type?: string;
    actions?: NonNullable<EmptyProps['actions']>;
    focus?: () => void | Promise<void>;
    retry?: () => Promise<unknown>;
    padded?: boolean;
};

const props = withDefaults(defineProps<DataNoDataBlockProps>(), {
    title: undefined,
    message: undefined,
    icon: 'i-mdi-magnify',
    type: undefined,
    actions: () => [],
    focus: undefined,
    retry: undefined,
    padded: true,
});

const { t } = useI18n();

const actions = computed<NonNullable<EmptyProps['actions']>>(() =>
    props.actions.length > 0
        ? props.actions
        : [
              props.focus
                  ? {
                        label: t('common.search'),
                        icon: 'i-mdi-search',
                        onClick: () => {
                            props.focus!();
                        },
                    }
                  : undefined,
              props.retry
                  ? {
                        label: t('common.refresh'),
                        icon: 'i-mdi-refresh',
                        onClick: () => {
                            props.retry!();
                        },
                    }
                  : undefined,
          ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const message = computed(() => props.message ?? $t('common.not_found', [props.type ?? $t('common.data')]));

const description = computed(() => (props.title ? message.value : undefined));
const title = computed(() => (props.title ? props.title : message.value));

async function click() {
    if (props.retry) {
        props.retry();
    } else if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <div :class="padded ? 'm-2' : ''">
        <UEmpty
            :icon="icon"
            variant="outline"
            :title="title"
            :description="description"
            :actions="actions"
            v-bind="$attrs"
            @click="click()"
        >
            <template v-if="$slots.title" #title>
                <slot name="title" />
            </template>

            <template v-if="$slots.description" #description>
                <slot name="description" />
            </template>
        </UEmpty>
    </div>
</template>
