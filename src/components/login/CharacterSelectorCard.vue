<script lang="ts" setup>
import { computed } from 'vue';
import { useStore } from '../../store/store';
import { useRoute, useRouter } from 'vue-router/auto';
import { getAuthClient } from '../../grpc/grpc';
import { ChooseCharacterRequest } from '@arpanet/gen/services/auth/auth_pb';
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
            console.log("Char Permissions: " + resp.getPermissionsList());
            const path = route.query.redirect?.toString() || "/overview";
            const url = new URL("https://example.com" + path);
            router.push({ path: url.pathname, query: parseQuery(url.search), hash: url.hash });
        });
}
</script>

<template>
    <div :key="char.getUserId()" class="flex flex-col divide-y rounded-lg bg-base-800">
        <div class="flex flex-1 flex-col p-8">
            <div class="flex flex-row mx-auto items-center gap-3">
                <h2 class="text-2xl font-medium text-neutral text-center">
                    {{ char.getFirstname() }}, {{ char.getLastname() }}
                </h2>
                <CharSexBadge :sex="char.getSex()" />
                <div v-if="lastCharID == char.getUserId()">
                    <span class="inline-flex items-center rounded-full bg-green-100 px-3 py-0.5 text-sm font-medium text-green-800">
                        Last Used
                    </span>
                </div>
            </div>
            <dl class="flex flex-grow flex-col justify-between mt-2 text-center">
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                            char.getJobLabel() }} (Rank: {{ char.getJobGradeLabel() }})</span>
                </dd>
                <dt class="text-sm text-neutral">Date of Birth</dt>
                <dd class="text-sm text-gray-300">{{ char.getDateofbirth() }}</dd>
                <dt class="text-sm text-neutral">Height</dt>
                <dd class="text-sm text-gray-300">{{ char.getHeight() }}cm</dd>
                <dt class="text-sm text-neutral">Visum</dt>
                <dd class="text-sm text-gray-300">{{ char.getVisum() }}</dd>
                <dt class="text-sm text-neutral">Playtime</dt>
                <dd class="text-sm text-gray-300">{{ getSecondsFormattedAsDuration(char.getPlaytime()) }}</dd>
            </dl>
        </div>
        <div>
            <div class="-mt-px flex">
                <div class="flex w-0 flex-1">
                    <button @click="chooseCharacter()"
                        class="relative inline-flex w-0 flex-1 items-center justify-center gap-x-3 rounded-b-lg border border-transparent py-4 text-sm font-semibold text-neutral bg-base-700 hover:bg-base-600 transition-colors">
                        Choose
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
