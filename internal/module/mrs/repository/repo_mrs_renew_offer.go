package repository

import (
	"codebase-app/internal/module/mrs/entity"
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *mrsRepository) RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error {
	var (
		newWACId string = ulid.Make().String()
		waccs           = make([]daoWACC, 0)
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to rollback transaction")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to commit transaction")
			}
		}
	}()

	// get WAC conditions
	err = r.getVehicleConditions(ctx, req, tx, &waccs)
	if err != nil {
		return err
	}

	isWACStillNeedFollowUp, WaccNotInterestedLeft, err := r.validateVehicleConditions(req, &waccs)
	if err != nil {
		return err
	}

	totalNewLeads := len(req.VehicleConditionIds)
	if totalNewLeads > 0 {
		err = r.createWACCopy(ctx, req, tx, req.WacId, newWACId, totalNewLeads)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create WAC copy")
			return err
		}

		err = r.moveVehicleConditionsToNewWAC(ctx, req, tx, newWACId, totalNewLeads)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to move WAC conditions")
			return err
		}

		err = r.updateTotalFollowUps(ctx, req.WacId, tx, totalNewLeads, WaccNotInterestedLeft)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update total follow ups")
			return err
		}

		if isWACStillNeedFollowUp {
			err = r.ExtendFollowUpAndRecountingTotalFollowUps(ctx, tx, req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to extend follow up and recounting total follow ups")
				return err
			}

			err = r.createFollowUpLog(ctx, tx, req.WacId, "perlu follow up lagi karena masih ada kondisi yang tidak tertarik")
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create follow up log")
				return err
			}
		} else {
			err := r.removeWACFromFollowUpList(ctx, tx, req.WacId)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to remove WAC from follow up list")
				return err
			}

			err = r.createFollowUpLog(ctx, tx, req.WacId, "semua kondisi berhasil menjadi leads")
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
				return err
			}
		}
	} else { // if empty condition ids
		err = r.updateFollowUpDeadline(ctx, tx, req.WacId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to update follow up deadline")
			return err
		}

		err = r.createFollowUpLog(ctx, tx, req.WacId, "follow up diperpanjang 7 hari ke depan")
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repository::RenewWAC - Failed to create log")
			return err
		}
		return nil
	}

	return nil
}
