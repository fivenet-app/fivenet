import { authUserTokenKey } from '~/stores/auth_session';

export function getGrpcRpcAuthToken(): string | null {
    return sessionStorage.getItem(authUserTokenKey);
}

export function getGrpcWebsocketAuthToken(): string | null {
    const { activeChar } = useAuth();
    if (!activeChar.value) return null;

    return sessionStorage.getItem(authUserTokenKey);
}
