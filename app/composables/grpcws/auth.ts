import { authUserTokenKey } from '../../stores/auth_session/constants';
import { useAuth } from '../useAuth';

export function getGrpcCharacterAuthToken(): string | null {
    return sessionStorage.getItem(authUserTokenKey);
}

export function getGrpcWebsocketAuthToken(): string | null {
    const { activeChar } = useAuth();
    if (!activeChar.value) return null;

    return sessionStorage.getItem(authUserTokenKey);
}
