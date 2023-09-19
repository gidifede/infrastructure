package manager

import (
	"context"
	"fmt"
	"parcel-sla-projection/internal"
	"parcel-sla-projection/internal/models"
	"parcel-sla-projection/internal/repository"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type AcceptedManager struct {
	logisticEvent models.Accepted
}

func newAcceptedManager(logisticEvent models.Accepted) IEventManager {
	return &AcceptedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *AcceptedManager) ManageEvent(ctx context.Context) error {
	/*
		Alla ricezione dell'evento viene letto il valore dello SLA nella tabella product e scritto il record nella tabella parcel_sla, il calcolo dell'estimated delivery date è la data di accettazione aggiungendo i giorni di SLA
		- parcle_id
		- exptected_delivery_date
		- delivery_rate all'inserimento questo campo è vuoto
	*/

	//Recupero lo SLA del prodotto
	filters := make(map[string]interface{})
	filters["name"] = a.logisticEvent.Parcel.Name
	products, err := internal.Repo.RetrieveProductDocument(ctx, repository.ProductCollectionEnum, filters)
	if err != nil {
		log.Err(err).Msgf("cannot get product doc with filters %s", filters)
		return err
	}
	var productSLA int
	if len(products) > 0 {
		productSLA, _ = strconv.Atoi(products[0].SLA)
	} else {
		err := fmt.Errorf("product %s not found", a.logisticEvent.Parcel.Name)
		log.Err(err).Msgf("")
		return err
	}

	//Calcolo l'exptected_delivery_date
	dayToAdd := time.Duration(productSLA) * 24 * time.Hour
	exptectedDeliveryDate := a.logisticEvent.Timestamp.Add(dayToAdd)

	parcelSLA := &repository.ParcelSLA{
		ParcelID:             a.logisticEvent.Parcel.ID,
		ExpectedDeliveryDate: exptectedDeliveryDate,
	}

	//Inserisco il document
	err = internal.Repo.InsertNewDocument(ctx, repository.ParcelSLACollectionEnum, parcelSLA)
	if err != nil {
		log.Err(err).Msgf("cannot insert new parcel_sla doc %s", parcelSLA)
		return err
	}
	return nil
}
