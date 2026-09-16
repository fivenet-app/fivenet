---
applyTo: "gen/go/**/*.pb.go"
---

These are generated Go protobuf files. The repository intentionally generates both
normal and `protoopaque` variants, selected by Go build tags (`!protoopaque` and
`protoopaque`). This is an intentional compatibility arrangement, not an incomplete
migration.

Do not recommend converting the project to protoopaque-only, removing the
protoopaque files, or treating the generated API differences as defects. Do not
review or request edits to generated files themselves; evaluate the `.proto` sources,
Buf configuration, and handwritten Go code instead. Only report generated-code
issues when they indicate a real generator or build failure.

When a `.proto` source changes, recommend the repository workflow documented in
`CONTRIBUTING.md`: run `make fmt-proto` for formatting and `make gen-proto` to
regenerate the checked-in Go and TypeScript APIs. Do not suggest hand-editing the
generated output.
