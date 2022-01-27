package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JorgeAdd/intro5/homeworks-jorge/homeworksAPI/internal/database"
	"github.com/gorilla/mux"
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

type HomeworkDeleteSt struct {
	Id string `json:"deletedId"`
}

func getConnection() *sql.DB {
	return database.GetDB()
}

func GetHomework(w http.ResponseWriter, r *http.Request) {

	db := getConnection()
	defer db.Close()
	var homework HomeworkSt
	var errHw ErrSt

	params := mux.Vars(r)
	homework.Id = params["homeworkId"]

	homeworkRes := db.QueryRow("SELECT * FROM homeworks.homeworks WHERE id = '" + homework.Id + "'")
	err := homeworkRes.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)
	if err == nil {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(homework)

	} else {

		w.Header().Set("Content-Type", "application/json")
		errHw.Error = "Homework id not found"
		json.NewEncoder(w).Encode(errHw)

	}

}

func GetHomeworks(w http.ResponseWriter, r *http.Request) {

	db := getConnection()
	defer db.Close()
	var homework HomeworkSt
	var homeworks []HomeworkSt

	homeworksRes, errHw := db.Query("SELECT * FROM homeworks.homeworks")
	if errHw != nil {
		panic(errHw)
	}
	defer homeworksRes.Close()
	for homeworksRes.Next() {
		errHw = homeworksRes.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)
		if errHw != nil {
			panic(errHw)
		}
		homeworks = append(homeworks, homework)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homeworks)
}

func CreateHomework(w http.ResponseWriter, r *http.Request) {

	db := getConnection()
	defer db.Close()
	var homework HomeworkSt
	var errHw ErrSt

	json.NewDecoder(r.Body).Decode(&homework)
	homework.Created_at = time.Now().Format("2006-01-02")
	homework.Updated_at = time.Now().Format("2006-01-02")
	homework.Status = "To do"

	if homework.Title != "" {

		insertHomework := db.QueryRow("INSERT INTO homeworks.homeworks (id,title,description,created_at,updated_at,status) VALUES" +
			"(gen_random_uuid(),'" + homework.Title + "','" + homework.Description + "','" + homework.Created_at + "','" + homework.Updated_at + "','" + homework.Status + "') RETURNING id")
		errInsHw := insertHomework.Scan(&homework.Id)
		if errInsHw != nil {
			panic(errInsHw)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(homework)

	} else {

		w.Header().Set("Content-Type", "application/json")
		errHw.Error = "Title shouldn't be empty"
		json.NewEncoder(w).Encode(errHw)

	}
}

func UpdateHomework(w http.ResponseWriter, r *http.Request) {
	db := getConnection()
	defer db.Close()
	var homework HomeworkSt
	var homeworkBody HomeworkSt
	var errHw ErrSt

	params := mux.Vars(r)
	homeworkBody.Id = params["homeworkId"]
	json.NewDecoder(r.Body).Decode(&homeworkBody)

	homeworkRes := db.QueryRow("SELECT * FROM homeworks.homeworks WHERE id = '" + homeworkBody.Id + "'")
	err := homeworkRes.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)
	if err == nil {
		if homeworkBody.Title == "" || homeworkBody.Status == "" {

			w.Header().Set("Content-Type", "application/json")
			errHw.Error = "Homework title and homework status shouldn't be empty"
			json.NewEncoder(w).Encode(errHw)

		} else {

			homeworkBody.Updated_at = time.Now().Format("2006-01-02")

			updatedHomework := db.QueryRow("UPDATE homeworks.homeworks SET " +
				"title = '" + homeworkBody.Title + "', " +
				"description = '" + homeworkBody.Description + "', " +
				"updated_at = '" + homeworkBody.Updated_at + "', " +
				"status = '" + homeworkBody.Status + "' " +
				"where id = '" + homeworkBody.Id + "' RETURNING id")
			errUpdtHw := updatedHomework.Scan(&homework.Id)
			if errUpdtHw != nil {
				panic(errUpdtHw)
			}

			//Select with new values
			homeworkNewValues := db.QueryRow("SELECT * FROM homeworks.homeworks WHERE id = '" + homework.Id + "'")
			homeworkNewValues.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(homework)

		}
	} else {

		w.Header().Set("Content-Type", "application/json")
		errHw.Error = "Homework id not found"
		json.NewEncoder(w).Encode(errHw)

	}
}

func DeleteHomework(w http.ResponseWriter, r *http.Request) {
	db := getConnection()
	defer db.Close()
	var homework HomeworkSt
	var deletedHomework HomeworkDeleteSt
	var errHw ErrSt

	params := mux.Vars(r)
	homework.Id = params["homeworkId"]

	homeworkRes := db.QueryRow("SELECT * FROM homeworks.homeworks WHERE id = '" + homework.Id + "'")
	err := homeworkRes.Scan(&homework.Id, &homework.Title, &homework.Description, &homework.Created_at, &homework.Updated_at, &homework.Status)

	if err == nil {

		deletedHomeworkScript := db.QueryRow("DELETE FROM homeworks.homeworks " +
			"WHERE id = '" + homework.Id + "' RETURNING id")
		errDelHw := deletedHomeworkScript.Scan(&deletedHomework.Id)
		if errDelHw == nil {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(deletedHomework)

		} else {
			w.Header().Set("Content-Type", "application/json")
			errHw.Error = "Could not delete homework"
			json.NewEncoder(w).Encode(errHw)

		}

	} else {

		w.Header().Set("Content-Type", "application/json")
		errHw.Error = "Homework id not found"
		json.NewEncoder(w).Encode(errHw)

	}
}
