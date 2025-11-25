import { describe, expect, it } from 'vitest';
import { generateDerefURL, generateDiscordConnectURL } from './url';

describe('generateDerefURL', () => {
    it('should generate a dereferer URL with the correct target and source', () => {
        const target = 'https://example.com';
        const mockLocation = 'https://current-page.com';

        const result = generateDerefURL(target, mockLocation);

        expect(result).toBe('/dereferer?target=https%3A%2F%2Fexample.com&source=https%3A%2F%2Fcurrent-page.com');
    });
});

describe('generateDiscordConnectURL', () => {
    it('should generate a Discord connect URL with the correct provider and default parameters', () => {
        const provider = 'discord';

        const result = generateDiscordConnectURL(provider);

        expect(result).toBe('/api/oauth2/login/discord?connect-only=true');
    });

    it('should include redirect parameter if provided', () => {
        const provider = 'discord';
        const redirect = 'https://redirect-url.com';

        const result = generateDiscordConnectURL(provider, redirect);

        expect(result).toBe('/api/oauth2/login/discord?connect-only=true&redirect=https%3A%2F%2Fredirect-url.com');
    });

    it('should include additional parameters if provided', () => {
        const provider = 'discord';
        const params = { scope: 'identify', state: '12345' };

        const result = generateDiscordConnectURL(provider, undefined, params);

        expect(result).toBe('/api/oauth2/login/discord?scope=identify&state=12345&connect-only=true');
    });
});
