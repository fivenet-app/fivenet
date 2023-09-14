<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { max, min, required } from '@vee-validate/rules';
import { useConfirmDialog } from '@vueuse/core';
import { CancelIcon, ContentSaveIcon, PencilIcon, TrashCanIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { LawBook } from '~~/gen/ts/resources/laws/laws';
import LawEntry from './LawEntry.vue';

const props = defineProps<{
    book: LawBook;
    startInEdit?: boolean;
}>();

const emit = defineEmits<{
    (e: 'deleted', id: bigint): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteLawBook(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        if (id < 0) {
            emit('deleted', id);
            return;
        }

        try {
            const call = $grpc.getRectorClient().deleteLawBook({
                id: id,
            });
            await call;

            emit('deleted', id);

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function saveLawBook(id: bigint, values: FormData): Promise<LawBook> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().createOrUpdateLawBook({
                id: BigInt(id < 0 ? 0 : id),
                name: values.name,
                description: values.description,
                laws: [],
            });
            const { response } = await call;

            props.book.id = response.id;
            props.book.createdAt = response.createdAt;
            props.book.updatedAt = response.updatedAt;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

interface FormData {
    id: bigint;
    name: string;
    description?: string;
}

const { handleSubmit, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});
setValues({
    name: props.book.name,
    description: props.book.description,
});

const onSubmit = handleSubmit(async (values): Promise<LawBook> => await saveLawBook(props.book.id, values));

function deletedLaw(id: bigint): void {
    const idx = props.book.laws.findIndex((b) => b.id === id);
    if (idx > -1) props.book.laws.splice(idx, 1);
}

const lastNewId = ref<bigint>(BigInt(-1));

function addLaw(): void {
    props.book.laws?.unshift({
        lawbookId: props.book.id,
        id: lastNewId.value,
        name: '',
        fine: BigInt(0),
        detentionTime: BigInt(0),
        stvoPoints: BigInt(0),
    });
    lastNewId.value--;
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteLawBook(id));

const editing = ref(props.startInEdit);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(book.id)" />

    <div class="my-2">
        <div v-if="!editing" class="flex text-white items-center gap-x-2">
            <button type="button" @click="editing = true" :title="$t('common.edit')">
                <PencilIcon class="w-6 h-6" />
            </button>
            <button type="button" @click="reveal()" :title="$t('common.delete')">
                <TrashCanIcon class="w-6 h-6" />
            </button>
            <h2 class="text-xl">{{ book.name }}</h2>
            <p v-if="book.description" class="pl-2">- {{ $t('common.description') }}: {{ book.description }}</p>
            <div class="pl-2">
                <div class="sm:flex-auto w-full">
                    <button
                        type="button"
                        @click="addLaw"
                        class="px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    >
                        {{ $t('pages.rector.laws.add_new_law') }}
                    </button>
                </div>
            </div>
        </div>
        <form v-else @submit="onSubmit" class="w-full flex flex-row gap-x-4 text-white items-start">
            <button type="submit" :title="$t('common.save')">
                <ContentSaveIcon class="w-6 h-6" />
            </button>
            <button
                type="button"
                @click="
                    editing = false;
                    book.id < BigInt(0) && $emit('deleted', book.id);
                "
                :title="$t('common.cancel')"
            >
                <CancelIcon class="w-6 h-6" />
            </button>

            <div class="flex-initial">
                <label for="name">
                    {{ $t('common.law_book') }}
                </label>
                <VeeField
                    name="name"
                    type="text"
                    :placeholder="$t('common.law_book')"
                    :label="$t('common.law_book')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                />
                <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <div class="flex-auto">
                <label for="description">
                    {{ $t('common.description') }}
                </label>
                <VeeField
                    name="description"
                    type="text"
                    :placeholder="$t('common.description')"
                    :label="$t('common.description')"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                />
                <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
            </div>
        </form>
        <table class="min-w-full divide-y divide-base-600">
            <thead>
                <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                        {{ $t('common.action', 2) }}
                    </th>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                        {{ $t('common.crime') }}
                    </th>
                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.fine') }}
                    </th>
                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.detention_time') }}
                    </th>
                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.traffic_infraction_points', 2) }}
                    </th>
                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                        {{ $t('common.description') }}
                    </th>
                </tr>
            </thead>
            <tbody class="divide-y divide-base-800">
                <LawEntry
                    v-for="law in book.laws"
                    :key="law.id.toString()"
                    :law="law"
                    :start-in-edit="law.id < BigInt(0)"
                    @deleted="deletedLaw($event)"
                />
            </tbody>
        </table>
    </div>
</template>
