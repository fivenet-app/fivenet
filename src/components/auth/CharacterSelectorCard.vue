<script lang="ts" setup>
import { computed } from 'vue';
import { useAuthStore } from '~/store/auth';
import { ChooseCharacterRequest } from '@fivenet/gen/services/auth/auth_pb';
import { User } from '@fivenet/gen/resources/users/users_pb';
import CharSexBadge from '~/components/citizens/CharSexBadge.vue';
import { fromSecondsToFormattedDuration } from '~/utils/time';
import { RpcError } from 'grpc-web';
import { parseQuery } from 'vue-router';

const { $grpc } = useNuxtApp();
const store = useAuthStore();
const route = useRoute();

const lastCharID = computed(() => store.$state.lastCharID);

const props = defineProps({
    char: {
        required: true,
        type: User,
    },
});

async function chooseCharacter(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new ChooseCharacterRequest();
        req.setCharId(props.char.getUserId());

        try {
            const resp = await $grpc.getAuthClient()
                .chooseCharacter(req, null);

            store.updateAccessToken(resp.getToken());
            store.updateActiveChar(props.char);
            store.updatePermissions(resp.getPermissionsList());

            const path = route.query.redirect?.toString() || "/overview";
            const url = new URL("https://example.com" + path);
            await navigateTo({ path: url.pathname, query: parseQuery(url.search), hash: url.hash });

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div :key="char.getUserId()" class="flex flex-col divide-y rounded-lg bg-base-800 shadow-float">
        <div class="flex flex-col flex-1 p-8">
            <div class="flex flex-row items-center gap-3 mx-auto">
                <h2 class="text-2xl font-medium text-center text-neutral">
                    {{ char.getFirstname() }}, {{ char.getLastname() }}
                </h2>
                <CharSexBadge :sex="char.getSex()" />
                <div v-if="lastCharID == char.getUserId()">
                    <span
                        class="inline-flex items-center rounded-full bg-success-100 px-3 py-0.5 text-sm font-medium text-success-800">
                        {{ $t('common.last_used') }}
                    </span>
                </div>
            </div>
            <dl class="flex flex-col justify-between flex-grow mt-2 text-center">
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                            char.getJobLabel() }} ({{ $t('common.rank') }}: {{ char.getJobGradeLabel() }})</span>
                </dd>
                <dt class="text-sm text-neutral">{{ $t('common.date_of_birth') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.getDateofbirth() }}</dd>
                <dt class="text-sm text-neutral">{{ $t('common.height') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.getHeight() }}cm</dd>
                <dt class="text-sm text-neutral">{{ $t('common.visum') }}</dt>
                <dd class="text-sm text-gray-300">{{ char.getVisum() }}</dd>
                <dt class="text-sm text-neutral">{{ $t('common.playtime') }}</dt>
                <dd class="text-sm text-gray-300">{{ fromSecondsToFormattedDuration(char.getPlaytime()) }}</dd>
            </dl>
        </div>
        <div>
            <div class="flex -mt-px">
                <div class="flex flex-1 w-0">
                    <button @click="chooseCharacter()"
                        class="relative inline-flex items-center justify-center flex-1 w-0 py-4 text-sm font-semibold transition-colors border border-transparent rounded-b-lg gap-x-3 text-neutral bg-base-700 hover:bg-base-600">
                        {{ $t('common.choose') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
