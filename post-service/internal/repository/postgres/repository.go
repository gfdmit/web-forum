package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gfdmit/web-forum/post-service/config"
	"github.com/gfdmit/web-forum/post-service/internal/model"
	"github.com/gfdmit/web-forum/post-service/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, conf *config.Postgres) (repository.Repository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.DB,
	)

	migrateURL := strings.Replace(url, "postgresql://", "pgx5://", 1)
	m, err := migrate.New(fmt.Sprintf("file://%v", conf.Migrations), migrateURL)
	if err != nil {
		return nil, fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()

	log.Println("applying migrations...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("nothing to migrate")
		} else {
			return nil, fmt.Errorf("m.Up: %w", err)
		}
	} else {
		log.Println("migrated successfully!")
	}

	poolConf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	poolConf.MaxConns = conf.Pool.MaxConns
	poolConf.MinConns = conf.Pool.MinConns
	poolConf.MaxConnLifetime = conf.Pool.MaxConnLifetime
	poolConf.MaxConnIdleTime = conf.Pool.MaxConnIdleTime
	poolConf.HealthCheckPeriod = conf.Pool.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return &postgresRepository{
		db: pool,
	}, nil
}

func (pr *postgresRepository) GetBoard(ctx context.Context, id int) (model.Board, error) {
	const query = `
        SELECT id, name, description, created_at, deleted_at
        FROM forum.boards
        WHERE id = $1
		AND deleted_at IS NULL
    `
	var board model.Board
	err := pr.db.QueryRow(ctx, query, id).Scan(
		&board.ID, &board.Name, &board.Description, &board.CreatedAt, &board.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Board{}, fmt.Errorf("board not found: %w", repository.ErrNotFound)
		}
		return model.Board{}, fmt.Errorf("GetBoard: %w", err)
	}
	return board, nil
}

func (pr *postgresRepository) GetBoards(ctx context.Context, includeDeleted bool) ([]model.Board, error) {
	const queryAll = `
		SELECT id, name, description, created_at, deleted_at
		FROM forum.boards
	`
	const queryActive = `
		SELECT id, name, description, created_at, deleted_at
		FROM forum.boards
		WHERE deleted_at IS NULL
	`
	query := queryActive
	if includeDeleted {
		query = queryAll
	}
	rows, err := pr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetBoards: %w", err)
	}
	defer rows.Close()

	boards := make([]model.Board, 0)
	for rows.Next() {
		var board model.Board
		err = rows.Scan(&board.ID, &board.Name, &board.Description, &board.CreatedAt, &board.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("GetBoards scan: %w", err)
		}
		boards = append(boards, board)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoards rows: %w", err)
	}

	return boards, nil
}

func (pr *postgresRepository) GetPost(ctx context.Context, id int) (model.Post, error) {
	const query = `
        SELECT 
			p.id, p.user_id, p.board_id, p.title, p.text, p.media_url, p.created_at,
			pr.firstname, pr.lastname
		FROM forum.posts p
		LEFT JOIN forum.profiles pr ON pr.user_id = p.user_id
        WHERE p.id = $1
		AND p.deleted_at IS NULL
    `
	var post model.Post
	var firstname, lastname *string
	err := pr.db.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.BoardID, &post.Title, &post.Text, &post.MediaURL, &post.CreatedAt, &firstname, &lastname,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Post{}, fmt.Errorf("post not found: %w", repository.ErrNotFound)
		}
		return model.Post{}, fmt.Errorf("GetPost: %w", err)
	}
	if firstname != nil && lastname != nil {
		post.Author = &model.Author{
			Firstname: *firstname,
			Lastname:  *lastname,
		}
	}
	return post, nil
}

