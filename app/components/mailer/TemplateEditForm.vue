<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import { useNotificatorStore } from '~/stores/notificator';
import type { Template } from '~~/gen/ts/resources/mailer/template';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateOrUpdateTemplateRequest } from '~~/gen/ts/services/mailer/mailer';

const props = defineProps<{
    template?: Template;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(3).max(255),
    content: z.string().min(3).max(1024),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    content: '',
});

watch(props, () => {
    state.title = props.template?.title ?? '';
    state.content = props.template?.content ?? '';
});

state.title = props.template?.title ?? '';
state.content = props.template?.content ?? '';

async function createOrUpdateTemplate(values: Schema): Promise<CreateOrUpdateTemplateRequest> {
    try {
        const call = $grpc.mailer.mailer.createOrUpdateTemplate({
            template: {
                id: props.template?.id ?? 0,
                emailId: selectedEmail.value!.id,
                title: values.title,
                content: values.content,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('refresh');
        emit('close');

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateTemplate(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <UForm
            class="mx-auto flex max-w-screen-xl flex-1 flex-col gap-y-2"
            :state="state"
            :schema="schema"
            @submit="onSubmitThrottle"
        >
            <UButtonGroup class="mb-2 flex">
                <UButton class="flex-1" type="submit" icon="i-mdi-pencil" :label="$t('common.save')" />

                <UButton icon="i-mdi-cancel" color="error" :label="$t('common.cancel')" @click="$emit('close')" />
            </UButtonGroup>

            <UFormGroup name="title" :label="$t('common.name')">
                <UInput v-model="state.title" type="text" />
            </UFormGroup>

            <UFormGroup
                class="flex flex-1 overflow-y-hidden"
                name="content"
                :ui="{ container: 'flex flex-1 flex-col mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                label="&nbsp;"
            >
                <ClientOnly>
                    <TiptapEditor
                        v-model="state.content"
                        class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                        wrapper-class="min-h-80"
                    />
                </ClientOnly>
            </UFormGroup>
        </UForm>
    </div>
</template>
