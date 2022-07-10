package beaconapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/pool"
	"strconv"
)

// Serves attestations known by the node but not necessarily incorporated into any block
func PoolAttestations(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/pool/attestations",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var opts []pool.AttSearchOption

			slotFilterVal, ok := req.Query("slot")
			if ok && len(slotFilterVal) != 0 {
				slot, err := strconv.ParseUint(slotFilterVal[0], 0, 64)
				if err != nil {
					return eth2api.RespondBadInput(err)
				}
				opts = append(opts, pool.WithSlot(common.Slot(slot)))
			}

			commIndexFilterVal, ok := req.Query("committee_index")
			if ok && len(commIndexFilterVal) != 0 {
				index, err := strconv.ParseUint(commIndexFilterVal[0], 0, 64)
				if err != nil {
					return eth2api.RespondBadInput(err)
				}
				opts = append(opts, pool.WithCommittee(common.CommitteeIndex(index)))
			}
			atts := backend.AttestationPool.Search(opts...)
			return eth2api.RespondOK(eth2api.Wrap(atts))
		})
}

// Handles publishing of attestations, stores them in the pool and sends them to the publisher.
func PublishAttestations(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/pool/attestations",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var atts []phase0.Attestation
			if err := req.DecodeBody(&atts); err != nil {
				return eth2api.RespondBadInput(err)
			}
			// TODO: test performance. This can probably be optimized.
			// Maybe even propagate before checking validity within own pool
			handleAtt := func(att *phase0.Attestation) error {
				targetSlot, err := backend.Spec.EpochStartSlot(att.Data.Target.Epoch)
				if err != nil {
					return err
				}
				start, ok := backend.Chain.ByBlockSlot(att.Data.Target.Root, targetSlot)
				if !ok {
					return fmt.Errorf("cannot get chain element at target %s for committee retrieval", att.Data.Target)
				}
				epc, err := start.EpochsContext(ctx)
				if err != nil {
					return err
				}
				comm, err := epc.GetBeaconCommittee(att.Data.Slot, att.Data.Index)
				if err != nil {
					return err
				}
				if err := backend.AttestationPool.AddAttestation(ctx, att, comm); err != nil {
					return err
				}
				return backend.Publisher.PublishAttestation(ctx, att)
			}
			var ierrs []eth2api.IndexedErrorMessageItem
			for i := range atts {
				if err := handleAtt(&atts[i]); err != nil {
					ierrs = append(ierrs, eth2api.IndexedErrorMessageItem{
						Index:   2,
						Message: err.Error(),
					})
				}
			}
			if len(ierrs) > 0 {
				return eth2api.RespondBadInputs("cannot store/publish", ierrs)
			}
			return eth2api.RespondOKMsg("stored and broadcast all attestations")
		})
}

// Retrieves attester slashings known by the node but not necessarily incorporated into any block
func PoolAttesterSlashings(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/pool/attester_slashings",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			return eth2api.RespondOK(eth2api.Wrap(backend.AttesterSlashingPool.All()))
		})
}

// Submits AttesterSlashing object to node's pool and if passes validation node MUST broadcast it to network.
func PublishAttesterSlashing(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/pool/attester_slashings",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var attSlashing phase0.AttesterSlashing
			if err := req.DecodeBody(&attSlashing); err != nil {
				return eth2api.RespondBadInput(err)
			}
			if err := backend.AttesterSlashingPool.AddAttesterSlashing(ctx, &attSlashing); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to add attester slashing to pool: %v", err))
			}
			if err := backend.Publisher.PublishAttesterSlashing(ctx, &attSlashing); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to publish attester slashing: %v", err))
			}
			return eth2api.RespondOK(eth2api.Wrap(backend.AttesterSlashingPool.All()))
		})
}

// Retrieves proposer slashings known by the node but not necessarily incorporated into any block
func PoolProposerSlashings(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/pool/proposer_slashings",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			return eth2api.RespondOK(eth2api.Wrap(backend.ProposerSlashingPool.All()))
		})
}

// Submits ProposerSlashing object to node's pool and if passes validation node MUST broadcast it to network.
func PublishProposerSlashing(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/pool/proposer_slashings",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var propSlashing phase0.ProposerSlashing
			if err := req.DecodeBody(&propSlashing); err != nil {
				return eth2api.RespondBadInput(err)
			}
			if err := backend.ProposerSlashingPool.AddProposerSlashing(ctx, &propSlashing); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to add proposer slashing to pool: %v", err))
			}
			if err := backend.Publisher.PublishProposerSlashing(ctx, &propSlashing); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to publish proposer slashing: %v", err))
			}
			return eth2api.RespondOK(eth2api.Wrap(backend.ProposerSlashingPool.All()))
		})
}

// Submits sync committee signature objects to the node.
// Sync committee signatures are not present in phase0, but are required for Altair networks.
// If a sync committee signature is validated successfully the node MUST publish that sync committee signature on all applicable subnets.
// If one or more sync committee signatures fail validation the node MUST return a 400 error with details of which sync committee signatures have failed, and why.
func PublishSyncCommittees(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/pool/sync_committees",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var atts []altair.SyncCommitteeMessage
			if err := req.DecodeBody(&atts); err != nil {
				return eth2api.RespondBadInput(err)
			}
			handleMsg := func(msg *altair.SyncCommitteeMessage) error {
				if err := backend.SyncCommitteePool.AddSyncCommitteeMessage(ctx, msg); err != nil {
					return err
				}
				return backend.Publisher.PublishSyncCommitteeMessage(ctx, msg)
			}
			var ierrs []eth2api.IndexedErrorMessageItem
			for i := range atts {
				if err := handleMsg(&atts[i]); err != nil {
					ierrs = append(ierrs, eth2api.IndexedErrorMessageItem{
						Index:   2,
						Message: err.Error(),
					})
				}
			}
			if len(ierrs) > 0 {
				return eth2api.RespondBadInputs("cannot store/publish", ierrs)
			}
			return eth2api.RespondOKMsg("stored and broadcast all sync committee messages")
		})
}

// Retrieves voluntary exits known by the node but not necessarily incorporated into any block
func PoolVoluntaryExits(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/pool/voluntary_exits",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			return eth2api.RespondOK(eth2api.Wrap(backend.VoluntaryExitPool.All()))
		})
}

// Submits SignedVoluntaryExit object to node's pool and if passes validation node MUST broadcast it to network.
func PublishVoluntaryExit(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/pool/voluntary_exits",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var signedExit phase0.SignedVoluntaryExit
			if err := req.DecodeBody(&signedExit); err != nil {
				return eth2api.RespondBadInput(err)
			}
			if err := backend.VoluntaryExitPool.AddVoluntaryExit(ctx, &signedExit); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to add proposer slashing to pool: %v", err))
			}
			if err := backend.Publisher.PublishVoluntaryExit(ctx, &signedExit); err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to publish proposer slashing: %v", err))
			}
			return eth2api.RespondOK(eth2api.Wrap(backend.VoluntaryExitPool.All()))
		})
}
