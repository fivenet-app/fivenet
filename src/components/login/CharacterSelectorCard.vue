<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from '../../store/store';
import { useRoute, useRouter } from 'vue-router/auto';
import { getAuthClient, handleGRPCError } from '../../grpc';
import { ChooseCharacterRequest } from '@arpanet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { parseQuery } from 'vue-router/auto';
import CharSexBadge from '../misc/CharSexBadge.vue';
import { getSecondsFormattedAsDuration } from '../../utils/time';

const store = useStore();
const route = useRoute();
const router = useRouter();

const lastCharID = computed(() => store.state.auth?.lastCharID);

const props = defineProps({
    char: {
        required: true,
        type: User,
    },
});

function chooseCharacter() {
    const req = new ChooseCharacterRequest();
    req.setCharId(props.char.getUserId());

    getAuthClient()
        .chooseCharacter(req, null)
        .then((resp) => {
            store.dispatch('auth/updateAccessToken', resp.getToken());
            store.dispatch('auth/updateActiveChar', props.char);
            store.dispatch('auth/updatePermissions', resp.getPermissionsList());
            console.log(resp.getPermissionsList());
            const path = route.query.redirect?.toString() || "/overview";
            const url = new URL("https://example.com" + path);
            router.push({ path: url.pathname, query: parseQuery(url.search), hash: url.hash });
        }).catch((err: RpcError) => {
            handleGRPCError(err);
        });
}
</script>

<template>
    <li :key="char.getUserId()"
        class="col-span-2 flex flex-col divide-y divide-white rounded-lg bg-gray-800 text-center shadow">
        <div class="flex flex-1 flex-col p-8">
            <h2 class="mt-6 text-2xl font-medium text-white">
                {{ char.getFirstname() }}, {{ char.getLastname() }}
                <CharSexBadge :sex="char.getSex()" />
            </h2>
            <dl class="mt-1 flex flex-grow flex-col justify-between">
                <dd>
                    <span v-if="lastCharID == char.getUserId()"
                        class="inline-flex items-center rounded-full bg-green-100 px-3 py-0.5 text-sm font-medium text-green-800">
                        Last Used
                    </span>
                    <br v-else />
                </dd>
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-md bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                            char.getJobLabel() }} (Rank: {{ char.getJobGradeLabel() }})</span>
                </dd>
                <dt class="text-sm text-white">Date of Birth</dt>
                <dd class="text-sm text-gray-300">{{ char.getDateofbirth() }}</dd>
                <dt class="text-sm text-white">Height</dt>
                <dd class="text-sm text-gray-300">{{ char.getHeight() }}cm</dd>
                <dt class="text-sm text-white">Visum</dt>
                <dd class="text-sm text-gray-300">{{ char.getVisum() }}</dd>
                <dt class="text-sm text-white">Playtime</dt>
                <dd class="text-sm text-gray-300">{{ getSecondsFormattedAsDuration(char.getPlaytime()) }}</dd>
            </dl>
        </div>
        <div>
            <div class="-mt-px flex divide-x divide-white">
                <div class="flex w-0 flex-1">
                    <button @click="chooseCharacter()"
                        class="relative -mr-px inline-flex w-0 flex-1 items-center justify-center gap-x-3 rounded-bl-lg border border-transparent py-4 text-sm font-semibold text-white bg-gray-600">
                        Choose
                    </button>
                </div>
            </div>
        </div>
    </li>
</template>
