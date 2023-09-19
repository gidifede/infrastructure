package manager

import (
	"context"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/pathfinder"
	"parcel-processing-projection/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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
		- Viene calcolato il percorso che il pacco deve fare (devo recuperare le rotte dalla collection routes per effettuale il calcolo)
		- Viene inserito lo stato
		- Viene inserita la posizione corrente
		- Aggiorno la percentuale di avanzamento
		- Si aggiunge lo stato alla history
		- Riportate le info che troviamo sull'evento
		- Inserisco il record
	*/
	path, err := estimatePath(ctx, a.logisticEvent.FacilityID, a.logisticEvent.Receiver.City)
	if err != nil {
		log.Err(err).Msgf("cannot estimate path for facility %s receiver city %s", a.logisticEvent.FacilityID, a.logisticEvent.Receiver.City)
		return err
	}
	pathCompleted := CalculatePathPercentage(path, a.logisticEvent.FacilityID)
	lastStatus := models.AcceptedEventStatus

	//Recupero il facility type dell'accettazione
	filtersFacilities := bson.M{"facility_id": a.logisticEvent.FacilityID}
	facilities, err := internal.Repo.RetrieveFacilityDocument(ctx, repository.FacilityCollectionEnum, filtersFacilities)
	if err != nil {
		log.Err(err).Msgf("cannot get facility doc with id %s", a.logisticEvent.FacilityID)
		return err
	}
	facilityType := facilities[0].FacilityType

	//PReparo il document per l'insert
	parcelAccepted := &repository.Parcel{
		Name:       a.logisticEvent.Parcel.Name,
		ID:         a.logisticEvent.Parcel.ID,
		Type:       a.logisticEvent.Parcel.Type,
		LastStatus: lastStatus,
		Position: repository.Position{
			PositionType: facilityType,
			PositionID:   a.logisticEvent.FacilityID,
		},
		ParcelPath: repository.ParcelPath{
			Path:          path,
			PathCompleted: pathCompleted,
		},
		History: []repository.Status{{
			Status: lastStatus,
			Date:   a.logisticEvent.Timestamp},
		},
		Sender: struct {
			Name    string "bson:\"name\""
			Address string "bson:\"address\""
			Zipcode string "bson:\"zipcode\""
			City    string "bson:\"city\""
			Nation  string "bson:\"nation\""
		}{
			Name:    a.logisticEvent.Sender.Name,
			Address: a.logisticEvent.Sender.Address,
			Zipcode: a.logisticEvent.Sender.Zipcode,
			City:    a.logisticEvent.Sender.City,
			Nation:  a.logisticEvent.Sender.Nation,
		},
		Receiver: struct {
			Name    string "bson:\"name\""
			Address string "bson:\"address\""
			Zipcode string "bson:\"zipcode\""
			City    string "bson:\"city\""
			Nation  string "bson:\"nation\""
			Number  string "bson:\"number\""
			Email   string "bson:\"email\""
		}{
			Name:    a.logisticEvent.Receiver.Name,
			Address: a.logisticEvent.Receiver.Address,
			Zipcode: a.logisticEvent.Receiver.Zipcode,
			City:    a.logisticEvent.Receiver.City,
			Nation:  a.logisticEvent.Receiver.Nation,
			Number:  a.logisticEvent.Receiver.Number,
			Email:   a.logisticEvent.Receiver.Email,
		},
	}

	err = internal.Repo.InsertNewDocument(ctx, repository.ParcelCollectionEnum, parcelAccepted)
	if err != nil {
		log.Err(err).Msgf("cannot insert parcel doc with id %s", a.logisticEvent.Parcel.ID)
		return err
	}
	return nil
}

// Calcolo del percorso del pacco
func estimatePath(ctx context.Context, sourceLocationID string, destinationCity string) ([]string, error) {
	/*
		Recupero le rotte e le facility dalle rispettive collection
		Metto insieme le info per creare un array di struct di NodeList per effettualre il calcolo dello shortherst path
	*/

	//Rotte
	logisticRoutes, err := internal.Repo.RetrieveRouteDocument(ctx, repository.RouteCollectionEnum, nil)
	if err != nil {
		log.Err(err).Msgf("cannot get all route docs")
		return nil, err
	}

	//Facilities
	logisticFacility, err := internal.Repo.RetrieveFacilityDocument(ctx, repository.FacilityCollectionEnum, nil)
	if err != nil {
		log.Err(err).Msgf("cannot get all facility docs")
		return nil, err
	}

	var nodelist []pathfinder.NodeList
	sourceNode := ""
	destinationNode := ""

	//TODO OTTIMIZZARE LA RICERCA ASSOLUTAMENTE!!! :-x
	for _, route := range logisticRoutes {
		cost := 0
		for _, facility := range logisticFacility {
			if facility.FacilityID == route.SourceFacilityID {
				for _, connection := range facility.Connections {
					if connection.FacilityDestinationID == route.DestFacilityID {
						cost = connection.Distance
					}
				}
			}
			if facility.FacilityID == sourceLocationID {
				sourceNode = facility.FacilityID
			}
			if facility.FacilityLocation.City == destinationCity {
				destinationNode = facility.FacilityID
			}
		}
		if cost != 0 {
			nodelist = append(nodelist, pathfinder.NodeList{Source: route.SourceFacilityID, Destination: route.DestFacilityID, Cost: cost})
		}
	}

	//Effettuo il calcolo
	path := pathfinder.GetPath(nodelist, sourceNode, destinationNode)

	return path, nil

}
