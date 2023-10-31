// https://github.com/timostamm/protobuf-ts/pull/233#issuecomment-1289053379
// From @smnbbrv on GitHub
import fs from 'fs';
import glob from 'glob';

const protoRoot = './gen/ts';

glob(protoRoot + '/**/*.ts', async (_, files) => {
    files.forEach((file) => {
        let content = fs.readFileSync(file, 'utf-8');

        content = content
            .split('\n')
            .map((s) => s.replace(/^(import .+? from ["']\..+?)(["'];)$/, '$1.js$2'))
            .join('\n');

        fs.writeFileSync(file, content, 'utf-8');
    });
});
