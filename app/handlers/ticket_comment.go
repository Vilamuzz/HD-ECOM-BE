package handlers

import (
 "app/domain/models"
 "app/domain/requests"
 "app/helpers"
 "net/http"
 "strconv"

 "github.com/gin-gonic/gin"
)

func (r *appRoute) TicketCommentRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-comments")
	api.POST("", r.createTicketComment)
	api.GET("", r.getTicketComments)
	api.GET("/:id", r.getTicketCommentByID)
	api.PUT("/:id", r.updateTicketComment)
	api.DELETE("/:id", r.deleteTicketComment)
}

// CreateTicketComment godoc
// @Summary Create a new ticket comment
// @Description Create a new comment on a ticket
// @Tags ticket-comments
// @Accept json
// @Produce json
// @Param comment body requests.CreateTicketCommentRequest true "Comment Data"
// @Success 201 {object} helpers.Response{data=requests.TicketCommentResponse}
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-comments [post]
func (r *appRoute) createTicketComment(c *gin.Context) {
 var req requests.CreateTicketCommentRequest
 if err := c.ShouldBindJSON(&req); err != nil {
 c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil))
 return
 }

 comment := models.TicketComment{
 IDTicket: req.IDTicket,
 IDUser: req.IDUser,
 IsiPesan: req.IsiPesan,
 TanggalDibuat: req.TanggalDibuat,
 }

 if err := r.Service.CreateTicketComment(&comment); err != nil {
 c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create comment", nil, nil))
 return
 }

 resp := requests.TicketCommentResponse{
 IDComment: comment.IDComment,
 IDTicket: comment.IDTicket,
 IDUser: comment.IDUser,
 IsiPesan: comment.IsiPesan,
 TanggalDibuat: comment.TanggalDibuat,
 }

 c.JSON(http.StatusCreated, helpers.NewResponse(http.StatusCreated, "Comment created successfully", nil, resp))
}

// GetTicketComments godoc
// @Summary Get all ticket comments
// @Description Get all comments, optionally filtered by ticket ID
// @Tags ticket-comments
// @Produce json
// @Param ticket_id query int false "Filter by Ticket ID"
// @Success 200 {object} helpers.Response{data=[]models.TicketComment}
// @Failure 500 {object} helpers.Response
// @Router /ticket-comments [get]
func (r *appRoute) getTicketComments(c *gin.Context) {
	ticketID, _ := strconv.Atoi(c.Query("ticket_id"))

	var comments []models.TicketComment
	var err error

	if ticketID > 0 {
		comments, err = r.Service.GetTicketCommentsByTicketID(ticketID)
	} else {
		comments, err = r.Service.GetTicketComments()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to get comments", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Comments retrieved successfully", nil, comments))
}

// GetTicketCommentByID godoc
// @Summary Get a ticket comment by ID
// @Description Get a specific comment by its ID
// @Tags ticket-comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} helpers.Response{data=models.TicketComment}
// @Failure 404 {object} helpers.Response
// @Router /ticket-comments/{id} [get]
func (r *appRoute) getTicketCommentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid comment ID", nil, nil))
		return
	}

	comment, err := r.Service.GetTicketCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewResponse(http.StatusNotFound, "Comment not found", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Comment retrieved successfully", nil, comment))
}

// UpdateTicketComment godoc
// @Summary Update a ticket comment
// @Description Update an existing comment
// @Tags ticket-comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param comment body requests.UpdateTicketCommentRequest true "Updated Comment Data"
// @Success 200 {object} helpers.Response{data=requests.TicketCommentResponse}
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-comments/{id} [put]
func (r *appRoute) updateTicketComment(c *gin.Context) {
 id, err := strconv.Atoi(c.Param("id"))
 if err != nil {
 c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid comment ID", nil, nil))
 return
 }

 var req requests.UpdateTicketCommentRequest
 if err := c.ShouldBindJSON(&req); err != nil {
 c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil))
 return
 }

 comment := models.TicketComment{
 IDComment: id,
 IDTicket: req.IDTicket,
 IDUser: req.IDUser,
 IsiPesan: req.IsiPesan,
 TanggalDibuat: req.TanggalDibuat,
 }

 if err := r.Service.UpdateTicketComment(&comment); err != nil {
 c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to update comment", nil, nil))
 return
 }

 resp := requests.TicketCommentResponse{
 IDComment: comment.IDComment,
 IDTicket: comment.IDTicket,
 IDUser: comment.IDUser,
 IsiPesan: comment.IsiPesan,
 TanggalDibuat: comment.TanggalDibuat,
 }

 c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Comment updated successfully", nil, resp))
}

// DeleteTicketComment godoc
// @Summary Delete a ticket comment
// @Description Delete an existing comment
// @Tags ticket-comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-comments/{id} [delete]
func (r *appRoute) deleteTicketComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid comment ID", nil, nil))
		return
	}

	if err := r.Service.DeleteTicketComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to delete comment", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Comment deleted successfully", nil, nil))
}