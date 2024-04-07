<script lang="ts" setup>
import { digits, max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useCompletorStore } from '~/store/completor';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    entry?: ConductEntry;
    userId?: number;
}>();

const emit = defineEmits<{
    (e: 'created', entry: ConductEntry): void;
    (e: 'update', entry: ConductEntry): void;
}>();

const { isOpen } = useModal();

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const usersLoading = ref(false);

interface FormData {
    targetUser?: number;
    type: ConductType;
    message: string;
    expiresAt?: string;
}

async function conductCreateOrUpdateEntry(values: FormData, id?: string): Promise<void> {
    try {
        const expiresAt = values.expiresAt ? toTimestamp(fromString(values.expiresAt)) : undefined;

        const req = {
            entry: {
                id: id ?? '0',
                job: '',
                type: values.type,
                message: values.message,
                creatorId: 1,
                targetUserId: values.targetUser!,
                expiresAt,
            },
        };

        if (id === undefined) {
            const call = $grpc.getJobsConductClient().createConductEntry(req);
            const { response } = await call;

            emit('created', response.entry!);
        } else {
            const call = $grpc.getJobsConductClient().updateConductEntry(req);
            const { response } = await call;

            emit('update', response.entry!);
        }

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const cTypes = ref<{ status: ConductType; selected?: boolean }[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const targetUser = ref<UserShort | undefined>();
watch(targetUser, () => {
    if (targetUser.value) {
        setFieldValue('targetUser', targetUser.value.userId);
    } else {
        setFieldValue('targetUser', undefined);
    }
});

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setValues, setFieldValue, resetForm } = useForm<FormData>({
    validationSchema: {
        targetUser: { required: true },
        type: { required: true },
        message: { required: true, min: 3, max: 2000 },
        expiresAt: { required: false },
    },
    initialValues: {
        type: ConductType.NOTE,
    },
    validateOnMount: true,
});

function setFormFromProps(): void {
    resetForm();

    if (props.entry) {
        targetUser.value = props.entry.targetUser;
    }

    setValues({
        targetUser: props.entry?.targetUserId,
        type: props.entry?.type ?? ConductType.NOTE,
        message: props.entry?.message,
        expiresAt: props.entry?.expiresAt ? toDatetimeLocal(toDate(props.entry?.expiresAt)) : undefined,
    });
}

watch(props, () => setFormFromProps());

onMounted(() => setFormFromProps());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await conductCreateOrUpdateEntry(values, props.entry?.id).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{
                            entry === undefined
                                ? $t('components.jobs.conduct.CreateOrUpdateModal.create.title')
                                : $t('components.jobs.conduct.CreateOrUpdateModal.update.title')
                        }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}">
                    <div class="flex flex-1 flex-col justify-between">
                        <div class="divide-y divide-gray-200 px-2 sm:px-6">
                            <div class="mt-1">
                                <dl class="divide-neutral/10 border-neutral/10 divide-y border-b">
                                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="type" class="block text-sm font-medium leading-6">
                                                {{ $t('common.type') }}
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                v-slot="{ field }"
                                                as="div"
                                                name="type"
                                                :placeholder="$t('common.type')"
                                                :label="$t('common.type')"
                                            >
                                                <select
                                                    v-bind="field"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                >
                                                    <option
                                                        v-for="mtype in cTypes"
                                                        :key="mtype.status"
                                                        :selected="mtype.selected"
                                                        :value="mtype.status"
                                                    >
                                                        {{ $t(`enums.jobs.ConductType.${ConductType[mtype.status ?? 0]}`) }}
                                                    </option>
                                                </select>
                                            </VeeField>
                                            <VeeErrorMessage name="type" as="p" class="mt-2 text-sm text-error-400" />
                                        </dd>
                                    </div>
                                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="targetUser" class="block text-sm font-medium leading-6">
                                                {{ $t('common.target') }}
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <UInputMenu
                                                v-model="targetUser"
                                                :search="
                                                    async (query: string) => {
                                                        usersLoading = true;
                                                        const colleagues = await completorStore.listColleagues({
                                                            pagination: { offset: 0 },
                                                            searchName: query,
                                                        });
                                                        usersLoading = false;
                                                        return colleagues;
                                                    }
                                                "
                                                :search-attributes="['firstname', 'lastname']"
                                                block
                                                :placeholder="
                                                    targetUser
                                                        ? `${targetUser?.firstname} ${targetUser?.lastname} (${targetUser?.dateofbirth})`
                                                        : $t('common.target')
                                                "
                                                trailing
                                                by="userId"
                                            >
                                                <template #option="{ option: user }">
                                                    {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                                </template>
                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>
                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                                </template>
                                            </UInputMenu>
                                        </dd>
                                    </div>
                                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="message" class="block text-sm font-medium leading-6">
                                                {{ $t('common.message') }}
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                as="textarea"
                                                name="message"
                                                class="placeholder:text-accent-200 block h-36 w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :placeholder="$t('common.message')"
                                                :label="$t('common.message')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                            <VeeErrorMessage name="message" as="p" class="mt-2 text-sm text-error-400" />
                                        </dd>
                                    </div>
                                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="expiresAt" class="block text-sm font-medium leading-6">
                                                {{ $t('common.expires_at') }}?
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                type="datetime-local"
                                                name="expiresAt"
                                                :placeholder="$t('common.expires_at')"
                                                :label="$t('common.expires_at')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                            <VeeErrorMessage name="expiresAt" as="p" class="mt-2 text-sm text-error-400" />
                                        </dd>
                                    </div>
                                </dl>
                            </div>
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ entry?.id === undefined ? $t('common.create') : $t('common.update') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
