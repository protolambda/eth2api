package beaconapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
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

// Handles publishing of attestations, stores them in the pool and send them to the publisher.
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
				if err := backend.AttestationPool.AddAttestation(att, comm); err != nil {
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
				return eth2api.RespondBadInputData(&eth2api.IndexedErrorMessage{
					Code:     400,
					Message:  "cannot store/publish",
					Failures: ierrs,
				})
			}
			return eth2api.RespondOKMsg("stored and broadcast all attestations")
		})
}

// Serves attester slashings known by the node but not necessarily incorporated into any block.
func PoolAttesterSlashings(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/pool/attester_slashings",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			return eth2api.RespondOK(eth2api.Wrap(backend.AttesterSlashingPool.All()))
		})
}
