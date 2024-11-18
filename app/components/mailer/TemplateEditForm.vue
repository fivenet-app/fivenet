<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import type { Template } from '~~/gen/ts/resources/mailer/template';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateOrUpdateTemplateRequest } from '~~/gen/ts/services/mailer/mailer';
import DocEditor from '../partials/DocEditor.vue';

const props = defineProps<{
    template?: Template;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

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
        const call = getGRPCMailerClient().createOrUpdateTemplate({
            template: {
                id: props.template?.id ?? '0',
                emailId: selectedEmail.value!.id,
                title: values.title,
                content: values.content,
            },
        });
        const { response } = await call;

        emits('refresh');

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emits('close');

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
        <UForm :state="state" :schema="schema" class="mx-auto flex max-w-screen-xl flex-col gap-y-2" @submit="onSubmitThrottle">
            <UButtonGroup class="mb-2 flex">
                <UButton type="submit" class="flex-1" icon="i-mdi-pencil" :label="$t('common.save')" />

                <UButton icon="i-mdi-cancel" color="red" :label="$t('common.cancel')" @click="$emit('close')" />
            </UButtonGroup>

            <UFormGroup name="title" :label="$t('common.name')">
                <UInput v-model="state.title" type="text" />
            </UFormGroup>

            <UFormGroup name="content">
                <ClientOnly>
                    <DocEditor v-model="state.content" />
                </ClientOnly>
            </UFormGroup>
        </UForm>
    </div>
</template>
