// composables/useFileUploader.ts
import type { ClientStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import pica from 'pica';
import { canEncodeWebp } from '~/utils/canEncodeWebp';
import { UploadFileRequest, UploadMeta, type UploadFileResponse } from '~~/gen/ts/resources/file/filestore';

const MAX_READ_SIZE = 128 * 1024; // 128 KB

export type UploadNamespaces = 'documents' | 'jobprops' | 'qualifications' | 'qualifications-exam-questions' | 'wiki';

/**
 * Factory that returns resize + upload helpers bound to a parent record.
 */
export function useFileUploader(
    filestore: (options?: RpcOptions) => ClientStreamingCall<UploadFileRequest, UploadFileResponse>,
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

        await pica().resize(from, to, { unsharpAmount: 150 });

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
        new Promise<UploadFileResponse>((resolve, reject) => {
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

            const metaPkt = UploadFileRequest.create({
                payload: {
                    oneofKind: 'meta',
                    meta: meta,
                },
            });
            stream.requests.send(metaPkt);

            // Chunk packets
            const rs = blob.stream();

            // Try BYOB; if it fails (e.g., Chrome 103 non-byte stream), fall back to default reader.
            let reader: ReadableStreamBYOBReader | ReadableStreamDefaultReader<Uint8Array>;
            let useBYOB = false;
            try {
                // This will throw on non-byte streams (Chrome 103)
                reader = rs.getReader({ mode: 'byob' }) as ReadableStreamBYOBReader;
                useBYOB = true;
            } catch {
                reader = rs.getReader(); // default reader
                useBYOB = false;
            }

            (async () => {
                try {
                    if (useBYOB) {
                        // BYOB path
                        while (true) {
                            const buffer = new Uint8Array(MAX_READ_SIZE);
                            const { value, done } = await (reader as ReadableStreamBYOBReader).read(buffer);
                            if (done) break;
                            // `value` is a view into `buffer` (may be a subarray)
                            const pkt = UploadFileRequest.create({
                                payload: { oneofKind: 'data', data: value as Uint8Array },
                            });
                            await stream.requests.send(pkt);
                        }
                    } else {
                        // Default reader path (Chrome 103)
                        let leftover: Uint8Array | null = null;

                        while (true) {
                            const { value, done } = await (reader as ReadableStreamDefaultReader<Uint8Array>).read();
                            if (done && !leftover) break;

                            let chunk = value;
                            if (leftover) {
                                // Concatenate leftover + current chunk
                                const merged: Uint8Array = new Uint8Array(leftover.length + (chunk?.length ?? 0));
                                merged.set(leftover, 0);
                                if (chunk) merged.set(chunk, leftover.length);
                                chunk = merged;
                                leftover = null;
                            }

                            let offset = 0;
                            while (chunk && offset < chunk.length) {
                                const end = Math.min(offset + MAX_READ_SIZE, chunk.length);
                                const slice = chunk.subarray(offset, end);
                                const pkt = UploadFileRequest.create({
                                    payload: { oneofKind: 'data', data: slice },
                                });
                                await stream.requests.send(pkt);
                                offset = end;
                            }

                            // If we somehow stopped mid-chunk (shouldn’t happen), keep remainder
                            if (chunk && offset < chunk.length) {
                                leftover = chunk.subarray(offset);
                            }

                            if (done) break;
                        }
                    }
                    await stream.requests.complete();
                } catch (e) {
                    reject(e);
                } finally {
                    try {
                        reader.releaseLock?.();
                    } catch {
                        // Ignore
                    }
                }
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
