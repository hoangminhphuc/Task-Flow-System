package model

import "strings"

// ProfileResponse represents the data returned in the user profile response.
type ProfileResponse struct {
    User  *User `json:"user"`
    Stats UserProfileStats   `json:"stats"`
}

// StatsUser holds statistics information for the user profile.
type UserProfileStats struct {
	UserID             int     `json:"user_id"`
	FullName           string  `json:"full_name"`
	DateJoined         string  `json:"date_joined"`
	TotalTasks         int     `json:"total_tasks"`
	CompletedTasks     int     `json:"completed_tasks"`
	PendingTasks       int     `json:"pending_tasks"`
	IncompletedTasks   int 		 `json:"incompleted_tasks"`
	OldestPendingTitle *string `json:"oldest_pending_title"`
	OldestPendingDate  *string `json:"oldest_pending_date"`
	RecentTasks        *string  `json:"recent_tasks"`
}

func (u *UserProfileStats) GetRecentTasks() []string {
	tasksList := strings.Split(*u.RecentTasks, ",")
	return tasksList
}
