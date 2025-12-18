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
	api.GET("/my-counts", r.Middleware.RequireRole(models.RoleSupport), r.getMyAssignedTicketCounts) // New route for support users
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
		PriorityID:        req.PriorityID,
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
// @Description Get ticket assignments where the admin_id matches the authenticated support user's ID, with cursor pagination and optional status filter
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Items per page (default: 10)"
// @Param cursor query string false "Cursor for next page"
// @Param status query string false "Filter by ticket status name (e.g., 'Open')"
// @Success 200 {object} helpers.Response{data=[]map[string]interface{}, meta=map[string]interface{}}
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

	// Parse query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	cursor := c.Query("cursor")
	statusName := c.Query("status") // Changed to string for status name

	// Call service with cursor pagination and status filter
	assignments, nextCursor, err := r.Service.GetTicketAssignmentsByAdminIDCursor(int(user.ID), limit, cursor, statusName) // Pass statusName
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket assignments", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respList []map[string]interface{}
	for _, a := range assignments {
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

	responseData := map[string]interface{}{
		"data": respList,
		"meta": map[string]interface{}{
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	}
	response := helpers.NewResponse(http.StatusOK, "My ticket assignments retrieved successfully", nil, responseData)
	c.JSON(http.StatusOK, response)
}

// GetMyAssignedTicketCounts godoc
// @Summary Get assigned ticket counts for support user
// @Description Get total and in-progress assigned ticket counts for the authenticated support user
// @Tags ticket-assignments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=map[string]int}
// @Router /ticket-assignments/my-counts [get]
func (r *appRoute) getMyAssignedTicketCounts(c *gin.Context) {
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

	// Get counts
	totalAssigned, assignedErr := r.Service.GetAssignedTicketCountByAdminID(int(user.ID))
	inProgressAssigned, inProgressAssignedErr := r.Service.GetAssignedTicketCountByAdminIDAndStatus(int(user.ID), 2) // Status ID 2 = "In Progress"

	responseData := map[string]int{
		"total":        totalAssigned,
		"in_progress":  inProgressAssigned,
	}

	// Add error info if counts failed
	if assignedErr != nil {
		responseData["total_error"] = 1 // Simple error flag
	}
	if inProgressAssignedErr != nil {
		responseData["in_progress_error"] = 1 // Simple error flag
	}

	response := helpers.NewResponse(http.StatusOK, "Assigned ticket counts retrieved successfully", nil, responseData)
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
	response := requests.TicketAssignmentResponse{
		AssignmentID:      a.ID,
		TicketID:          a.TicketID,
		AdminID:           a.AdminID,
		PriorityID:        a.PriorityID,
		TanggalDitugaskan: a.TanggalDitugaskan.Format("2006-01-02T15:04:05Z"),
	}

	// Map Priority if exists
	if a.Priority != nil {
		response.Priority = &requests.PriorityResponse{
			ID:           a.Priority.ID,
			NamaPriority: a.Priority.NamaPriority,
		}
	}

	return response
}

// mapTicketAssignmentToResponseWithTicket maps TicketAssignment with Ticket data to response DTO
func mapTicketAssignmentToResponseWithTicket(a *models.TicketAssignment, ticket *models.Ticket) requests.TicketAssignmentResponse {
	response := requests.TicketAssignmentResponse{
		AssignmentID:      a.ID,
		TicketID:          a.TicketID,
		AdminID:           a.AdminID,
		PriorityID:        a.PriorityID,
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

	// Map Priority if exists
	if a.Priority != nil {
		response.Priority = &requests.PriorityResponse{
			ID:           a.Priority.ID,
			NamaPriority: a.Priority.NamaPriority,
		}
	}

	return response
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
		PriorityID:        req.PriorityID,
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
