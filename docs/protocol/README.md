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

Generated Go output is configured to land in `core/pkg/gen`.

## Contract Notes

Current implementation exposes one control method constant in `core/pkg/api`:

- `forward.start` (`MethodForwardStart`)

Current proto service support matches that scope:

- `ControlService.ForwardStart`

Pending methods are **not implemented yet** in code or proto for this MVP slice:

- `forward.stop`
- `forward.get_stats`

When those methods are added, keep proto RPC names and Go method constants aligned.
