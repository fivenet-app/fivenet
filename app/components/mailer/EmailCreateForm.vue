<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';

const props = withDefaults(
    defineProps<{
        personalEmail?: boolean;
    }>(),
    {
        personalEmail: false,
    },
);

const { activeChar } = useAuth();

const mailerStore = useMailerStore();

const firstname = slugify(activeChar.value!.firstname).replaceAll('-', '.').replaceAll('..', '.');
const lastname = slugify(activeChar.value!.lastname).replaceAll('-', '.').replaceAll('..', '.');

const schema = z.object({
    email: z.string().min(6).max(50).includes(firstname).includes(lastname),
    label: z.string().max(128),
    internal: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    email: '',
    label: '',
    internal: false,
});

onBeforeMount(() => {
    if (props.personalEmail) {
        state.email = firstname + '.' + lastname;
    }
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    const values = event.data;
    await mailerStore
        .createOrUpdateEmail({
            email: {
                id: '0',
                email: values.email,
                internal: values.internal,
                label: values.label !== '' ? values.label : undefined,
                disabled: false,
                job: !props.personalEmail ? activeChar.value!.job : undefined,
                userId: props.personalEmail ? activeChar.value!.userId : undefined,
            },
        })
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :state="state" :schema="schema" class="flex flex-col gap-y-2" @submit="onSubmitThrottle">
        <UFormGroup name="email" :label="$t('common.mail')">
            <div class="inline-flex">
                <UInput v-model="state.email" type="text" />
                <span class="text-gray-900 dark:text-white">@fivenet.app</span>
            </div>
        </UFormGroup>

        <UFormGroup v-if="!personalEmail" name="label" :label="$t('common.label')">
            <UInput v-model="state.label" type="text" />
        </UFormGroup>

        <UFormGroup>
            <UButton type="submit" block :label="$t('common.create')" />
        </UFormGroup>
    </UForm>
</template>
