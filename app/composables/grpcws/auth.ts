import { authUserTokenKey } from '~/stores/auth_session';

export function getGrpcAuthToken(): string | null {
    const { activeChar } = useAuth();
    if (!activeChar.value) return null;

    return sessionStorage.getItem(authUserTokenKey);
}
