import { onMounted, onUnmounted, ref } from 'vue';
import type { Awareness } from 'y-protocols/awareness';

interface UserState {
    id: number; // Yjs clientId
    name: string;
    color: string;
}

export function useAwarenessUsers(awareness: Awareness) {
    const users = ref<UserState[]>([]);

    const rebuild = () => {
        const arr: UserState[] = [];

        awareness.getStates().forEach((state, id) => {
            if (state.user) {
                const { name, color } = state.user as { name: string; color: string };
                arr.push({ id, name, color });
            }
        });

        // sort alphabetically for stable UI
        arr.sort((a, b) => a.name.localeCompare(b.name));
        users.value = arr; // ref assignment triggers re-render
    };

    onMounted(() => {
        rebuild(); // initial
        awareness.on('change', rebuild);
    });

    onUnmounted(() => {
        awareness.off('change', rebuild);
    });

    return { users };
}
