# Eth2 API bindings

Fully typed API bindings, for both client and server, to implement [the standard Eth2.0 API specification](https://github.com/ethereum/eth2.0-APIs).

**Work in progress, incomplete**

## How is this different from [`prysmaticlabs/ethereumapis`](https://github.com/prysmaticlabs/ethereumapis)?

- Stricter typing, bazed on [ZRNT](https://github.com/protolambda/zrnt)
- Full transport abstraction, no protobufs, implement it how you like
- Focus on full JSON compatibility with Lighthouse and Teku
- Avoids allocations, optimized requests
- Minimal dependencies
- Designed for customization. Swap the transport, change the encoding, add Eth2 Phase1, whatever. 

## License

MIT, see [`LICENSE`](./LICENSE) file.
