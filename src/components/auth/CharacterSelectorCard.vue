<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useClipboardStore } from 'store/clipboard';
import { parseQuery } from 'vue-router';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import { useAuthStore } from '~/store/auth';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import { User } from '~~/gen/ts/resources/users/users';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();
const route = useRoute();

const { lastCharID } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setPermissions, setJobProps } = authStore;

const props = defineProps<{
    char: User;
}>();

async function chooseCharacter(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            if (authStore.lastCharID !== props.char.userId) {
                clipboardStore.clear();
            }

            const call = $grpc.getAuthClient().chooseCharacter({
                charId: props.char.userId,
            });
            const { response } = await call;

            setAccessToken(response.token, toDate(response.expires) as null | Date);
            setActiveChar(props.char);
            setPermissions(response.permissions);
            if (response.jobProps) {
                setJobProps(response.jobProps!);
            } else {
                setJobProps(null);
            }

            const path = route.query.redirect?.toString() || '/overview';
            const url = new URL('https://example.com' + path);
            await navigateTo({
                path: url.pathname,
                query: parseQuery(url.search),
                hash: url.hash,
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div :key="char.userId" class="flex flex-col divide-y rounded-lg bg-base-800 shadow-float">
        <div class="flex flex-col flex-1 p-8">
            <div class="flex flex-row items-center gap-3 mx-auto">
                <h2 class="text-2xl font-medium text-center text-neutral">{{ char.firstname }}, {{ char.lastname }}</h2>
                <CharSexBadge :sex="char.sex!" />
                <div v-if="lastCharID === char.userId">
                    <span
                        class="inline-flex items-center rounded-full bg-success-100 px-3 py-0.5 text-sm font-medium text-success-800"
                    >
                        {{ $t('common.last_used') }}
                    </span>
                </div>
            </div>
            <dl class="flex flex-col justify-between flex-grow mt-2 text-center">
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800"
                        >{{ char.jobLabel }} ({{ $t('common.rank') }}: {{ char.jobGradeLabel }})</span
                    >
                </dd>
                <dt class="text-sm text-neutral">
                    {{ $t('common.date_of_birth') }}
                </dt>
                <dd class="text-sm text-gray-300">{{ char.dateofbirth }}</dd>
                <dt class="text-sm text-neutral">{{ $t('common.height') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.height }}cm</dd>
                <dt class="text-sm text-neutral">{{ $t('common.visum') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.visum }}</dd>
                <dt class="text-sm text-neutral">
                    {{ $t('common.playtime') }}
                </dt>
                <dd class="text-sm text-gray-300">
                    {{ fromSecondsToFormattedDuration(char.playtime!) }}
                </dd>
            </dl>
        </div>
        <div>
            <div class="flex -mt-px">
                <div class="flex flex-1 w-0">
                    <button
                        @click="chooseCharacter"
                        class="relative inline-flex items-center justify-center flex-1 w-0 py-4 text-sm font-semibold transition-colors border border-transparent rounded-b-lg gap-x-3 text-neutral bg-base-700 hover:bg-base-600"
                    >
                        {{ $t('common.choose') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
