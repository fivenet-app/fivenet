<script lang="ts" setup>
defineEmits<{
    (e: 'clicked'): void;
}>();

const props = withDefaults(
    defineProps<{
        title: string;
        message?: string | undefined;
        icon?: string;
        type?: 'error' | 'info' | 'success';
        callbackMessage?: string;
    }>(),
    {
        message: undefined,
        icon: 'i-mdi-close-circle',
        type: 'error',
        callbackMessage: undefined,
    },
);

const bgClass = ref('');
const textClass = ref('');
const messageClass = ref('');
const iconClass = ref('');
const buttonClass = ref('');

switch (props.type) {
    case 'error':
        bgClass.value = 'bg-error-100';
        textClass.value = 'text-error-800';
        messageClass.value = 'text-error-700';
        iconClass.value = 'text-error-400';
        buttonClass.value = 'bg-error-200 hover:bg-error-400 focus:ring-error-600 focus-visible:outline-error-500';
        break;

    case 'success':
        bgClass.value = 'bg-success-100';
        textClass.value = 'text-success-800';
        messageClass.value = 'text-success-700';
        iconClass.value = 'text-success-400';
        buttonClass.value = 'bg-success-200 hover:bg-success-400 focus:ring-success-600 focus-visible:outline-success-500';
        break;

    case 'info':
    default:
        bgClass.value = 'bg-info-100';
        textClass.value = 'text-info-800';
        messageClass.value = 'text-info-700';
        iconClass.value = 'text-info-400';
        buttonClass.value = 'bg-info-200 hover:bg-info-400 focus:ring-info-600 focus-visible:outline-info-500';
        break;
}
</script>

<template>
    <div class="mt-6 rounded-md p-4" :class="bgClass">
        <div class="flex">
            <div class="shrink-0">
                <component :is="icon" class="size-5" :class="iconClass" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium" :class="textClass">
                    {{ title }}
                </h3>
                <div v-if="message" class="mt-2 text-sm" :class="messageClass">
                    <p>{{ message }}</p>
                </div>
                <div v-if="callbackMessage" class="mt-4">
                    <div class="-mx-2 -my-1.5 flex">
                        <UButton
                            class="flex justify-center rounded-md px-2 py-1.5 text-sm font-medium focus:ring-2 focus:ring-offset-2"
                            :class="[textClass, buttonClass]"
                            @click="$emit('clicked')"
                        >
                            {{ callbackMessage ?? $t('common.retry') }}
                            <UIcon name="i-mdi-gesture-tap" class="ml-2 size-5" />
                        </UButton>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
