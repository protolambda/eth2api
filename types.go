package eth2api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/ztyp/view"
	"strings"
)

type GenesisResponse struct {
	GenesisTime           common.Timestamp `json:"genesis_time"`
	GenesisValidatorsRoot common.Root      `json:"genesis_validators_root"`
	GenesisForkVersion    common.Version   `json:"genesis_fork_version"`
}

type DepositContractResponse struct {
	// Id of Eth1 chain on which contract is deployed.
	ChainID view.Uint64View `json:"chain_id"`
	// Hex encoded deposit contract address with 0x prefix
	Address common.Eth1Address `json:"address"`
}

type ChainHead struct {
	Root common.Root `json:"root"`
	Slot common.Slot `json:"slot"`
}

type RootResponse struct {
	Root common.Root `json:"root"`
}

type FinalityCheckpoints struct {
	PreviousJustified common.Checkpoint `json:"previous_justified"`
	CurrentJustified  common.Checkpoint `json:"current_justified"`
	Finalized         common.Checkpoint `json:"finalized"`
}

type BeaconBlockHeaderAndInfo struct {
	Root      common.Root                    `json:"root"`
	Canonical bool                           `json:"canonical"`
	Header    common.SignedBeaconBlockHeader `json:"header"`
}

type SyncingStatus struct {
	// Head slot node is trying to reach
	HeadSlot common.Slot `json:"head_slot"`
	// How many slots node needs to process to reach head. 0 if synced.
	SyncDistance common.Slot `json:"sync_distance"`
}

// Wrapper around the original ProposerDuty response
type DependentProposerDuty struct {
	// Duties are valid only on the chain with this given block root
	DependentRoot common.Root    `json:"dependent_root"`
	Data          []ProposerDuty `json:"data"`
}

type ProposerDuty struct {
	Pubkey common.BLSPubkey `json:"pubkey"`
	// Index of validator in validator registry
	ValidatorIndex common.ValidatorIndex `json:"validator_index"`
	// The slot at which the validator must propose block.
	Slot common.Slot
}

type SyncCommitteeDuty struct {
	Pubkey                        common.BLSPubkey      `json:"pubkey"`
	ValidatorIndex                common.ValidatorIndex `json:"validator_index"`
	ValidatorSyncCommitteeIndices []view.Uint64View     `json:"validator_sync_committee_indices"`
}

// Wrapper around the original AttesterDuty response
type DependentAttesterDuties struct {
	// Duties are valid only on the chain with this given block root
	DependentRoot common.Root    `json:"dependent_root"`
	Data          []AttesterDuty `json:"data"`
}

type AttesterDuty struct {
	Pubkey common.BLSPubkey `json:"pubkey"`
	// Index of validator in validator registry
	ValidatorIndex common.ValidatorIndex `json:"validator_index"`
	// The committee index
	CommitteeIndex common.CommitteeIndex `json:"committee_index"`
	// Number of validators in committee
	CommitteeLength view.Uint64View `json:"committee_length"`
	// Number of committees at the provided slot
	CommitteesAtSlot view.Uint64View `json:"committees_at_slot"`
	// Index of validator in committee
	ValidatorCommitteeIndex view.Uint64View `json:"validator_committee_index"`
	// The slot at which the validator must attest.
	Slot common.Slot `json:"slot"`
}

type BeaconCommitteeSubscribeSignal struct {
	ValidatorIndex common.ValidatorIndex `json:"validator_index"`
	CommitteeIndex common.CommitteeIndex `json:"committee_index"`
	// Number of committees at the returned slot
	CommitteesAtSlot view.Uint64View `json:"committees_at_slot"`
	// Should be slot at which validator is assigned to attest
	Slot common.Slot `json:"slot"`
	// Signals to BN that a validator on the VC has been chosen for aggregator role.
	IsAggregator view.BoolView `json:"is_aggregator"`
}

type SyncCommitteeSubscribeSignal struct {
	ValidatorIndex common.ValidatorIndex `json:"validator_index"`
	// positions of the validator-index relative to the complete sync committee (a validator can have multiple)
	SyncCommitteeIndices []view.Uint64View `json:"sync_committee_indices"`
	UntilEpoch           common.Epoch      `json:"until_epoch"`
}

type NodeVersionResponse struct {
	Version string `json:"version"`
}

type ValidatorResponse struct {
	// Index of validator in validator registry.
	Index common.ValidatorIndex `json:"index"`
	// Current validator balance in gwei
	Balance common.Gwei `json:"balance"`
	// TODO
	Status ValidatorStatus `json:"status"`
	// The validator as defined in the registry in the BeaconState
	Validator phase0.Validator `json:"validator"`
}

type ValidatorBalanceResponse struct {
	// Index of validator in validator registry.
	Index common.ValidatorIndex `json:"index"`

	// Current validator balance in gwei
	Balance common.Gwei `json:"balance"`
}

// TODO enum with or without additional values?
type ValidatorStatus string

