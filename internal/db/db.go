package db

import (
	"assignment/internal/config"
	employeeerror "assignment/internal/errors"
	"assignment/internal/models"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type postgres struct{ db *sql.DB }

type EmployeeDBService interface {
	// EmployeeDBService
	CreateEmployee(*gin.Context, models.Employee) (string, *employeeerror.EmployeeError)
	DeleteEmployee(*gin.Context, models.Employee) *employeeerror.EmployeeError
	GetEmployeeByID(*gin.Context) ([]models.Employee, *employeeerror.EmployeeError)
	UpdateEmployee(*gin.Context) (string, *employeeerror.EmployeeError)
}

func (p postgres) CreateEmployee(ctx *gin.Context, notes models.Employee) (string, *employeeerror.EmployeeError) {
	return "", nil
}
func (p postgres) DeleteEmployee(ctx *gin.Context, notes models.Employee) *employeeerror.EmployeeError {
	return nil
}

func (p postgres) GetEmployeeByID(ctx *gin.Context) ([]models.Employee, *employeeerror.EmployeeError) {
	return []models.Employee{}, nil
}
func (p postgres) UpdateEmployee(ctx *gin.Context) (string, *employeeerror.EmployeeError) {
	return "", nil
}

func New() (postgres, error) {
	cfg := config.GetConfig()
	connString := "host=" + cfg.Database.Host + " " + "dbname=" + cfg.Database.DBname + " " + "password=" +
		cfg.Database.Password + " " + "user=" + cfg.Database.User + " " + "port=" + fmt.Sprint(cfg.Database.Port)

	conn, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to connect: %v\n", err))
		return postgres{}, err
	}

	log.Println("Connected to database")

	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot Ping the database")
		return postgres{}, err
	}
	log.Println("pinged database")

	return postgres{db: conn}, nil
}
