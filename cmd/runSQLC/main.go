package main

import (
	"context"
	"database/sql"

	"github.com/JoaoPedroVicentin/sqlc/internal/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConnection, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	queries := db.New(dbConnection)

	// err = queries.CreateCategory(ctx, db.CreateCategoryParams{
	// 	ID:   uuid.New().String(),
	// 	Name: "Category 1",
	// 	Description: sql.NullString{
	// 		String: "This is a description",
	// 		Valid:  true,
	// 	},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// err = queries.UpdateCategory(ctx, db.UpdateCategoryParams{
	// 	ID:          "a545e280-e654-4513-9d9d-1a264671d075",
	// 	Name:        "Updated Category Name",
	// 	Description: sql.NullString{String: "Updated description", Valid: true},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	err = queries.DeleteCategory(ctx, "a545e280-e654-4513-9d9d-1a264671d075")
	if err != nil {
		panic(err)
	}

	categories, err := queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		println("Category Name:", category.Name)
	}
}
