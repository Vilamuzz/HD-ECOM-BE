package handlers

import (
	"app/domain/models"
	"app/domain/requests"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/tickets")

	// Public endpoints (no auth required)
	api.GET("/:id", r.getTicketByID)

	// Authenticated user endpoints
	api.POST("", r.Middleware.Auth(), r.createTicket)
	api.GET("/my-tickets", r.Middleware.Auth(), r.getMyTickets)
	api.PUT("/:id", r.Middleware.Auth(), r.updateTicket)
	api.DELETE("/:id", r.Middleware.Auth(), r.deleteTicket)

	// Admin-only endpoints
	api.GET("", r.Middleware.Auth(), r.Middleware.RequireRole(models.RoleAdmin), r.getTickets)
}

// CreateTicket godoc
// @Summary Create a new ticket
// @Description Create a new ticket (tipe_pengaduan will be auto-filled based on user role)
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ticket body requests.TicketCreateRequest true "Ticket Data"
// @Success 201 {object} helpers.Response{data=requests.TicketResponse}
// @Router /tickets [post]
func (r *appRoute) createTicket(c *gin.Context) {
	var req requests.TicketCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get user from context (required - no fallback)
	userData, exists := c.Get("userData")
	if !exists {
		response := helpers.NewResponse(http.StatusUnauthorized, "User authentication required", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	user, ok := userData.(models.User)
	if !ok {
		response := helpers.NewResponse(http.StatusUnauthorized, "Invalid user data", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Use authenticated user's data
	tipePengaduan := user.Role
	if tipePengaduan == "" {
		tipePengaduan = models.RoleCustomer // default fallback
	}

	ticket := models.Ticket{
		// KodeTiket will be auto-generated in service
		UserID:        user.ID, // Always use authenticated user's ID
		Judul:         req.Judul,
		Deskripsi:     req.Deskripsi,
		CategoryID:    req.CategoryID,
		PriorityID:    3, // Default priority ID
		StatusID:      1, // Default status ID
		TipePengaduan: tipePengaduan,
	}

	if err := r.Service.CreateTicket(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	resp := requests.TicketResponse{
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
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket created successfully", nil, resp)
	c.JSON(http.StatusCreated, response)
}

// GetTickets godoc
// @Summary Get all tickets
// @Description Get a list of all tickets (Admin only), cursor-based pagination
// @Tags tickets
// @Produce json
// @Security BearerAuth
// @Param role query string false "Filter by tipe_pengaduan (customer, seller, admin, support)"
// @Param status query int false "Filter by status ID"
// @Param priority query int false "Filter by priority ID"
// @Param category query int false "Filter by category ID"
// @Param limit query int false "Items per page (default: 10)"
// @Param cursor query string false "Cursor for next page"
// @Success 200 {object} helpers.Response{data=requests.TicketListResponse}
// @Router /tickets [get]
func (r *appRoute) getTickets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	cursor := c.Query("cursor")
	filterType := c.Query("role")

	// Get filter params for status, priority, and category
	statusID, _ := strconv.Atoi(c.Query("status"))
	priorityID, _ := strconv.Atoi(c.Query("priority"))
	categoryID, _ := strconv.Atoi(c.Query("category"))

	// Call service with all filters - filtering happens at DB level
	tickets, nextCursor, err := r.Service.GetTicketsCursor(limit, cursor, filterType, statusID, priorityID, categoryID)
	if err != nil {
		response := helpers.NewResponse(500, "Failed to get tickets", nil, nil)
		c.JSON(500, response)
		return
	}

	// No need to filter here anymore - already filtered by database
	var resp []requests.TicketResponse
	for _, ticket := range tickets {
		username := ""
		if ticket.User.Username != "" {
			username = ticket.User.Username
		}
		resp = append(resp, requests.TicketResponse{
			ID:                ticket.ID,
			KodeTiket:         ticket.KodeTiket,
			UserID:            ticket.UserID,
			Username:          username,
			Judul:             ticket.Judul,
			Deskripsi:         ticket.Deskripsi,
			CategoryID:        ticket.CategoryID,
			PriorityID:        ticket.PriorityID,
			StatusID:          ticket.StatusID,
			TipePengaduan:     ticket.TipePengaduan,
			TanggalDibuat:     ticket.TanggalDibuat.Format("2006-01-02T15:04:05Z07:00"),
			TanggalDiperbarui: ticket.TanggalDiperbarui.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	responseData := map[string]interface{}{
		"data": resp,
		"meta": map[string]interface{}{
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	}
	response := helpers.NewResponse(200, "Tickets retrieved successfully", nil, responseData)
	c.JSON(200, response)
}

// GetTicketByID godoc
// @Summary Get a ticket by ID
// @Description Get a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} helpers.Response{data=requests.TicketResponse}
// @Router /tickets/{id} [get]
func (r *appRoute) getTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ticket, err := r.Service.GetTicketByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Ticket not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	username := ""
	if ticket.User.Username != "" {
		username = ticket.User.Username
	}

	resp := requests.TicketResponse{
		ID:                ticket.ID,
		KodeTiket:         ticket.KodeTiket,
		UserID:            ticket.UserID,
		Username:          username,
		Judul:             ticket.Judul,
		Deskripsi:         ticket.Deskripsi,
		CategoryID:        ticket.CategoryID,
		PriorityID:        ticket.PriorityID,
		StatusID:          ticket.StatusID,
		TipePengaduan:     ticket.TipePengaduan,
		TanggalDibuat:     ticket.TanggalDibuat.Format("2006-01-02T15:04:05Z07:00"),
		TanggalDiperbarui: ticket.TanggalDiperbarui.Format("2006-01-02T15:04:05Z07:00"),
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket retrieved successfully", nil, resp)
	c.JSON(http.StatusOK, response)
}

// UpdateTicket godoc
// @Summary Update a ticket
// @Description Update a ticket by its ID (tipe_pengaduan will be auto-filled based on user role)
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body requests.TicketCreateRequest true "Updated Ticket Data"
// @Success 200 {object} helpers.Response{data=requests.TicketResponse}
// @Router /tickets/{id} [put]
func (r *appRoute) updateTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var req requests.TicketCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Ambil user dari context
	userData, exists := c.Get("userData")
	var tipePengaduan models.UserRole
	var userID uint64
	if exists {
		user, ok := userData.(models.User)
		if ok && user.Role != "" {
			tipePengaduan = user.Role
			userID = user.ID
		} else {
			// fallback ke default customer jika role kosong
			tipePengaduan = models.RoleCustomer
			userID = req.UserID
		}
	} else {
		// fallback jika tidak ada user di context - default customer
		tipePengaduan = models.RoleCustomer
		userID = req.UserID
	}

	ticket := models.Ticket{
		ID:            id,
		KodeTiket:     req.KodeTiket,
		UserID:        userID,
		Judul:         req.Judul,
		Deskripsi:     req.Deskripsi,
		CategoryID:    req.CategoryID,
		PriorityID:    req.PriorityID,
		StatusID:      req.StatusID,
		TipePengaduan: tipePengaduan,
	}

	if err := r.Service.UpdateTicket(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	resp := requests.TicketResponse{
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
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket updated successfully", nil, resp)
	c.JSON(http.StatusOK, response)
}

// DeleteTicket godoc
// @Summary Delete a ticket
// @Description Delete a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} helpers.Response
// @Router /tickets/{id} [delete]
func (r *appRoute) deleteTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicket(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

// GetMyTickets godoc
// @Summary Get current user's tickets
// @Description Get all tickets belonging to the authenticated user
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Success 200 {object} helpers.Response{data=[]requests.TicketResponse}
// @Failure 401 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /tickets/my-tickets [get]
func (r *appRoute) getMyTickets(c *gin.Context) {
	// Get authenticated user from context
	userData, exists := c.Get("userData")
	if !exists {
		response := helpers.NewResponse(http.StatusUnauthorized, "User authentication required", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	user, ok := userData.(models.User)
	if !ok {
		response := helpers.NewResponse(http.StatusUnauthorized, "Invalid user data", nil, nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Get tickets by user ID
	tickets, err := r.Service.GetTicketsByUserID(int(user.ID))
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get user tickets", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var resp []requests.TicketResponse
	for _, ticket := range tickets {
		resp = append(resp, requests.TicketResponse{
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
		})
	}

	response := helpers.NewResponse(http.StatusOK, "User tickets retrieved successfully", nil, resp)
	c.JSON(http.StatusOK, response)
}
