import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn } from '@vueuse/core';
import { applyAwarenessUpdate, Awareness, encodeAwarenessUpdate } from 'y-protocols/awareness';
import * as Y from 'yjs';
import { ClientPacket, type ServerPacket } from '~~/gen/ts/resources/collab/collab';

export type StreamConnectFn = (options?: RpcOptions) => DuplexStreamingCall<ClientPacket, ServerPacket>;

const logger = useLogger('yjs:grpc');

interface GrpcProviderOpts {
    /** Target id (e.g., document id, page id) */
    targetId: number;

    /**
     * Reconnect function – called with {attempt}. Return delay (ms).
     * Default: exponential back-off 1 s → 32 s.
     */
    reconnectDelay?: (attempt: number) => number;
}

export default class GrpcProvider {
    public readonly yDoc: Y.Doc;
    public readonly awareness: Awareness;

    private readonly opts: GrpcProviderOpts;
    private clientId: number | undefined;
    private abort: AbortController | undefined;
    private streamConnect: StreamConnectFn;
    private stream: DuplexStreamingCall<ClientPacket, ServerPacket> | undefined;
    private connected = true;
    private reconnectAttempt = 0;
    private destroyed = false;

    constructor(doc: Y.Doc, streamProvider: StreamConnectFn, opts: GrpcProviderOpts) {
        this.yDoc = doc;
        this.opts = opts;
        this.awareness = new Awareness(doc);

        this.streamConnect = streamProvider;

        /* Setup local listeners */
        this.yDoc.on('update', this.handleDocUpdate);
        const throttledAwarenessUpdate = useThrottleFn(
            (changes: { added: number[]; updated: number[]; removed: number[] }, origin: unknown) => {
                // forward the payload unchanged
                this.handleAwarenessUpdate(changes, origin);
            },
            200, // Delay in ms
            true, // Trailing (run at end of delay window)
            true, // Leading  (run immediately on first call)
        );

        this.awareness.on('update', throttledAwarenessUpdate);

        /* Kick off first connection */
        this.connect();
    }

    /* Public helpers */
    destroy() {
        if (this.destroyed) return;

        this.destroyed = true;
        // Clear local user state
        this.awareness.setLocalState(null);

        setTimeout(() => {
            this.abort?.abort();
            this.yDoc.off('update', this.handleDocUpdate);
            this.awareness.off('update', this.handleAwarenessUpdate);
        }, 0);
    }

    /* Internal */
    private connect() {
        if (this.destroyed) return;

        if (this.abort) {
            this.abort.abort();
            this.abort = undefined;
        }

        this.abort = new AbortController();
        this.stream = this.streamConnect({
            abort: this.abort.signal,
        });

        this.sendHello();

        this.stream.responses.onError((_) => {
            this.connected = false;
            this.stream = undefined;
            this.scheduleReconnect();
        });

        this.stream.responses.onMessage((msg: ServerPacket) => {
            if (msg.msg.oneofKind === 'clientId' && !this.clientId) {
                this.clientId = msg.msg.clientId;
                return;
            }

            if (!this.clientId) {
                logger.warn('Received message before clientId was set', msg);
                return;
            }

            // Ignore our own echoes (server broadcasts them to everyone incl. sender)
            if (msg.senderId === this.clientId) return;
            logger.debug('Received message from', msg.senderId);

            switch (msg.msg.oneofKind) {
                case 'awareness': {
                    if (msg.msg.awareness.data.length > 0) {
                        applyAwarenessUpdate(this.awareness, msg.msg.awareness.data, 'remote');
                    }
                    break;
                }
                case 'yjsUpdate': {
                    if (msg.msg.yjsUpdate.data.length > 0) {
                        Y.applyUpdate(this.yDoc, msg.msg.yjsUpdate.data);
                    }
                    break;
                }
            }
        });
    }

    private scheduleReconnect() {
        if (this.destroyed) return;

        const delay = this.opts.reconnectDelay?.(this.reconnectAttempt) ?? Math.min(1000 * 2 ** this.reconnectAttempt, 32000);
        this.reconnectAttempt++;
        setTimeout(() => this.connect(), delay);
    }

    /* Yjs → Server */
    private handleDocUpdate = (update: Uint8Array) => {
        if (!this.connected) return;

        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'yjsUpdate',
                yjsUpdate: {
                    data: update,
                },
            },
        });
        this.send(msg);
    };

    private handleAwarenessUpdate = (
        { added, updated, removed }: { added: number[]; updated: number[]; removed: number[] },
        _origin: unknown,
    ) => {
        const changed = added.concat(updated, removed);
        if (changed.length === 0 || !this.connected) return;

        const update = encodeAwarenessUpdate(this.awareness, changed);
        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'awareness',
                awareness: {
                    data: update,
                },
            },
        });

        logger.debug('Send awareness update', update.length);
        this.send(msg);
    };

    private sendHello() {
        const msg = ClientPacket.create({
            msg: {
                oneofKind: 'hello',
                hello: {
                    targetId: this.opts.targetId,
                },
            },
        });
        // Optionally: send encodeStateVector / encodeStateAsUpdate to prime server
        this.send(msg);
    }

    private send(msg: ClientPacket) {
        try {
            this.stream?.requests.send(msg);
        } catch (_) {
            /* swallow if stream closed mid-send */
        }
    }
}
