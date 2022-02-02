package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arkus-mind/intro5/homeworks-lambda-jorge/homeworksLambdaAPI/internal/database"
)

type HomeworkSt struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Status      string `json:"status"`
}

type ErrSt struct {
	Error string `json:"error"`
}

type ReqBody struct {
	Title string `json:"title"`
}

func getConnection() *sql.DB {
	return database.GetDB()
}

func GetHomeworkByTitle(ctx context.Context, title ReqBody) (HomeworkSt, error) {

	db := getConnection()
	defer db.Close()
	var homework HomeworkSt

	homeworkRes := db.QueryRow("SELECT * FROM homeworks.homeworks WHERE title = '" + title.Title + "'")
	err := homeworkRes.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)
	if err == nil {
		fmt.Println("Homework with title " + title.Title + "found, id: " + homework.Id)
		return homework, nil
	} else {
		fmt.Println("Homework with title " + title.Title + "doesn't exists")
		homework.Title = "not found"
		return homework, nil
	}

}