func (vs ValidatorStatus) String() string {
	return string(vs)
}

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
	Index common.CommitteeIndex `json:"index"`

	// Slot the committee performs duties on
	Slot common.Slot `json:"slot"`

	// List of validator indices assigned to this committee
	Validators []common.ValidatorIndex `json:"validators"`
}

type SyncCommittees struct {
	// All of the validator indices in the current sync committee
	Validators []common.ValidatorIndex `json:"validators"`
	// Subcommittee slices of the current sync committee
	ValidatorAggregates [][]common.ValidatorIndex `json:"validator_aggregates"`
}

type versionStruct struct {
	Version string `json:"version"`
}

type blockDataStruct struct {
	Data common.SpecObj `json:"data"`
}

type VersionedBeaconBlock struct {
	Version string `json:"version"`
	// Data is *phase0.BeaconBlock, *altair.BeaconBlock,
	// *merge.BeaconBlock or *sharding.BeaconBlock.
	Data common.SpecObj `json:"data"`
}

func (v *VersionedBeaconBlock) UnmarshalJSON(b []byte) error {
	var version versionStruct
	if err := json.Unmarshal(b, &version); err != nil {
		return err
	}
	var data blockDataStruct
	switch version.Version {
	case "phase0":
		data.Data = new(phase0.BeaconBlock)
	case "altair":
		data.Data = new(altair.BeaconBlock)
	case "merge":
		data.Data = new(merge.BeaconBlock)
	case "sharding":
		data.Data = new(sharding.BeaconBlock)
	default:
		return fmt.Errorf("unrecognized version: %q", version.Version)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	v.Data = data.Data
	v.Version = version.Version
	return nil
}

type SignedBeaconBlock interface {
	common.SpecObj
	common.EnvelopeBuilder
}

type signedBlockDataStruct struct {
	Data SignedBeaconBlock `json:"data"`
}

type VersionedSignedBeaconBlock struct {
	Version string `json:"version"`
	// Data is *phase0.SignedBeaconBlock, *altair.SignedBeaconBlock,
	// *merge.SignedBeaconBlock or *sharding.SignedBeaconBlock.
	Data SignedBeaconBlock `json:"data"`
}

func (v *VersionedSignedBeaconBlock) UnmarshalJSON(b []byte) error {
	var version versionStruct
	if err := json.Unmarshal(b, &version); err != nil {
		return err
	}
	var data signedBlockDataStruct
	switch version.Version {
	case "phase0":
		data.Data = new(phase0.SignedBeaconBlock)
	case "altair":
		data.Data = new(altair.SignedBeaconBlock)
	case "merge":
		data.Data = new(merge.SignedBeaconBlock)
	case "sharding":
		data.Data = new(sharding.SignedBeaconBlock)
	default:
		return fmt.Errorf("unrecognized version: %q", version.Version)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	v.Data = data.Data
	v.Version = version.Version
	return nil
}

type stateDataStruct struct {
	Data common.SpecObj `json:"data"`
}

type VersionedBeaconState struct {
	Version string `json:"version"`
	// Data is *phase0.BeaconState, *altair.BeaconState, *merge.BeaconState or *sharding.BeaconState.
	// See the Tree(spec) method to transform into a binary-tree backed state for advanced processing.
	Data common.SpecObj `json:"data"`
}

func (v *VersionedBeaconState) Tree(spec *common.Spec) (common.BeaconState, error) {
	if v.Data == nil {
		return nil, fmt.Errorf("no state (version: %q)", v.Version)
	}
	var buf bytes.Buffer
	w := codec.NewEncodingWriter(&buf)
	if err := v.Data.Serialize(spec, w); err != nil {
		return nil, err
	}
	data := buf.Bytes()
	r := codec.NewDecodingReader(bytes.NewReader(data), uint64(len(data)))
	switch v.Version {
	case "phase0":
		return phase0.AsBeaconStateView(phase0.BeaconStateType(spec).Deserialize(r))
	case "altair":
		return altair.AsBeaconStateView(altair.BeaconStateType(spec).Deserialize(r))
	case "merge":
		return merge.AsBeaconStateView(merge.BeaconStateType(spec).Deserialize(r))
	case "sharding":
		return sharding.AsBeaconStateView(sharding.BeaconStateType(spec).Deserialize(r))
	default:
		return nil, fmt.Errorf("unrecognized version: %q", v.Version)
	}
}

func (v *VersionedBeaconState) UnmarshalJSON(b []byte) error {
	var version versionStruct
	if err := json.Unmarshal(b, &version); err != nil {
		return err
	}
	var data stateDataStruct
	switch version.Version {
	case "phase0":
		data.Data = new(phase0.BeaconState)
	case "altair":
		data.Data = new(altair.BeaconState)
	case "merge":
		data.Data = new(merge.BeaconState)
	case "sharding":
		data.Data = new(sharding.BeaconState)
	default:
		return fmt.Errorf("unrecognized version: %q", version.Version)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	v.Data = data.Data
	v.Version = version.Version
	return nil
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
	Metadata common.MetaData `json:"metadata"`
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
