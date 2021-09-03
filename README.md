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
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/client/beaconapi"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/configs"
	"net/http"
	"os"
	"time"
)

func main() {
	// Make an HTTP client (reuse connections!)
	client := &eth2api.Eth2HttpClient{
		Addr: "http://localhost:5052",
		Cli: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 123,
			},
			Timeout: 10 * time.Second,
		},
	  Codec: eth2api.JSONCodec{},
	}


	//// e.g. cancel requests with a context.WithTimeout/WithCancel/WithDeadline
	ctx := context.Background()

	var genesis eth2api.GenesisResponse
	if exists, err := beaconapi.Genesis(ctx, client, &genesis); !exists {
		fmt.Println("chain did not start yet")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("failed to get genesis", err)
		os.Exit(1)
	}
	spec := configs.Mainnet
	// or load testnet config info from a YAML file
	// yaml.Unmarshal(data, &spec.Config)

	// every fork has a digest. Blocks are versioned by name in the API,
	// but wrapped with digest info in ZRNT to do enable different kinds of processing
	altairForkDigest := common.ComputeForkDigest(spec.ALTAIR_FORK_VERSION, genesis.GenesisValidatorsRoot)


	id := eth2api.BlockHead
	// or try other block ID types:
	// eth2api.BlockIdSlot(common.Slot(12345))
	// eth2api.BlockIdRoot(common.Root{0x.....})
	// eth2api.BlockGenesis

	// standard errors are part of the API.
	var versionedBlock eth2api.VersionedSignedBeaconBlock
	if exists, err := beaconapi.BlockV2(ctx, client, id, &versionedBlock); !exists {
		fmt.Println("block not found")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("failed to get block", err)
		os.Exit(1)
	} else {
		fmt.Println("block version:", versionedBlock.Version)
		// get the block (any fork)
		fmt.Printf("data: %v\n", versionedBlock.Data)

		// add digest:
		envelope := versionedBlock.Data.Envelope(spec, altairForkDigest)
		fmt.Println("got block:", envelope.BlockRoot)
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
