import type { Ref } from 'vue';
import type { User } from '~~/gen/ts/resources/users/user';

export function useCitizenProperties(user: Ref<User>) {
    const profileProperties = computed(() => user.value.properties?.properties ?? []);

    return { profileProperties };
}
