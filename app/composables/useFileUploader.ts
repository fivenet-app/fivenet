// composables/useFileUploader.ts
import type { ClientStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import pica from 'pica';
import { canEncodeWebp } from '~/utils/canEncodeWebp';
import { UploadMeta, UploadPacket, type UploadResponse } from '~~/gen/ts/resources/file/filestore';

const MAX_READ_SIZE = 128 * 1024; // 128 KB

export type UploadNamespaces = 'documents' | 'qualifications' | 'qualifications-exam-questions' | 'wiki';

/**
 * Factory that returns resize + upload helpers bound to a parent record.
 */
export function useFileUploader(
    filestore: (options?: RpcOptions) => ClientStreamingCall<UploadPacket, UploadResponse>,
    namespace: UploadNamespaces,
    parentId: number,
) {
    // 1. Resize and convert to webp (if available) in browser
    const resizeImage = async (file: File): Promise<{ blob: Blob; fileName: string; mime: string }> => {
        const img = await createImageBitmap(file);

        const minW = 320; // Smallest we'll ever upscale to
        const maxW = 2250; // Largest we'll ever down-scale to

        // Decide the target width
        let tgtW = img.width;
        if (img.width < minW) {
            tgtW = minW; // Upscale small avatars
        } else if (img.width > maxW) {
            tgtW = maxW; // Shrink huge photos
        }

        // 2) derive scale & height
        const scale = tgtW / img.width; // in (0, …]
        const [w, h] = [tgtW, Math.round(img.height * scale)];

        // Prepare canvases
        const from = document.createElement('canvas');
        from.width = img.width;
        from.height = img.height;
        from.getContext('2d')!.drawImage(img, 0, 0);

        const to = document.createElement('canvas');
        to.width = w;
        to.height = h;

        await pica().resize(from, to, { unsharpAmount: 80 });

        const wantsWebp = await canEncodeWebp();
        const mime = wantsWebp ? 'image/webp' : 'image/png';
        const blob = await pica().toBlob(to, mime, 0.9);

        // Final file name (keep original base)
        const base = file.name.replace(/\.[^.]+$/, '');
        const ext = wantsWebp ? '.webp' : '.png';

        return { blob, fileName: base + ext, mime };
    };

    // 2. gRPC streaming upload
    const upload = (blob: Blob, originalName: string, mime?: string, reason?: string) =>
        new Promise<UploadResponse>((resolve, reject) => {
            const stream = filestore({});
            stream.response
                .then((resp) => resolve(resp))
                .catch((err) => {
                    if (err && (err as Error).name === 'AbortError') {
                        // Upload was aborted
                        return;
                    }
                    reject(err);
                });

            // Meta packet
            const meta = UploadMeta.create();
            meta.parentId = parentId;
            meta.namespace = namespace;
            meta.originalName = originalName;
            meta.contentType = mime || blob.type;
            meta.size = blob.size;

            meta.reason = reason ?? '';

            const metaPkt = UploadPacket.create({
                payload: {
                    oneofKind: 'meta',
                    meta: meta,
                },
            });
            stream.requests.send(metaPkt);

            // Chunk packets
            const reader = blob.stream().getReader({ mode: 'byob' });
            (async () => {
                while (true) {
                    const buffer = new Uint8Array(MAX_READ_SIZE);

                    const { value, done } = await reader.read(buffer);
                    if (done) break;
                    const pkt = UploadPacket.create({
                        payload: {
                            oneofKind: 'data',
                            data: value,
                        },
                    });
                    await stream.requests.send(pkt);
                }
                await stream.requests.complete();
            })();
        });

    const resizeAndUpload = async (file: File, reason?: string) => {
        const { blob, fileName, mime } = await resizeImage(file).catch(() => ({
            blob: file,
            fileName: file.name,
            mime: file.type || 'application/octet-stream',
        }));
        return upload(blob, fileName, mime, reason);
    };

    return { resize: resizeImage, upload, resizeAndUpload };
}
