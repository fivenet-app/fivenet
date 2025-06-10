import type { DuplexStreamingCall, RpcOptions } from '@protobuf-ts/runtime-rpc';
import { ObservableV2 } from 'lib0/observable';
import { applyAwarenessUpdate, Awareness, encodeAwarenessUpdate, removeAwarenessStates } from 'y-protocols/awareness';
import * as Y from 'yjs';
import { ClientPacket, type ServerPacket } from '~~/gen/ts/resources/collab/collab';

export type StreamConnectFn = (options?: RpcOptions) => DuplexStreamingCall<ClientPacket, ServerPacket>;

const logger = useLogger('📞 yjs:grpc');

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
    sync(synced: boolean, doc: Y.Doc): void;
    saved(): void;
};

export default class GrpcProvider extends ObservableV2<Events> {
    public readonly ydoc: Y.Doc;
    public readonly awareness: Awareness;

    private readonly opts: GrpcProviderOpts;
    private clientId: number | undefined;
    private streamConnect: StreamConnectFn;
    private stream: DuplexStreamingCall<ClientPacket, ServerPacket> | undefined;
    private connected = false;
    private reconnectAttempt = 0;
    private authorative = false;
    private synced = false;
    private destroyed = false;

    constructor(doc: Y.Doc, streamProvider: StreamConnectFn, opts: GrpcProviderOpts) {
        super();

        this.opts = opts;
        this.ydoc = doc;
        this.awareness = new Awareness(doc);

        this.streamConnect = streamProvider;

        // Setup local listeners
        this.ydoc.on('update', this.handleDocUpdate);

        this.awareness.on('update', (changes: { added: number[]; updated: number[]; removed: number[] }, origin: unknown) =>
            this.handleAwarenessUpdate(changes, origin),
        );
    }

    public get isAuthoritative() {
        return this.authorative;
    }

    // Public helpers
    override destroy() {
        if (this.destroyed) return;

        super.destroy();

        this.destroyed = true;
        // Clear local user state
        this.clientId && removeAwarenessStates(this.awareness, [this.clientId], 'app closed');

        setTimeout(() => {
            this.ydoc.off('update', this.handleDocUpdate);
            this.awareness.off('update', this.handleAwarenessUpdate);
            this.stream?.requests.complete();
        }, 0);

        logger.debug('Destroyed grpc provider');
    }

    // Internal
    public connect() {
        if (this.destroyed || this.connected) return;
        logger.info('Connecting to collab gRPC stream');

        try {
            this.stream = this.streamConnect({});
        } catch (err) {
            logger.error('Failed to connect to collab gRPC stream', err);
            this.scheduleReconnect();
            return;
        }
        this.sendHello();

        this.stream.responses.onError(() => {
            this.connected = false;
            this.clientId = undefined;

            this.stream = undefined;
            this.scheduleReconnect();
        });

        this.stream.responses.onMessage((msg: ServerPacket) => {
            if (msg.msg.oneofKind === 'handshake' && !this.clientId) {
                logger.info('Received handshake message from server', msg.msg.handshake);
                this.clientId = msg.msg.handshake.clientId;
                this.connected = true;

                if (msg.msg.handshake.first) {
                    this.authorative = true;
                    this.synced = true;

                    this.ydoc.whenSynced.then(() => this.authorative && this.emit('loadContent', []));

                    this.ydoc.emit('sync', [true, this.ydoc]);
                    this.emit('sync', [true, this.ydoc]);
                } else {
                    const sv = Y.encodeStateVector(this.ydoc);
                    this.authorative = false;

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
            logger.debug('Received message from', msg.senderId, 'oneofKind:', msg.msg.oneofKind);

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
                        const diff = Y.encodeStateAsUpdate(this.ydoc, msg.msg.syncStep.data);

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
                        Y.applyUpdate(this.ydoc, msg.msg.syncStep.data);

                        if (!this.synced) {
                            this.synced = true;
                            this.ydoc.emit('sync', [true, this.ydoc]);
                            this.emit('sync', [true, this.ydoc]);
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
                        Y.applyUpdate(this.ydoc, msg.msg.yjsUpdate.data);
                    }
                    break;
                }

                case 'targetSaved': {
                    logger.info('Received target saved message', msg.msg.targetSaved);

                    this.emit('saved', []);
                }
            }
        });

        logger.debug('Connect call completed, waiting for handshake');
    }

    private scheduleReconnect() {
        if (this.destroyed) return;
        logger.info('Scheduling reconnect', {
            reconnectAttempt: this.reconnectAttempt,
            destroyed: this.destroyed,
            connected: this.connected,
            synced: this.synced,
            authorative: this.authorative,
            clientId: this.clientId,
        });

        this.emit('sync', [false, this.ydoc]);

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

        logger.debug('Send yjs update', update.length);
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

        logger.debug('Send hello message', this.opts.targetId);
        this.send(msg);
    }

    private async send(msg: ClientPacket) {
        try {
            await this.stream?.requests.send(msg);
        } catch (_) {
            // swallow if stream closed mid-send
        }
    }
}
