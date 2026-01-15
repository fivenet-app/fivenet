import { describe, expect, it } from 'vitest';
import type { File } from '~~/gen/ts/resources/file/file';
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

describe('useImageURL', () => {
    it('should return a broken image URL when the filePath is undefined and no fallback is set', () => {
        const filePath = ref(undefined);
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe(undefined);
    });

    it('should return a broken image URL when the filePath is undefined and fallback is set', () => {
        const filePath = ref(undefined);
        const imageURL = useImageURL(filePath, '/images/broken_image.png');

        expect(imageURL.value).toBe('/images/broken_image.png');
    });

    it('should return the resolved path when filePath is a valid object', () => {
        const filePath = ref<File>({
            id: 0,
            filePath: 'example/path/to/image.png',
            isDir: false,
            byteSize: 0,
            contentType: '',
        });
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/filestore/example/path/to/image.png');
    });

    it('should return the resolved path when filePath is a string', () => {
        const filePath = ref('example/path/to/image.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/filestore/example/path/to/image.png');
    });

    it('should return the original path when filePath starts with http', () => {
        const filePath = ref('http://example.com/image.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/image_proxy/http%3A%2F%2Fexample.com%2Fimage.png');
    });

    it('should return the original path when filePath starts with /images', () => {
        const filePath = ref('/images/example.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/images/example.png');
    });

    it('should return the original path when filePath starts with /api/filestore', () => {
        const filePath = ref('/api/filestore/example.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/filestore/example.png');
    });

    it('should return the original path when filePath starts with /api/image_proxy', () => {
        const filePath = ref('/api/image_proxy/example.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/image_proxy/example.png');
    });

    it('should return the original "fixed" path when filePath starts with /api/filestore but is missing a slash', () => {
        const filePath = ref('/api/filestoreexample/path/to/image.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/filestore/example/path/to/image.png');
    });

    it('should update the imageURL when filePath changes', async () => {
        const filePath = ref('example/path/to/image.png');
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe('/api/filestore/example/path/to/image.png');

        filePath.value = 'http://example.com/new-image.png';
        await nextTick();

        expect(imageURL.value).toBe('/api/image_proxy/http%3A%2F%2Fexample.com%2Fnew-image.png');
    });

    it('should return the original data base64 url', () => {
        const filePath = ref(
            // Taken from "ashleedawg" on <https://gist.github.com/ondrek/7413434?permalink_comment_id=3914686#gistcomment-3914686>
            // Thanks for the tiny image! :-)
            'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII=',
        );
        const imageURL = useImageURL(filePath);

        expect(imageURL.value).toBe(
            'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII=',
        );
    });
});
