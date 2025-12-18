package services

import (
    "app/domain/models"
    "app/helpers"
    "net/http"
)

func (s *appService) GetSupportUsers() helpers.Response {
    users, err := s.repo.GetUsersByRole(models.RoleSupport)
    if err != nil {
        return helpers.NewResponse(http.StatusInternalServerError, "Failed to get support users", nil, nil)
    }

    type SupportUserResponse struct {
        UserID   uint64 `json:"user_id"`
        Username string `json:"username"`
    }

    var resp []SupportUserResponse
    for _, user := range users {
        resp = append(resp, SupportUserResponse{
            UserID:   user.ID,
            Username: user.Username,
        })
    }

    return helpers.NewResponse(http.StatusOK, "Support users retrieved successfully", nil, resp)
}