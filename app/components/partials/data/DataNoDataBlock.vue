<script lang="ts" setup>
import type { AlertProps } from '@nuxt/ui';

type DataNoDataBlockProps = {
    title?: AlertProps['title'];
    message?: AlertProps['description'];
    icon?: AlertProps['icon'];
    type?: string;
    actions?: NonNullable<AlertProps['actions']>;
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

const actions = computed<NonNullable<AlertProps['actions']>>(() =>
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
        <UAlert
            :icon="icon"
            variant="outline"
            :title="title"
            :description="message ?? $t('common.not_found', [type ?? $t('common.data')])"
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
        </UAlert>
    </div>
</template>
