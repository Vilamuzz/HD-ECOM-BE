package handlers

import (
	"app/domain/models"
	"app/domain/requests"
	"app/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketAssignmentRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-assignments")
	api.Use(r.Middleware.Auth())
	api.POST("", r.createTicketAssignment)
	api.GET("", r.getTicketAssignments)
	api.GET("/my-assignments", r.getMySupportAssignments)
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
// @Security BearerAuth
// @Param assignment body requests.CreateTicketAssignmentRequest true "Ticket Assignment Data"
// @Success 201 {object} helpers.Response{data=requests.TicketAssignmentResponse}
// @Router /ticket-assignments [post]
func (r *appRoute) createTicketAssignment(c *gin.Context) {
	var req requests.CreateTicketAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tglDitugaskan, err := time.Parse("2006-01-02T15:04:05Z07:00", req.TanggalDitugaskan)
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid date format for tanggal_ditugaskan", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	assignment := models.TicketAssignment{
		TicketID:          req.TicketID,
		AdminID:           req.AdminID,
		TanggalDitugaskan: tglDitugaskan,
	}

	if err := r.Service.CreateTicketAssignment(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket assignment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	resp := mapTicketAssignmentToResponse(&assignment)
	response := helpers.NewResponse(http.StatusCreated, "Ticket assignment created successfully", nil, resp)
	c.JSON(http.StatusCreated, response)
}

