package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"washo.com/main/middleware"
	"washo.com/main/repositories"
	"washo.com/main/services"
)

type VehicleHandler struct {
	service *services.VehicleService
}

func NewVehicleHandler(service *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "create.html", nil)
		return
	}

	var vehicle repositories.CreateVehicleRequest
	if err := c.ShouldBind(&vehicle); err != nil {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	claims := middleware.JwtClaims(c)
	if claims["id"] == nil {
		c.HTML(http.StatusUnauthorized, "create.html", gin.H{
			"Error": "Unauthorized",
		})
		return
	}
	idFloat, ok := claims["id"].(float64)
	if !ok {
		c.HTML(http.StatusUnauthorized, "create.html", gin.H{
			"Error": "Unauthorized",
		})
		return
	}
	vehicle.UID = fmt.Sprintf("%d", uint(idFloat))
	vehicleID, err := h.service.CreateVehicle(&vehicle)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "create.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/vehicles/%s", vehicleID))
}

func (h *VehicleHandler) GetVehiclesByProcess(c *gin.Context) {
	process := []string{"Waiting", "Washing", "Finish"}

	// Call the service to get the vehicles grouped by status
	groupedVehicles, err := h.service.GetVehiclesByProcess(process)
	if err != nil {
		// Handle error by showing it on the page
		c.HTML(http.StatusInternalServerError, "list.html", gin.H{"error": err.Error()})
		return
	}

	// Pass the grouped vehicles to the template
	c.HTML(http.StatusOK, "list.html", gin.H{
		"groupedVehicles": groupedVehicles,
	})
}

func (h *VehicleHandler) GetVehicleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Redirect(http.StatusSeeOther, "/vehicles")
		return
	}
	vehicle, err := h.service.GetVehicleByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/vehicles")
		return
	}
	var currentUserID uint
	var username string
	claims := middleware.JwtClaims(c)
	if claims != nil {
		if idFloat, ok := claims["id"].(float64); ok {
			currentUserID = uint(idFloat)
		}
		if usernameClaim, ok := claims["username"].(string); ok {
			username = usernameClaim
		}
	}
	c.HTML(http.StatusOK, "viewvehicle.html", gin.H{
		"Name":          vehicle.Name,
		"Username":      vehicle.User.Username,
		"Process":       vehicle.Process,
		"Contact":       vehicle.Contact,
		"Plate":         vehicle.Plate,
		"Date":          vehicle.Date,
		"Time":     	 vehicle.Time,
		"ID":            vehicle.ID,
		"IsOwner":       currentUserID == vehicle.UserID,
		"IsAdmin":       username,
		"CurrentUser":   username,
	})
}

func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	id := c.Param("id")

	// Show edit form for GET requests
	if c.Request.Method == http.MethodGet {
		vehicle, err := h.service.GetVehicleByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Check if user owns this vehicle
		claims := middleware.JwtClaims(c)
		if claims == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		var username string
		if usernameClaim, ok := claims["username"].(string); ok {
			username = usernameClaim
			if !strings.Contains(username, "@admin") {
				c.HTML(http.StatusForbidden, "edit.html", gin.H{
					"Error": "Not authorized to edit this vehicle",
				})
				return
			}
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"ID":      vehicle.ID,
			"Name":    vehicle.Name,
			"Contact": vehicle.Contact,
			"Process": vehicle.Process,
			"Plate":   vehicle.Plate,
		})
		return
	}

	// Handle POST request to update vehicle
	var updatedVehicle repositories.CreateVehicleRequest
	if err := c.ShouldBind(&updatedVehicle); err != nil {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"Error":   err.Error(),
			"Name":    updatedVehicle.Name,
			"Contact": updatedVehicle.Contact,
			"Process": updatedVehicle.Process,
			"Plate":   updatedVehicle.Plate,
		})
		return
	}

	if err := h.service.UpdateVehicle(id, updatedVehicle); err != nil {
		c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
			"Error":   err.Error(),
			"ID":      id,
			"Name":    updatedVehicle.Name,
			"Contact": updatedVehicle.Contact,
			"Process": updatedVehicle.Process,
			"Plate":   updatedVehicle.Plate,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/vehicles/%s", id))
}

func (h *VehicleHandler) DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	// Show delete confirmation for GET requests
	if c.Request.Method == http.MethodGet {
		vehicle, err := h.service.GetVehicleByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Check if user owns this vehicle
		claims := middleware.JwtClaims(c)
		if claims == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		if idFloat, ok := claims["id"].(float64); ok {
			currentUserID := uint(idFloat)
			if currentUserID != vehicle.UserID {
				c.HTML(http.StatusForbidden, "mylist.html", gin.H{
					"Error": "Not authorized to delete this vehicle",
				})
				return
			}
		}

		c.Redirect(http.StatusOK, "/vehicles")
		return
	}

	// Handle DELETE request
	if err := h.service.DeleteVehicle(id); err != nil {
		c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicles")
}
