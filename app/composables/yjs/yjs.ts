import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { ObservableV2 } from 'lib0/observable';
import { applyAwarenessUpdate, Awareness, encodeAwarenessUpdate, removeAwarenessStates } from 'y-protocols/awareness';
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

type Events = {
    loadContent(): void;
};

export default class GrpcProvider extends ObservableV2<Events> {
    public readonly yDoc: Y.Doc;
    public readonly awareness: Awareness;

    private readonly opts: GrpcProviderOpts;
    private clientId: number | undefined;
    private streamConnect: StreamConnectFn;
    private stream: DuplexStreamingCall<ClientPacket, ServerPacket> | undefined;
    private connected = true;
    private reconnectAttempt = 0;
    private synced = false;
    private destroyed = false;

    constructor(doc: Y.Doc, streamProvider: StreamConnectFn, opts: GrpcProviderOpts) {
        super();

        this.opts = opts;
        this.yDoc = doc;
        this.awareness = new Awareness(doc);

        this.streamConnect = streamProvider;

        /* Setup local listeners */
        this.yDoc.on('update', this.handleDocUpdate);

        this.awareness.on('update', (changes: { added: number[]; updated: number[]; removed: number[] }, origin: unknown) =>
            this.handleAwarenessUpdate(changes, origin),
        );

        // Kick off connection
        this.connect();
    }

    // Public helpers
    override destroy() {
        if (this.destroyed) return;

        this.destroyed = true;
        // Clear local user state
        this.clientId && removeAwarenessStates(this.awareness, [this.clientId], 'app closed');

        setTimeout(() => {
            this.yDoc.off('update', this.handleDocUpdate);
            this.awareness.off('update', this.handleAwarenessUpdate);
        }, 0);

        super.destroy();
    }

    // Internal
    private connect() {
        if (this.destroyed) return;

        this.stream = this.streamConnect({});
        this.sendHello();

        this.connected = true;

        this.stream.responses.onError((_) => {
            this.connected = false;
            this.stream = undefined;
            this.scheduleReconnect();
        });

        this.stream.responses.onMessage((msg: ServerPacket) => {
            if (msg.msg.oneofKind === 'handshake' && !this.clientId) {
                logger.info('Received handshake message from server', msg.msg.handshake);

                this.clientId = msg.msg.handshake.clientId;
                if (msg.msg.handshake.first) {
                    this.emit('loadContent', []);

                    this.synced = true;
                    this.yDoc.emit('sync', [true, this.yDoc]);
                } else {
                    const sv = Y.encodeStateVector(this.yDoc);

                    this.send(
                        ClientPacket.create({
                            msg: {
                                oneofKind: 'syncStep',
                                syncStep: {
                                    step: 1,
                                    data: sv,
                                },
                            },
                        }),
                    );
                }

                return;
            }

            if (!this.clientId) {
                logger.warn('Received message before clientId was set', msg);
                return;
            }

            // Ignore our own echoes (server broadcasts them to everyone incl. sender)
            if (msg.senderId === this.clientId) return;
            logger.debug('Received message from', msg.senderId, ' - oneofKind: ', msg.msg.oneofKind);

            switch (msg.msg.oneofKind) {
                case 'syncStep': {
                    if (msg.msg.syncStep.data.length === 0 || msg.msg.syncStep.step < 1 || msg.msg.syncStep.step > 2) {
                        logger.warn('Received invalid sync step', msg.msg.syncStep);
                        break;
                    }

                    logger.debug(
                        'Received sync step',
                        msg.msg.syncStep.step,
                        'from',
                        msg.senderId,
                        'length',
                        msg.msg.syncStep.data.length,
                    );
                    if (msg.msg.syncStep.step === 1) {
                        const diff = Y.encodeStateAsUpdate(this.yDoc, msg.msg.syncStep.data);

                        this.send({
                            msg: {
                                oneofKind: 'syncStep',
                                syncStep: {
                                    step: 2,
                                    data: diff,
                                    receiverId: msg.senderId,
                                },
                            },
                        });
                    } else if (msg.msg.syncStep.step === 2) {
                        Y.applyUpdate(this.yDoc, msg.msg.syncStep.data);

                        if (!this.synced) {
                            this.synced = true;
                            this.yDoc.emit('sync', [true, this.yDoc]);
                        }
                    }

                    break;
                }

                case 'awareness': {
                    if (msg.msg.awareness.data.length > 0) {
                        applyAwarenessUpdate(this.awareness, msg.msg.awareness.data, 'remote');
                    }
                    break;
                }
                case 'yjsUpdate': {
                    if (msg.msg.yjsUpdate.data.length > 0) {
                        logger.debug('Received Yjs update', msg.msg.yjsUpdate.data.length, 'from', msg.senderId);
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

    // Yjs to Server
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

        logger.debug('Send awareness update', update.length);
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
        this.send(msg);
        // Optionally: send encodeStateVector / encodeStateAsUpdate to prime server
    }

    private send(msg: ClientPacket) {
        try {
            this.stream?.requests.send(msg);
        } catch (_) {
            // swallow if stream closed mid-send
        }
    }
}
