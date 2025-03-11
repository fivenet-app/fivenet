# grpcws

Based upon [GitHub improbable-eng/grpc-web](https://github.com/improbable-eng/grpc-web/tree/master/go/grpcweb) and the modifications made by [borissmidt](https://github.com/borissmidt) in [GitHub borissmidt/grpc-websocket](https://github.com/borissmidt/grpc-websocket/tree/master/go/grpcweb).

## Additional Modifications

- It has been adjusted to work with `protobuf-ts` in the context of this project.
- CORS logics have been adjusted for websockets calling (only) the "root" path.
- Registered Endpoints func is only called once during server creation.

## License

[GitHub improbable-eng/grpc-web is licensed under Apache 2.0](https://github.com/improbable-eng/grpc-web/blob/master/LICENSE.txt)