// GetTicketAssignments godoc
// @Summary Get all ticket assignments
// @Description Get a list of all ticket assignments
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=[]requests.TicketAssignmentResponse}
// @Router /ticket-assignments [get]
func (r *appRoute) getTicketAssignments(c *gin.Context) {
	assignments, err := r.Service.GetTicketAssignments()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket assignments", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respList []requests.TicketAssignmentResponse
	for _, a := range assignments {
		respList = append(respList, mapTicketAssignmentToResponse(&a))
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket assignments retrieved successfully", nil, respList)
	c.JSON(http.StatusOK, response)
}

// GetMySupportAssignments godoc
// @Summary Get my ticket assignments as support user
// @Description Get ticket assignments where the admin_id matches the authenticated support user's ID
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=[]requests.TicketAssignmentResponse}
// @Router /ticket-assignments/my-assignments [get]
func (r *appRoute) getMySupportAssignments(c *gin.Context) {
	userInterface, exists := c.Get("userData")
	if !exists || userInterface == nil {
		response := helpers.NewResponse(http.StatusUnauthorized, "User not authenticated", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	user, ok := userInterface.(models.User)
	if !ok {
		response := helpers.NewResponse(http.StatusUnauthorized, "Invalid user data", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	if user.Role != "support" {
		response := helpers.NewResponse(http.StatusForbidden, "Access denied. Only support users can access this endpoint", nil, nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	assignments, err := r.Service.GetTicketAssignments()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket assignments", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respList []map[string]interface{}
	for _, a := range assignments {
		if a.AdminID == int(user.ID) {
			ticket, err := r.Service.GetTicketByID(a.TicketID)
			if err != nil {
				continue
			}

			// Fetch category - try Nama field (Indonesian naming)
			categoryName := ""
			if ticket.CategoryID > 0 {
				category, err := r.Service.GetTicketCategoryByID(ticket.CategoryID)
				if err == nil && category != nil {
					categoryName = category.NamaCategory // Change from Name to Nama
				}
			}

			// Fetch priority - try Nama field
			priorityName := ""
			if ticket.PriorityID > 0 {
				priority, err := r.Service.GetTicketPriorityByID(ticket.PriorityID)
				if err == nil && priority != nil {
					priorityName = priority.NamaPriority // Change from Name to Nama
				}
			}

			// Fetch status - try Nama field
			statusName := ""
			if ticket.StatusID > 0 {
				status, err := r.Service.GetTicketStatusByID(ticket.StatusID)
				if err == nil && status != nil {
					statusName = status.NamaStatus // Change from Name to Nama
				}
			}

			username := ""
			if ticket.User.Username != "" {
				username = ticket.User.Username
			}

			respList = append(respList, map[string]interface{}{
				"id_assignment":      a.ID,
				"id_ticket":          a.TicketID,
				"id_admin":           a.AdminID,
				"tanggal_ditugaskan": a.TanggalDitugaskan.Format("2006-01-02T15:04:05Z"),
				"ticket": map[string]interface{}{
					"id_ticket":      ticket.ID,
					"kode_ticket":    ticket.KodeTiket,
					"username":       username,
					"judul":          ticket.Judul,
					"deskripsi":      ticket.Deskripsi,
					"category_name":  categoryName,
					"priority_name":  priorityName,
					"status_name":    statusName,
					"tipe_pengaduan": ticket.TipePengaduan,
				},
			})
		}
	}

	response := helpers.NewResponse(http.StatusOK, "My ticket assignments retrieved successfully", nil, respList)
	c.JSON(http.StatusOK, response)
}

// GetTicketAssignmentByID godoc
// @Summary Get a ticket assignment by ID
// @Description Get a ticket assignment by its ID
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
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
	return requests.TicketAssignmentResponse{
		AssignmentID:      a.ID,
		TicketID:          a.TicketID,
		AdminID:           a.AdminID,
		TanggalDitugaskan: a.TanggalDitugaskan.Format("2006-01-02T15:04:05Z"),
	}
}

// mapTicketAssignmentToResponseWithTicket maps TicketAssignment with Ticket data to response DTO
func mapTicketAssignmentToResponseWithTicket(a *models.TicketAssignment, ticket *models.Ticket) requests.TicketAssignmentResponse {
	return requests.TicketAssignmentResponse{
		AssignmentID:      a.ID,
		TicketID:          a.TicketID,
		AdminID:           a.AdminID,
		TanggalDitugaskan: a.TanggalDitugaskan.Format("2006-01-02T15:04:05Z"),
		Ticket: &requests.TicketResponse{
			ID:                ticket.ID,
			KodeTiket:         ticket.KodeTiket,
			UserID:            ticket.UserID,
			Judul:             ticket.Judul,
			Deskripsi:         ticket.Deskripsi,
			CategoryID:        ticket.CategoryID,
			PriorityID:        ticket.PriorityID,
			StatusID:          ticket.StatusID,
			TipePengaduan:     ticket.TipePengaduan,
			TanggalDibuat:     ticket.TanggalDibuat.Format("2006-01-02T15:04:05Z07:00"),
			TanggalDiperbarui: ticket.TanggalDiperbarui.Format("2006-01-02T15:04:05Z07:00"),
		},
	}
}

// UpdateTicketAssignment godoc
// @Summary Update a ticket assignment
// @Description Update a ticket assignment by its ID
// @Tags ticket-assignments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assignment ID"
// @Param assignment body requests.CreateTicketAssignmentRequest true "Updated Assignment Data"
// @Success 200 {object} helpers.Response{data=requests.TicketAssignmentResponse}
// @Router /ticket-assignments/{id} [put]
func (r *appRoute) updateTicketAssignment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid assignment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req requests.CreateTicketAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tglDitugaskan, err := time.Parse("2006-01-02T15:04:05Z07:00", req.TanggalDitugaskan)
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid date format for tanggal_ditugaskan", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	assignment := models.TicketAssignment{
		ID:                id,
		TicketID:          req.TicketID,
		AdminID:           req.AdminID,
		TanggalDitugaskan: tglDitugaskan,
	}

	if err := r.Service.UpdateTicketAssignment(&assignment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket assignment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	resp := mapTicketAssignmentToResponse(&assignment)
	response := helpers.NewResponse(http.StatusOK, "Ticket assignment updated successfully", nil, resp)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketAssignment godoc
// @Summary Delete a ticket assignment
// @Description Delete a ticket assignment by its ID
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
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
