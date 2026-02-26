# XShare Protocol (MVP)

This directory defines the protobuf contracts used by the MVP control/data/OTA channels.

## Layout

- `protocol/buf.yaml`: Buf module + lint/breaking config.
- `protocol/buf.gen.yaml`: code generation plugins and outputs.
- `protocol/proto/**`: source `.proto` files.

## Prerequisites

- `buf` CLI installed and on `PATH`.
- Network access for remote Buf plugins configured in `buf.gen.yaml`.

## Generate Code

From repository root:

```bash
bash tools/proto-gen.sh
```

or directly:

```bash
cd protocol
buf generate
```

Generated Go output is configured to land in `core/go/pkg/gen`.

## Contract Notes

MVP control method names are defined in Go constants under `core/go/pkg/api` and include:

- `forward.start`
- `forward.stop`
- `forward.get_stats`

Keep proto RPC names and these method constants aligned.
