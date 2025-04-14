package models

import (
	"database/sql"
)

type Notification struct {
	Id         int     `json:"id"`
	RelatedId  int     `json:"related_id"`
	Type       string  `json:"type"`
	SenderName string  `json:"sender_name"`
	SenderId   int     `json:"sender_id"`
	ReceverId  int     `json:"recever_id"`
	GroupId    *int    `json:"group_id"`
	GroupName  *string `json:"group_name"`
}

// InsertNotification inserts a new notification into the database
func InsertNotification(db *sql.DB, notif *Notification) (Notification, error) {
	query := `INSERT INTO notifications (type, related_id, recever_id, sender_id, is_read, created_at) VALUES (?, ?, ?, ?, 0, CURRENT_TIMESTAMP)`
	res, err := db.Exec(query, notif.Type, notif.RelatedId, notif.ReceverId, notif.SenderId)
	if err != nil {return Notification{}, err}
	notifId, err := res.LastInsertId()
	if err != nil {return Notification{}, err}
	notifToSend, err := getNotification(db, notifId)
	return notifToSend, err
}

// DeleteNotification deletes a notification by ID
func DeleteNotification(db *sql.DB, notifID, follower, following int) error {
	_, err := db.Exec("DELETE FROM notifications WHERE id = ? or (sender_id = ? and recever_id = ? and type='follow request')", notifID, follower, following)
	return err
}

func DeleteGroupMemberNotifications(db *sql.DB, groupID, userID int) error {
	_, err := db.Exec(`DELETE FROM notifications WHERE related_id IN (SELECT id FROM group_members WHERE group_id = ? AND user_id = ? )`, groupID, userID)
	return err
}

func getNotification(db *sql.DB, notifId int64)(Notification, error) {
	query := `
		SELECT 
		    n.id, 
		    n.type, 
		    n.related_id,
		    n.sender_id,
			n.recever_id,
		    COALESCE(ge.group_id, gm.group_id) AS group_id,
		    u.first_name || ' ' || u.last_name AS sender_name,
		    g.title AS group_name
		FROM notifications AS n
		JOIN users As u ON n.sender_id = u.id
		LEFT JOIN group_events AS ge ON (n.related_id = ge.id AND n.type = 'group event')
		LEFT JOIN group_members AS gm ON (n.related_id = gm.id AND (n.type = 'group join request' OR n.type = 'group join invitation'))
		LEFT JOIN groups AS g ON g.id = COALESCE(ge.group_id, gm.group_id)
		WHERE n.id = ?;
	`
	var notif Notification
	err := db.QueryRow(query, notifId).Scan(&notif.Id, &notif.Type, &notif.RelatedId, &notif.SenderId, &notif.ReceverId, &notif.GroupId, &notif.SenderName, &notif.GroupName)
	return notif, err
}

// GetNotifications retrieves unread notifications for a user
func GetNotifications(db *sql.DB, userId int) ([]Notification, error) {
	query := `
		SELECT 
		    n.id, 
		    n.type, 
		    n.related_id,
		    n.sender_id,
			n.recever_id,
		    COALESCE(ge.group_id, gm.group_id) AS group_id,
		    u.first_name || ' ' || u.last_name AS sender_name,
		    g.title AS group_name
		FROM notifications AS n
		JOIN users As u ON n.sender_id = u.id
		LEFT JOIN group_events AS ge ON (n.related_id = ge.id AND n.type = 'group event')
		LEFT JOIN group_members AS gm ON (n.related_id = gm.id AND (n.type = 'group join request' OR n.type = 'group join invitation'))
		LEFT JOIN groups AS g ON g.id = COALESCE(ge.group_id, gm.group_id)
		WHERE n.recever_id = ? order by n.id desc;

	`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notif Notification
		err := rows.Scan(&notif.Id, &notif.Type, &notif.RelatedId, &notif.SenderId, &notif.ReceverId, &notif.GroupId, &notif.SenderName, &notif.GroupName)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

// MarkNotificationAsRead updates a notification as read
func MarkNotificationAsRead(db *sql.DB, notifID int) error {
	_, err := db.Exec("UPDATE notifications SET is_read = 1 WHERE recever_id = ?", notifID)
	return err
}

// GetUnreadCount returns the number of unread notifications for a user
func GetUnreadCount(db *sql.DB, userId int) (int, error) {
	var notifsNum int
	err := db.QueryRow("SELECT COUNT(*) FROM notifications WHERE recever_id = ? AND is_read = 0;", userId).Scan(&notifsNum)
	return notifsNum, err
}
