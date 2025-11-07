package handlers

import (
	"app/domain/models"
	"app/domain/requests"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketAssignmentRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-assignments")
	api.POST("", r.createTicketAssignment)
	api.GET("", r.getTicketAssignments)
	api.GET("/:id", r.getTicketAssignmentByID)
	api.PUT("/:id", r.updateTicketAssignment)
	api.DELETE("/:id", r.deleteTicketAssignment)
}

// CreateTicketAssignment godoc
// @Summary Create a new ticket assignment
// @Description Create a new ticket assignment
// @Tags ticket-assignments
// @Accept json
// @Produce json
// @Param assignment body models.TicketAssignment true "Ticket Assignment Data"
// @Success 201 {object} helpers.Response{data=requests.TicketAssignmentResponse}
// @Router /ticket-assignments [post]
func (r *appRoute) createTicketAssignment(c *gin.Context) {
	var assignment models.TicketAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.CreateTicketAssignment(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket assignment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket assignment created successfully", nil, assignment)
	c.JSON(http.StatusCreated, response)
}

// GetTicketAssignments godoc
// @Summary Get all ticket assignments
// @Description Get a list of all ticket assignments
// @Tags ticket-assignments
// @Produce json
// @Success 200 {object} helpers.Response{data=[]requests.TicketAssignmentResponse}
// @Router /ticket-assignments [get]
func (r *appRoute) getTicketAssignments(c *gin.Context) {
	assignments, err := r.Service.GetTicketAssignments()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket assignments", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Map to DTO
	var respList []requests.TicketAssignmentResponse
	for _, a := range assignments {
		respList = append(respList, mapTicketAssignmentToResponse(&a))
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket assignments retrieved successfully", nil, respList)
	c.JSON(http.StatusOK, response)
}

// GetTicketAssignmentByID godoc
// @Summary Get a ticket assignment by ID
// @Description Get a ticket assignment by its ID
// @Tags ticket-assignments
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} helpers.Response{data=requests.TicketAssignmentResponse}
// @Router /ticket-assignments/{id} [get]
func (r *appRoute) getTicketAssignmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid assignment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	assignment, err := r.Service.GetTicketAssignmentByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Assignment not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	resp := mapTicketAssignmentToResponse(assignment)
	response := helpers.NewResponse(http.StatusOK, "Ticket assignment retrieved successfully", nil, resp)
	c.JSON(http.StatusOK, response)
}

// mapTicketAssignmentToResponse maps TicketAssignment model to TicketAssignmentResponse DTO
func mapTicketAssignmentToResponse(a *models.TicketAssignment) requests.TicketAssignmentResponse {
	resp := requests.TicketAssignmentResponse{
		IDAssignment:      a.IDAssignment,
		IDTicket:          a.IDTicket,
		IDAdmin:           a.IDAdmin,
		TanggalDitugaskan: a.TanggalDitugaskan.Format("2006-01-02T15:04:05Z"),
	}
	if a.Ticket != nil {
		resp.Ticket = &requests.TicketResponse{
			IDTicket:          a.Ticket.IDTicket,
			KodeTiket:         a.Ticket.KodeTiket,
			IDUser:            a.Ticket.IDUser,
			Judul:             a.Ticket.Judul,
			Deskripsi:         a.Ticket.Deskripsi,
			IDCategory:        a.Ticket.IDCategory,
			IDPriority:        a.Ticket.IDPriority,
			IDStatus:          a.Ticket.IDStatus,
			TipePengaduan:     a.Ticket.TipePengaduan,
			TanggalDibuat:     a.Ticket.TanggalDibuat.Format("2006-01-02T15:04:05Z"),
			TanggalDiperbarui: a.Ticket.TanggalDiperbarui.Format("2006-01-02T15:04:05Z"),
		}
	}
	if a.Admin != nil {
		resp.Admin = &requests.UserSimpleResponse{
			IDUser:    a.Admin.IDUser,
			Nama:      a.Admin.Nama,
			Email:     a.Admin.Email,
			Username:  a.Admin.Username,
			Role:      a.Admin.Role,
			CreatedAt: a.Admin.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return resp
}

// UpdateTicketAssignment godoc
// @Summary Update a ticket assignment
// @Description Update a ticket assignment by its ID
// @Tags ticket-assignments
// @Accept json
// @Produce json
// @Param id path int true "Assignment ID"
// @Param assignment body models.TicketAssignment true "Updated Assignment Data"
// @Success 200 {object} helpers.Response{data=requests.TicketAssignmentResponse}
// @Router /ticket-assignments/{id} [put]
func (r *appRoute) updateTicketAssignment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid assignment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var assignment models.TicketAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	assignment.IDAssignment = id
	if err := r.Service.UpdateTicketAssignment(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket assignment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket assignment updated successfully", nil, assignment)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketAssignment godoc
// @Summary Delete a ticket assignment
// @Description Delete a ticket assignment by its ID
// @Tags ticket-assignments
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} helpers.Response
// @Router /ticket-assignments/{id} [delete]
func (r *appRoute) deleteTicketAssignment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid assignment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicketAssignment(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket assignment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket assignment deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}
