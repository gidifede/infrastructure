package api

import (
	"facility_queryhandler/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GetNetwork(c *gin.Context) {
	res, err := service.GetNetwork(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilityShortStats(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilityShortStats(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilityParcelDetails(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilityParcelDetails(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilityParcelStats(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilityParcelStats(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilitySortingMachineStats(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilitySortingMachineStats(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilityVehicleStats(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilityVehicleStats(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetFacilityVehicleDetails(c *gin.Context) {
	facilityID := c.Param("facility_id")
	if facilityID == "" {
		errorMsg := "bad facility id"
		log.Error().Msgf(errorMsg)
		c.JSON(http.StatusBadGateway, fmt.Errorf(errorMsg))
		return
	}
	res, err := service.GetFacilityVehicleDetails(c, facilityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}
