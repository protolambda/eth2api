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

type ChainHead struct {
	Root beacon.Root `json:"root"`
	Slot beacon.Slot `json:"slot"`
}

type RootResponse struct {
	Root beacon.Root `json:"root"`
}

type FinalityCheckpoints struct {
	PreviousJustified beacon.Checkpoint `json:"previous_justified"`
	CurrentJustified  beacon.Checkpoint `json:"current_justified"`
	Finalized         beacon.Checkpoint `json:"finalized"`
}

type SyncingStatus struct {
	// Head slot node is trying to reach
	HeadSlot beacon.Slot `json:"head_slot"`
	// How many slots node needs to process to reach head. 0 if synced.
	SyncDistance beacon.Slot `json:"sync_distance"`
}

// Wrapper around the original ProposerDuty response
type DependentProposerDuty struct {
	// Duties are valid only on the chain with this given block root
	DependentRoot beacon.Root    `json:"dependent_root"`
	Data          []ProposerDuty `json:"data"`
}

type ProposerDuty struct {
	Pubkey beacon.BLSPubkey `json:"pubkey"`
	// Index of validator in validator registry
	ValidatorIndex beacon.ValidatorIndex `json:"validator_index"`
	// The slot at which the validator must propose block.
	Slot beacon.Slot
}

// Wrapper around the original AttesterDuty response
type DependentAttesterDuties struct {
	// Duties are valid only on the chain with this given block root
	DependentRoot beacon.Root    `json:"dependent_root"`
	Data          []AttesterDuty `json:"data"`
}

type AttesterDuty struct {
	Pubkey beacon.BLSPubkey `json:"pubkey"`
	// Index of validator in validator registry
	ValidatorIndex beacon.ValidatorIndex `json:"validator_index"`
	// The committee index
	CommitteeIndex beacon.CommitteeIndex `json:"committee_index"`
	// Number of validators in committee
	CommitteeLength view.Uint64View `json:"committee_length"`
	// Number of committees at the provided slot
	CommitteesAtSlot view.Uint64View `json:"committees_at_slot"`
	// Index of validator in committee
	ValidatorCommitteeIndex view.Uint64View `json:"validator_committee_index"`
	// The slot at which the validator must attest.
	Slot beacon.Slot `json:"slot"`
}

type BeaconCommitteeSubscribeSignal struct {
	ValidatorIndex beacon.ValidatorIndex `json:"validator_index"`
	CommitteeIndex beacon.CommitteeIndex `json:"committee_index"`
	// Number of committees at the returned slot
	CommitteesAtSlot view.Uint64View `json:"committees_at_slot"`
	// Should be slot at which validator is assigned to attest
	Slot beacon.Slot `json:"slot"`
	// Signals to BN that a validator on the VC has been chosen for aggregator role.
	IsAggregator view.BoolView `json:"is_aggregator"`
}

type NodeVersionResponse struct {
	Version string `json:"version"`
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

// TODO enum with or without additional values?
type ValidatorStatus string

// status names known from lighthouse, differ from spec, TODO
const (
	ValidatorStatusUnknown                     ValidatorStatus = "unknown"
	ValidatorStatusWaitingForEligibility       ValidatorStatus = "waiting_for_eligibility"
	ValidatorStatusWaitingForFinality          ValidatorStatus = "waiting_for_finality"
	ValidatorStatusWaitingInQueue              ValidatorStatus = "waiting_in_queue"
	ValidatorStatusStandbyForActive            ValidatorStatus = "standby_for_active"
	ValidatorStatusActive                      ValidatorStatus = "active"
	ValidatorStatusActiveAwaitingVoluntaryExit ValidatorStatus = "active_awaiting_voluntary_exit"
	ValidatorStatusActiveAwaitingSlashedExit   ValidatorStatus = "active_awaiting_slashed_exit"
	ValidatorStatusExitedVoluntarily           ValidatorStatus = "exited_voluntarily"
	ValidatorStatusExitedSlashed               ValidatorStatus = "exited_slashed"
	ValidatorStatusWithdrawable                ValidatorStatus = "withdrawable"
	ValidatorStatusWithdrawn                   ValidatorStatus = "withdrawn"
)

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

// Retrieves number of known peers.
type PeerCountResponse struct {
	Disconnected  view.Uint64View `json:"disconnected"`
	Connecting    view.Uint64View `json:"connecting"`
	Connected     view.Uint64View `json:"connected"`
	Disconnecting view.Uint64View `json:"disconnecting"`
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
