import slugify from 'slugify';

slugify.extend({
    '.': '-',
    '"': '-',
    '”': '-',
    '“': '-',
    '„': '-',
    '⹂': '-',
    "'": '-',
    '!': '-',
    '?': '-',
    '&': 'and',
});

export default function slug(input: string): string {
    return slugify(input, {
        lower: true,
        trim: true,
    });
}
