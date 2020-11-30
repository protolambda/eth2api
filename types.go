package eth2api

import (
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/ztyp/view"
	"strings"
)

type GenesisResponse struct {
	GenesisTime           beacon.Timestamp `json:"genesis_time"`
	GenesisValidatorsRoot beacon.Root      `json:"genesis_validators_root"`
	GenesisForkVersion    beacon.Version   `json:"genesis_fork_version"`
}

type DepositContractResponse struct {
	// Id of Eth1 chain on which contract is deployed.
	ChainID view.Uint64View `json:"chain_id"`
	// Hex encoded deposit contract address with 0x prefix
	Address beacon.Eth1Address `json:"address"`
}

type FinalityCheckpoints struct {
	PreviousJustified beacon.Checkpoint `json:"previous_justified"`
	CurrentJustified  beacon.Checkpoint `json:"current_justified"`
	Finalized         beacon.Checkpoint `json:"finalized"`
}

type ValidatorResponse struct {
	// Index of validator in validator registry.
	Index beacon.ValidatorIndex `json:"index"`
	// Current validator balance in gwei
	Balance beacon.Gwei `json:"balance"`
	// TODO
	Status ValidatorStatus `json:"status"`
	// The validator as defined in the registry in the BeaconState
	Validator beacon.Validator `json:"validator"`
}

type ValidatorBalanceResponse struct {
	// Index of validator in validator registry.
	Index beacon.ValidatorIndex `json:"index"`

	// Current validator balance in gwei
	Balance beacon.Gwei `json:"balance"`
}

type ValidatorStatus struct {
	// TODO enum with or without additional values?
	//- pending_initialized
	//- pending_queued
	//- active_ongoing
	//- active_exiting
	//- active_slashed
	//- exited_unslashed
	//- exited_slashed
	//- withdrawal_possible
	//- withdrawal_done
}

type Committee struct {
	// Committee index at a slot
	Index beacon.CommitteeIndex `json:"index"`

	// Slot the committee performs duties on
	Slot beacon.Slot `json:"slot"`

	// List of validator indices assigned to this committee
	Validators []beacon.ValidatorIndex `json:"validators"`
}

// Network identity data, not typed in detail,
// using libp2p and discv5 dependencies for further processing is optional.
type NetworkIdentity struct {
	PeerID ApiPeerId `json:"peer_id"`

	// Ethereum node record. [Read more](https://eips.ethereum.org/EIPS/eip-778)
	ENR string `json:"enr"`

	// Node's addresses on which eth2 rpc requests are served. [Read more](https://docs.libp2p.io/reference/glossary/#multiaddr)
	P2PAddresses []string `json:"p2p_addresses"`

	// Node's addresses on which is listening for discv5 requests. [Read more](https://docs.libp2p.io/reference/glossary/#multiaddr)
	DiscoveryAddresses []string `json:"discovery_addresses"`

	// Based on eth2 Metadata object
	Metadata beacon.MetaData `json:"metadata"`
}

type Peer struct {
	PeerID             ApiPeerId           `json:"peer_id"`
	ENR                string              `json:"enr"`
	LastSeenP2PAddress string              `json:"last_seen_p2p_address"`
	State              ConnectionState     `json:"state"`
	Direction          ConnectionDirection `json:"direction"`
}

type ConnectionDirection string

const (
	ConnectionDirectionInbound  ConnectionDirection = "inbound"
	ConnectionDirectionOutbound ConnectionDirection = "outbound"
)

type ConnectionState string

const (
	ConnectionStateDisconnected  ConnectionState = "disconnected"
	ConnectionStateConnecting    ConnectionState = "connecting"
	ConnectionStateConnected     ConnectionState = "connected"
	ConnectionStateDisconnecting ConnectionState = "disconnecting"
)

// Cryptographic hash of a peer’s public key. [Read more](https://docs.libp2p.io/concepts/peer-id/)
// Further (optional) processing in Go with: https://github.com/libp2p/go-libp2p-core/blob/a39b84ea2e340466d57fdb342c7d62f12957d972/peer/peer.go#L42
type ApiPeerId string

type ErrorMessage struct {
	// Either specific error code in case of invalid request or http status code
	CodeValue uint `json:"code"`
	// Message describing error
	Message string `json:"message"`
	// Optional stacktraces, sent when node is in debug mode
	Stacktraces []string `json:"stacktraces"`
}

func (err *ErrorMessage) Code() uint {
	return err.CodeValue
}

func (err *ErrorMessage) Error() string {
	stack := ""
	if len(err.Stacktraces) > 0 {
		stack = "\n" + strings.Join(err.Stacktraces, "\n")
	}
	return fmt.Sprintf("Error(%d): %s", err.CodeValue, err.Message+stack)
}

type IndexedErrorMessageItem struct {
	// Index of item in the request list that caused the error
	Index uint `json:"index"`
	// Message describing error
	Message string `json:"message"`
}

type IndexedErrorMessage struct {
	// Either specific error code in case of invalid request or http status code
	Code uint `json:"code"`
	// Message describing error
	Message string `json:"message"`
	// List of individual items that have failed
	Failures []IndexedErrorMessageItem `json:"failures"`
}

func (m *IndexedErrorMessage) IndexedErrors() []IndexedErrorMessageItem {
	return m.Failures
}

// Optional error enhancement, applicable for API methods that can have a partial success.
type IndexedError interface {
	IndexedErrors() []IndexedErrorMessageItem
}
