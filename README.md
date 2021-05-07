# Eth2 API bindings

Fully typed API bindings, for both client and server, to implement [the standard Eth2.0 API specification](https://github.com/ethereum/eth2.0-APIs).

**Work in progress, testing in progress**

TODO:
- [ ] Client
  - [x] Types for full API spec
  - [x] Bindings for full API spec
      - [x] Beacon API
      - [x] Debug API
      - [x] Config API
      - [x] Node API
      - [x] Validator API
  - [x] Abstraction of requests/responses
  - [x] HTTP client implementation
  - [ ] Testing: API Integration test-suite against test vectors (generated from Lighthouse API, verified with spec)
    - [x] Beacon API
    - [ ] Debug API
    - [ ] Config API
    - [ ] Node API
    - [ ] Validator API
  - [ ] Tests for the util methods
- [ ] Server
  - [ ] (WIP) Interfaces for serverside API
  - [ ] (WIP) Abstract server that consumes above interfaces, runs API server

The API design is not definite yet, current bindings are based on Eth2.0-apis commit `ceb555f9b40ff9c2094d038e9f70a19419e5b652`.

## Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/protolambda/zrnt/eth2/configs"
    "github.com/protolambda/ztyp/tree"
    "net/http"
    "time"
    "github.com/protolambda/zrnt/eth2/beacon/common"
    "github.com/protolambda/zrnt/eth2/beacon/phase0"
    "github.com/protolambda/eth2api"
    "github.com/protolambda/eth2api/beaconapi"
)

func main() {
    // Make an HTTP client (reuse connections!)
    client := &eth2api.HttpClient{
        Addr: "http://localhost:5052",
        Cli: &http.Client{
            Transport: &http.Transport{
                MaxIdleConnsPerHost: 123,
            },
            Timeout: 10 * time.Second,
        },
    }

    // e.g. cancel requests on demand if you don't need the block anymore.
    ctx, cancel := context.WithCancel(context.Background())

    slot := common.Slot(127) // strict Eth2 types from ZRNT fully integrated
    // or try eth2api.BlockIdRoot(beacon.Root{0x.....}), eth2api.BlockHead, eth2api.BlockGenesis, etc. as BlockId
  
    // standard errors are part of the API.
    if blockEnvelop, err := beaconapi.Block(ctx, client, eth2api.BlockIdSlot(slot)); blockEnvelop == nil {
        fmt.Println("block not found")
    } else if err != nil {
    	fmt.Println("failed to get block", err)
    } else {
        // Easy access to optimized Eth2 spec functions 
        blockRoot := blockEnvelop.SignedBlock.(*phase0.SignedBeaconBlock).Message.HashTreeRoot(configs.Mainnet, tree.GetHashFn())
        // Or just use the block envelop fields, the same between all Eth2 forks
        blockRoot = blockEnvelop.BlockRoot
        fmt.Println("got block: ", blockRoot)
    }
}
```

## Testing

Testing is not fully automated yet, awaiting first test vector release.

For now, you can:
- Use [`protolambda/eth2-api-testgen`](https://github.com/protolambda/eth2-api-testgen) to generate test vectors.
- Copy the `output` dir to the `tests` dir in this repo. (`cp -r ../eth2-api-testgen/output tests`)
- Run the Go tests in this repo (`go test ./...`)

## Architecture

- Strictly typed client bindings that send a `PreparedRequest` via a `Client` and decode the `Response`
- Strictly typed server routes, backed by chain/db interfaces, with handlers which take a `Request` and produce a `PreparedResponse`
- A `Server` is a mux of routes, calling the handlers, and encoding the `PreparedResponse` to serve to the requesting client

The `Client` and `Server` essentially do the encoding and transport, and are fully replaceable, without rewriting any API.

```
                                    __________________ Client __________________
                                   |                                            |
call binding ---PreparedRequest---> Eth2HttpClient ---http.Request--> HttpClient

                          _______________________ Server _________________________
                         |                                                        |
create route ---Route---> Eth2HttpMux ---http.Handler---> julienschmidt/httprouter


  ________________ Server ________________                                                     _________________ Server _______________ 
 |                                        |                                                   |                                        |
 ---http.Request---> Eth2HttpMux & Decoder ---Request---> route handler ---PreparedResponse---> Http server Encoder ---http.Response--->

```

## How is this different from [`prysmaticlabs/ethereumapis`](https://github.com/prysmaticlabs/ethereumapis)?

- Stricter typing, bazed on [ZRNT](https://github.com/protolambda/zrnt)
- Full transport abstraction, no protobufs, implement it how you like
- Focus on full JSON compatibility with Lighthouse and Teku
- Avoids allocations, optimized requests
- Minimal dependencies
- Designed for customization. Swap the transport, change the encoding, add Eth2 Phase1, whatever. 

## License

MIT, see [`LICENSE`](./LICENSE) file.
