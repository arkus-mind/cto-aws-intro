package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func getConnection() (*sql.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	return sql.Open("postgres", uri)
}

type Todo struct {
	Id            string `json:"tdid"`
	Title         string `json:"tdtitle"`
	Description   string `json:"tddescription"`
	Creation_date string `json:"tdcreation"`
	Updated_date  string `json:"tdupdated"`
	Status        string `json:"tdstatus"`
}

var valid_status = map[string]bool{
	"En Progreso":  true,
	"Terminado":    true,
	"Por Trabajar": true,
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return get_task(req)
	case "POST":
		return create_task(req)
	case "PUT":
		return update_task(req)
	case "DELETE":
		return delete_task(req)
	default:
		apiResp := events.APIGatewayProxyResponse{Body: req.HTTPMethod, StatusCode: 200}
		return apiResp, nil
	}
}

func get_task(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["id"]
	title := req.QueryStringParameters["title"]
	query := "SELECT * FROM public.task LIMIT 1"
	search := ""
	log.Printf(id)
	if id != "" {
		query = "SELECT * FROM public.task WHERE task_id = $1"
		search = id
	} else if title != "" {
		query = "SELECT * FROM public.task WHERE title = $1"
		search = title
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Not found",
		}, nil
	}
	uri := os.Getenv("URI")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}

	if search != "" {
		var task_db Todo
		queryErr := db.QueryRow(query, search).
			Scan(&task_db.Id, &task_db.Title, &task_db.Description, &task_db.Creation_date, &task_db.Updated_date, &task_db.Status)
		if queryErr != nil {
			fmt.Println(queryErr)
		}
		if task_db.Id == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "Not found",
			}, nil
		}

		js, err := json.Marshal(task_db)
		log.Printf(string(js))
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(js),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       http.StatusText(http.StatusNotFound),
	}, nil
}

func create_task(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	title := req.QueryStringParameters["title"]
	description := req.QueryStringParameters["description"]

	if title != "" {
		uri := os.Getenv("URI")
		db, err := sql.Open("postgres", uri)
		if err != nil {
			log.Fatal(err)
		}

		query := `
		INSERT INTO public.task(
			title, description)
			VALUES ($1, $2)
			RETURNING task_id`
		var id string = ""
		err = db.QueryRow(query, title, description).Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Println("New record ID is:", id)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "task created " + id,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       http.StatusText(http.StatusBadRequest),
	}, nil
}

func update_task(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["id"]
	status := req.QueryStringParameters["status"]

	if id != "" && status != "" && valid_status[status] {
		query := `
		UPDATE public.task
		SET updated_date = $2, status = $3
		WHERE task_id = $1;`

		uri := os.Getenv("URI")
		db, err := sql.Open("postgres", uri)
		if err != nil {
			log.Fatal(err)
		}
		res, err := db.Exec(query, id, time.Now(), status)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprint(count) + " rows afected",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       http.StatusText(http.StatusNotFound),
	}, nil
}

func delete_task(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.QueryStringParameters["id"]

	if id != "" {
		query := `
		DELETE FROM public.task
		WHERE task_id = $1;`

		uri := os.Getenv("URI")
		db, err := sql.Open("postgres", uri)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(query, id)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       http.StatusText(http.StatusInternalServerError),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       id + "user deleted",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       http.StatusText(http.StatusNotFound),
	}, nil

}

func main() {
	lambda.Start(router)
}
