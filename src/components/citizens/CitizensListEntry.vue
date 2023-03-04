<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { defineComponent } from 'vue';
import CitizenInfoSlideOver from './CitizenInfoSlideOver.vue';

export default defineComponent({
    methods: {
        toTitleCase(value: string) {
            return value.replace(/(?:^|\s|-)\S/g, x => x.toUpperCase());
        },
        toggleSlideOver() {
            this.open = !this.open;
        },
    },
    props: {
        "user": {
            required: true,
            type: Character,
        },
    },
    data() {
        return {
            'open': false,
        };
    },
    components: {
        CitizenInfoSlideOver,
    },
});
</script>

<template>
    <tr :key="user.getIdentifier()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
            {{ user.getFirstname() }}, {{ user.getLastname() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ toTitleCase(user.getJob()) }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getSex().toUpperCase() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getDateofbirth() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ user.getHeight() }}cm
        </td>
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div v-can="'users-findusers'">
                <button @click="toggleSlideOver()" class="text-indigo-400 hover:text-indigo-300">VIEW</button>
                <CitizenInfoSlideOver @close="toggleSlideOver()" :open="open" :user="user" />
            </div>
        </td>
    </tr>
</template>