func (pr *postgresRepository) GetPosts(ctx context.Context, boardID int, includeDeleted bool, limit, offset int) ([]model.Post, error) {
	const queryAll = `
		SELECT 
			p.id, p.user_id, p.board_id, p.title, p.text, p.media_url, p.created_at,
			pr.firstname, pr.lastname
		FROM forum.posts p
		LEFT JOIN forum.profiles pr ON pr.user_id = p.user_id
		WHERE p.board_id = $1
		ORDER BY p.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	const queryActive = `
		SELECT 
			p.id, p.user_id, p.board_id, p.title, p.text, p.media_url, p.created_at,
			pr.firstname, pr.lastname
		FROM forum.posts p
		LEFT JOIN forum.profiles pr ON pr.user_id = p.user_id
		WHERE p.board_id = $1
		AND p.deleted_at IS NULL
		ORDER BY p.updated_at DESC
		LIMIT $2 OFFSET $3
	`

	query := queryActive
	if includeDeleted {
		query = queryAll
	}
	rows, err := pr.db.Query(ctx, query, boardID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetPosts: %w", err)
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		var firstname, lastname *string

		err = rows.Scan(
			&post.ID, &post.UserID, &post.BoardID, &post.Title, &post.Text, &post.MediaURL, &post.CreatedAt,
			&firstname, &lastname,
		)
		if err != nil {
			return nil, fmt.Errorf("GetPosts scan: %w", err)
		}

		if firstname != nil && lastname != nil {
			post.Author = &model.Author{
				Firstname: *firstname,
				Lastname:  *lastname,
			}
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetPosts rows: %w", err)
	}

	return posts, nil
}

func (pr *postgresRepository) GetComment(ctx context.Context, id int) (model.Comment, error) {
	const query = `
        SELECT 
			c.id, c.user_id, c.post_id, c.text, c.media_url, c.created_at, 
			pr.firstname, pr.lastname
        FROM forum.comments c
		LEFT JOIN forum.profiles pr ON pr.user_id = c.user_id
        WHERE c.id = $1
		AND c.deleted_at IS NULL
    `
	var comment model.Comment
	var firstname, lastname *string
	err := pr.db.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.UserID, &comment.PostID, &comment.Text, &comment.MediaURL, &comment.CreatedAt, &firstname, &lastname,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Comment{}, fmt.Errorf("comment not found: %w", repository.ErrNotFound)
		}
		return model.Comment{}, fmt.Errorf("GetComment: %w", err)
	}
	if firstname != nil && lastname != nil {
		comment.Author = &model.Author{
			Firstname: *firstname,
			Lastname:  *lastname,
		}
	}
	return comment, nil
}

func (pr *postgresRepository) GetComments(ctx context.Context, postID int, includeDeleted bool, limit, offset int) ([]model.Comment, error) {
	const queryAll = `
		SELECT 
			c.id, c.user_id, c.post_id, c.text, c.media_url, c.created_at, 
			pr.firstname, pr.lastname
        FROM forum.comments c
		LEFT JOIN forum.profiles pr ON pr.user_id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at
		LIMIT $2 OFFSET $3
	`
	const queryActive = `
		SELECT 
			c.id, c.user_id, c.post_id, c.text, c.media_url, c.created_at, 
			pr.firstname, pr.lastname
        FROM forum.comments c
		LEFT JOIN forum.profiles pr ON pr.user_id = c.user_id
		WHERE c.post_id = $1
		AND c.deleted_at IS NULL
		ORDER BY c.created_at
		LIMIT $2 OFFSET $3
	`

	query := queryActive
	if includeDeleted {
		query = queryAll
	}

	rows, err := pr.db.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetComments: %w", err)
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		var firstname, lastname *string
		err = rows.Scan(
			&comment.ID, &comment.UserID, &comment.PostID, &comment.Text, &comment.MediaURL, &comment.CreatedAt, &firstname, &lastname,
		)
		if err != nil {
			return nil, fmt.Errorf("GetComments scan: %w", err)
		}
		if firstname != nil && lastname != nil {
			comment.Author = &model.Author{
				Firstname: *firstname,
				Lastname:  *lastname,
			}
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetComments rows: %w", err)
	}

	return comments, nil
}

func (pr *postgresRepository) CreateBoard(ctx context.Context, input model.CreateBoardInput) (model.Board, error) {
	const query = `
        INSERT INTO forum.boards (name, description)
        VALUES ($1, $2)
        RETURNING id, name, description, created_at, deleted_at
    `
	var board model.Board
	err := pr.db.QueryRow(ctx, query, input.Name, input.Description).Scan(
		&board.ID, &board.Name, &board.Description, &board.CreatedAt, &board.DeletedAt,
	)
	if err != nil {
		return model.Board{}, fmt.Errorf("CreateBoard: %w", err)
	}
	return board, nil
}

func (pr *postgresRepository) DeleteBoard(ctx context.Context, id int) error {
	const query = `
		UPDATE forum.boards
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`
	result, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteBoard: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DeleteBoard: %w", repository.ErrNotFound)
	}

	return nil
}

func (pr *postgresRepository) RestoreBoard(ctx context.Context, id int) error {
	const query = `
		UPDATE forum.boards
		SET deleted_at = NULL
		WHERE id = $1
		AND deleted_at IS NOT NULL
	`
	result, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("RestoreBoard: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("RestoreBoard: %w", repository.ErrNotFound)
	}

	return nil
}

func (pr *postgresRepository) CreatePost(ctx context.Context, input model.CreatePostInput) (model.Post, error) {
	const query = `
		INSERT INTO forum.posts (user_id, board_id, title, text, media_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, board_id, title, text, media_url, created_at, deleted_at
	`
	var post model.Post
	err := pr.db.QueryRow(ctx, query, input.UserID, input.BoardID, input.Title, input.Text, input.MediaURL).Scan(
		&post.ID, &post.UserID, &post.BoardID, &post.Title, &post.Text, &post.MediaURL, &post.CreatedAt, &post.DeletedAt,
	)
	if err != nil {
		return model.Post{}, fmt.Errorf("CreatePost: %w", err)
	}
	return post, nil
}

func (pr *postgresRepository) DeletePost(ctx context.Context, id int) error {
	const query = `
		UPDATE forum.posts
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`
	result, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeletePost: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DeletePost: %w", repository.ErrNotFound)
	}

	return nil
}

func (pr *postgresRepository) CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error) {
	const query = `
		INSERT INTO forum.comments (user_id, post_id, text, media_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, post_id, text, media_url, created_at, deleted_at
	`
	var comment model.Comment
	err := pr.db.QueryRow(ctx, query, input.UserID, input.PostID, input.Text, input.MediaURL).Scan(
		&comment.ID, &comment.UserID, &comment.PostID, &comment.Text, &comment.MediaURL, &comment.CreatedAt, &comment.DeletedAt,
	)
	if err != nil {
		return model.Comment{}, fmt.Errorf("CreateComment: %w", err)
	}
	return comment, nil
}

func (pr *postgresRepository) DeleteComment(ctx context.Context, id int) error {
	const query = `
		UPDATE forum.comments
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`
	result, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteComment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DeleteComment: %w", repository.ErrNotFound)
	}

	return nil
}

func (pr *postgresRepository) GetProfile(ctx context.Context, userID int) (model.Profile, error) {
	const query = `
        SELECT id, user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status
        FROM forum.profiles
        WHERE user_id = $1
		AND deleted_at IS NULL
    `
	var profile model.Profile
	err := pr.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.UniversityID, &profile.Firstname, &profile.Lastname, &profile.Middlename, &profile.Birthday, &profile.Faculty, &profile.Grade, &profile.Group, &profile.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Profile{}, fmt.Errorf("profile not found: %w", repository.ErrNotFound)
		}
		return model.Profile{}, fmt.Errorf("GetProfile: %w", err)
	}
	return profile, nil
}

func (pr *postgresRepository) GetProfiles(ctx context.Context, includeDeleted bool) ([]model.Profile, error) {
	const queryAll = `
        SELECT id, user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status
        FROM forum.profiles
    `
	const queryActive = `
        SELECT id, user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status
        FROM forum.profiles
        WHERE deleted_at IS NULL
    `

	query := queryActive
	if includeDeleted {
		query = queryAll
	}
	rows, err := pr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetProfiles: %w", err)
	}
	defer rows.Close()

	profiles := make([]model.Profile, 0)
	for rows.Next() {
		var profile model.Profile
		err = rows.Scan(
			&profile.ID, &profile.UserID, &profile.UniversityID, &profile.Firstname, &profile.Lastname, &profile.Middlename, &profile.Birthday, &profile.Faculty, &profile.Grade, &profile.Group, &profile.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("GetProfiles scan: %w", err)
		}
		profiles = append(profiles, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetProfiles rows: %w", err)
	}

	return profiles, nil
}
