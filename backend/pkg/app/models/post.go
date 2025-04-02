package models

import (
	"database/sql"
	"html"
	"time"
)

type Post struct {
	Pid             int
	Uid             int
	Username        string
	Avatar          *string
	Image           *string
	Content         string
	Privacy         string
	CreatedAt       string
	Likes           int
	Dislikes        int
	NbComment       int
	UserInteraction int
}

// CreatePost inserts a new post into the database and returns the post ID
func CreatePost(db *sql.DB, userID int, content, privacy string, imagePath *string, groupID int) (int, error) {
	var postID int
	err := db.QueryRow(`
		INSERT INTO posts (user_id, content, privacy, image_path, group_id)
		VALUES (?, ?, ?, ?, ?) RETURNING id`,
		userID, html.EscapeString(content), privacy, imagePath, groupID,
	).Scan(&postID)

	if err != nil {
		return 0, err
	}
	return postID, nil
}

// AddPostPrivacy inserts privacy settings for selected users
func AddPostPrivacy(db *sql.DB, postID int, selectedUsers []int) error {
	if len(selectedUsers) == 0 {
		return nil
	}

	for _, userID := range selectedUsers {
		_, err := db.Exec(`INSERT INTO post_privacy_users (user_id, post_id, created_at) VALUES (?, ?, ?)`, userID, postID, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUserPostCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ? AND group_id = 0", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetPostsWithPrivacy(userID int, baseWhere string, args []interface{}, offset int, db *sql.DB) ([]Post, error) {
	privacyWhere := `
        AND (
            posts.privacy = 'public' 
            OR posts.user_id = ? 
            OR (posts.privacy = 'almost_private' AND EXISTS (
                SELECT 1 FROM follow_requests 
                WHERE follower_id = ? AND following_id = posts.user_id AND status = 'approved'
            )) 
            OR (posts.privacy = 'private' AND EXISTS (
                SELECT 1 FROM post_privacy_users 
                WHERE post_id = posts.id AND user_id = ?
            ))
        )`

	args = append(args, userID, userID, userID)

	query := `
        SELECT posts.id, posts.user_id, posts.created_at, posts.content, 
               posts.image_path, posts.privacy, 
               users.first_name, users.last_name, users.avatar,
               (SELECT COUNT(*) FROM comments WHERE post_id = posts.id) AS comment_count
        FROM posts
        JOIN users ON posts.user_id = users.id
        ` + baseWhere + privacyWhere + `
        ORDER BY posts.id DESC LIMIT 4 OFFSET ?`
	args = append(args, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var firstName, lastName string

		err := rows.Scan(&post.Pid, &post.Uid, &timeCreated, &post.Content,
			&post.Image, &post.Privacy, &firstName, &lastName, &post.Avatar, &post.NbComment)
		if err != nil {
			return nil, err
		}

		post.Content = html.UnescapeString(post.Content)
		post.Username = firstName + " " + lastName
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
		post.Likes = getPostLikeDisLike(post.Pid, db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, db)

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPost(db *sql.DB, postID, userID int) (Post, error) {
	var p Post
	var timeCreated time.Time
	var fistname, lastname string
	err := db.QueryRow(`SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.image_path,users.first_name, users.last_name, users.avatar FROM posts
        JOIN users ON posts.user_id = users.id WHERE posts.id = ? `, postID).Scan(&p.Pid, &p.Uid, &timeCreated, &p.Content, &p.Image, &fistname, &lastname, &p.Avatar)
	if err != nil {
		return Post{}, err
	}
	p.Content = html.UnescapeString(p.Content)
	p.Username = fistname + " " + lastname
	err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", p.Pid).Scan(&p.NbComment)
	if err != nil {
		return Post{}, err
	}
	p.Likes = getPostLikeDisLike(p.Pid, db)
	p.UserInteraction = getPostInteractions(p.Pid, userID, db)
	p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
	return p, nil
}

func getPostLikeDisLike(post_id int, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id=? AND interaction=1", post_id).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func getPostInteractions(pid, uid int, db *sql.DB) int {
	var ii int
	db.QueryRow(`select interaction from post_interactions where user_id = ? and post_id=?`, uid, pid).Scan(&ii)
	return ii
}

// InsertOrUpdateLike handles inserting or updating likes in the database
func InsertOrUpdateLike(db *sql.DB, userID, postID, likeStatus int) error {
	query := `
		INSERT INTO post_interactions (user_id, post_id, interaction)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET interaction=excluded.interaction
	`
	_, err := db.Exec(query, userID, postID, likeStatus)
	return err
}
