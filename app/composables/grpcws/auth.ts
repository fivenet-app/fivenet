import { useAuthSessionStore } from '~/stores/auth_session';

export function getGrpcCharacterAuthToken(): string | null {
    return useAuthSessionStore().getUserToken();
}

export function getGrpcWebsocketAuthToken(): string | null {
    return useAuthSessionStore().getUserToken();
}
