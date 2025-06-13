package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JoaoPedroVicentin/sqlc/internal/db"
	_ "github.com/go-sql-driver/mysql"
)

type CourseDB struct {
	dbConnection *sql.DB
	*db.Queries
}

func NewCourseDB(dbConnection *sql.DB) *CourseDB {
	return &CourseDB{
		dbConnection: dbConnection,
		Queries:      db.New(dbConnection),
	}
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price       float64
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := db.New(tx)
	err = fn(queries)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("failed to rollback transaction: %v, original error: %w", errRb, err)
		}
		return err
	}

	return tx.Commit()
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, courseParams CourseParams, categoryParams CategoryParams) error {
	err := c.callTx(ctx, func(queries *db.Queries) error {
		var error error
		error = queries.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          categoryParams.ID,
			Name:        categoryParams.Name,
			Description: categoryParams.Description,
		})
		if error != nil {
			return error
		}

		error = queries.CreateCourse(ctx, db.CreateCourseParams{
			ID:          courseParams.ID,
			CategoryID:  categoryParams.ID,
			Name:        courseParams.Name,
			Description: courseParams.Description,
			Price:       courseParams.Price,
		})
		if error != nil {
			return error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	dbConnection, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	queries := db.New(dbConnection)

	courses, err := queries.ListCourses(ctx)
	if err != nil {
		panic(err)
	}
	for _, course := range courses {
		fmt.Printf("Course ID: %s, Name: %s, Category: %s, Price: %.2f\n, Description: %s\n",
			course.ID, course.Name, course.CategoryName, course.Price, course.Description.String)
	}

	// courseParams := CourseParams{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Go",
	// 	Description: sql.NullString{String: "Learn the basics of GO", Valid: true},
	// 	Price:       19.99,
	// }

	// categoryParams := CategoryParams{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Backend",
	// 	Description: sql.NullString{String: "Courses related to backend", Valid: true},
	// }

	// courseDb := NewCourseDB(dbConnection)
	// err = courseDb.CreateCourseAndCategory(ctx, courseParams, categoryParams)
	// if err != nil {
	// 	panic(err)
	// }
}
