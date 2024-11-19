<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import AccessManager from '../partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '../partials/access/helpers';
import { defaultEmailDomain } from './helpers';

const props = withDefaults(
    defineProps<{
        modelValue?: Email;
        personalEmail?: boolean;
        disabled?: boolean;
        hideLabel?: boolean;
    }>(),
    {
        modelValue: undefined,
        personalEmail: false,
        disabled: false,
        hideLabel: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', email: Email | undefined): void;
    (e: 'refresh'): void;
}>();

const notifications = useNotificatorStore();

const email = computed({
    get() {
        return (
            props.modelValue ?? {
                id: '0',
                email: props.personalEmail ? firstname + '.' + lastname + defaultEmailDomain : '',
                internal: false,
                disabled: false,
                job: !props.personalEmail ? activeChar.value!.job : undefined,
                userId: props.personalEmail ? activeChar.value!.userId : undefined,
                access: {
                    jobs: [],
                    qualifications: [],
                    users: [],
                },
            }
        );
    },
    set(value: Email | undefined) {
        emit('update:modelValue', value);
    },
});

const { activeChar } = useAuth();

const mailerStore = useMailerStore();

const firstname = slugify(activeChar.value!.firstname).replaceAll('-', '.').replaceAll('..', '.');
const lastname = slugify(activeChar.value!.lastname).replaceAll('-', '.').replaceAll('..', '.');

const schema = z.object({
    email: props.personalEmail
        ? z.string().min(6).max(50).includes(firstname).includes(lastname).endsWith(defaultEmailDomain)
        : z.string().min(6).max(50),
    label: z.string().max(128).optional(),
    internal: z.boolean(),
});

type Schema = z.output<typeof schema>;

async function createOrUpdateEmail(): Promise<undefined> {
    const response = await mailerStore.createOrUpdateEmail({
        email: {
            id: email.value.id,
            email: email.value.email,
            internal: email.value.internal,
            label: email.value.label !== '' ? email.value.label : undefined,
            disabled: email.value.disabled,
            job: email.value.job,
            userId: email.value.userId,
            access: {
                jobs: [],
                qualifications: [],
                users: [],
            },
        },
    });

    if (response.email) {
        email.value.email = response.email.email;
        email.value.label = response.email.label;
    }

    notifications.add({
        title: { key: 'notifications.action_successfull.title', parameters: {} },
        description: { key: 'notifications.action_successfull.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    emit('refresh');
}

watch(
    () => props.modelValue,
    () => {
        if (email.value.access === undefined) {
            email.value.access = {
                jobs: [],
                qualifications: [],
                users: [],
            };
        }
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (_: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await createOrUpdateEmail().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :state="email" :schema="schema" class="flex flex-col gap-y-2" @submit="onSubmitThrottle">
        <UFormGroup name="email" :label="$t('common.mail')">
            <UInput v-model="email.email" type="text" :placeholder="$t('common.mail')" :disabled="disabled" />
        </UFormGroup>

        <UFormGroup v-if="!hideLabel" name="label" :label="$t('common.label')">
            <UInput v-model="email.label" type="text" />
        </UFormGroup>

        <UFormGroup v-if="!personalEmail && email.access" name="access" :label="$t('common.access')">
            <AccessManager
                v-model:jobs="email.access!.jobs"
                v-model:users="email.access!.users"
                v-model:qualifications="email.access!.qualifications"
                :target-id="email.id ?? '0'"
                :access-types="[
                    { type: 'user', name: $t('common.citizen', 2) },
                    { type: 'job', name: $t('common.job', 2) },
                    { type: 'qualification', name: $t('common.qualification', 2) },
                ]"
                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.mailer.AccessLevel')"
                :disabled="disabled"
            />
        </UFormGroup>

        <UFormGroup v-if="!disabled">
            <UButton type="submit" block :label="email.id === '0' ? $t('common.create') : $t('common.update')" />
        </UFormGroup>
    </UForm>
</template>
