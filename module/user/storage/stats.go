package storage

import (
	"context"
	"first-proj/module/user/model"
)

func (s *sqlStore) GetStats(ctx context.Context, userId int) (*model.UserProfileStats, error) {
	db := s.db.Table(model.User{}.TableName())

	var userProfile model.UserProfileStats

	query := `
	SELECT 
    u.id AS user_id,
    CONCAT(u.first_name, ' ', u.last_name) AS full_name,
    u.created_at AS date_joined,
    COUNT(t.id) AS total_tasks,
    COUNT(CASE WHEN t.status = 'Done' THEN 1 END) AS completed_tasks,
    COUNT(CASE WHEN t.status = 'Doing' THEN 1 END) AS pending_tasks,
    COUNT(CASE WHEN t.status = 'Deleted' THEN 1 END) AS incompleted_tasks,
    (SELECT t2.title 
      FROM todo_items t2 
      WHERE t2.user_id = u.id AND t2.status = 'Doing'
      ORDER BY t2.created_at ASC LIMIT 1) AS oldest_pending_title, 
    (SELECT t2.created_at 
      FROM todo_items t2 
      WHERE t2.user_id = u.id AND t2.status = 'Doing'
      ORDER BY t2.created_at ASC LIMIT 1) AS oldest_pending_date,
    GROUP_CONCAT(t.title ORDER BY t.created_at DESC SEPARATOR ', ') AS recent_tasks
FROM users u
LEFT JOIN todo_items t ON u.id = t.user_id
WHERE u.id = ?
GROUP BY u.id

	`
	if err := db.Raw(query, userId).Scan(&userProfile).Error; err != nil {
		return nil, err
	}

	return &userProfile, nil
}
