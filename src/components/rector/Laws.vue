<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { LawBook } from '~~/gen/ts/resources/laws/laws';

const { $grpc } = useNuxtApp();

const { data: lawBooks, pending, refresh, error } = useLazyAsyncData(`rector-lawbooks`, () => listLawBooks());

async function listLawBooks(): Promise<LawBook[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().listLawBooks({});
            const { response } = await call;

            return res(response.books);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="text-white px-4">
        <ul class="max-w-md space-y-1 text-gray-500 list-disc list-inside dark:text-gray-400">
            <li v-for="book in lawBooks">
                {{ book.name }} - {{ book.description }}
                <ul v-if="book.laws.length > 0" class="pl-5 mt-2 space-y-1 list-disc list-inside">
                    <li v-for="law in book.laws">{{ law.name }} - {{ law.description }}</li>
                </ul>
            </li>
        </ul>
    </div>
</template>
