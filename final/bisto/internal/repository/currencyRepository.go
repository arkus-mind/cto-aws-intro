package repository

import (
	"bisto/internal/models"
	"database/sql"
)

type currencyRepository struct {
	db *sql.DB
}

type CurrencyRepository interface {
	NewCurrency(currency models.Currency) string
	ExistCurrency(IdCrypto string) bool
	GetCurrenciesByDate(dateIni string, dateEnd string) ([]models.Currency, error)
	GetCurrenciesByType(filter string) ([]models.Currency, error)
	GetCurrenciesByAllParams(dateIni string, dateEnd string, filter string) ([]models.Currency, error)
	GetAllCurrencies() ([]models.Currency, error)
	CloseConnection() bool
}

func (ur *currencyRepository) NewCurrency(currency models.Currency) string {
	//TODO: validate the new schedule not exist into data base.
	// close database
	//defer ur.db.Close()

	insertStmt := `INSERT INTO "practices"."Currency" ("Id", "IdCrypto", "CreatedAt", "Book", "Volume", "High", "Last", "Low", "Vwap", "Ask", "Bid", "Change_24", "USDToMXN", "HKDToMXN") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING "Id"`
	var Id string
	// Scan function will save the insert id in the id
	err := ur.db.QueryRow(insertStmt, currency.Id, currency.IdCrypto, currency.CreatedAt, currency.Book, currency.Volume, currency.High, currency.Last, currency.Low, currency.Vwap, currency.Ask, currency.Bid, currency.Change_24, currency.USDToMXN, currency.HKDToMXN).Scan(&Id)
	if err != nil {
		panic(err)
	}
	return Id
}

func (ur *currencyRepository) ExistCurrency(IdCrypto string) bool {
	var crypto models.Currency
	sqlStatement := `SELECT "Id", "IdCrypto" FROM "practices"."Currency" WHERE "IdCrypto"=$1`
	rows := ur.db.QueryRow(sqlStatement, IdCrypto)
	err := rows.Scan(&crypto.Id, &crypto.IdCrypto)
	return err == nil
}

func (ur *currencyRepository) GetCurrenciesByDate(dateIni string, dateEnd string) ([]models.Currency, error) {
	var currencies []models.Currency
	// close database
	defer ur.db.Close()
	// create the select sql query
	sqlStatement := `SELECT * FROM "practices"."Currency" WHERE "CreatedAt"  >= $1 AND "CreatedAt" < $2`
	// execute the sql statement
	rows, err := ur.db.Query(sqlStatement, dateIni, dateEnd)
	CheckError(err)
	// close the statement
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var currency models.Currency
		// unmarshal the row object to user
		err = rows.Scan(&currency.Id, &currency.IdCrypto, &currency.CreatedAt, &currency.Book, &currency.Volume, &currency.High, &currency.Last, &currency.Low, &currency.Vwap, &currency.Ask, &currency.Bid, &currency.Change_24, &currency.USDToMXN, &currency.HKDToMXN)
		CheckError(err)
		currencies = append(currencies, currency)
	}
	return currencies, err
}

func (ur *currencyRepository) GetCurrenciesByType(filter string) ([]models.Currency, error) {
	currencies := []models.Currency{}
	// close database
	defer ur.db.Close()
	// create the select sql query
	sqlStatement := `SELECT * FROM "practices"."Currency"`
	// execute the sql statement
	rows, err := ur.db.Query(sqlStatement)
	CheckError(err)
	// close the statement
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var currency models.Currency
		// unmarshal the row object to user
		err = rows.Scan(&currency.Id, &currency.IdCrypto, &currency.CreatedAt, &currency.Book, &currency.Volume, &currency.High, &currency.Last, &currency.Low, &currency.Vwap, &currency.Ask, &currency.Bid, &currency.Change_24, &currency.USDToMXN, &currency.HKDToMXN)
		CheckError(err)
		if filter == "USD" {
			currency.Volume = currency.Volume / currency.USDToMXN
			currency.High = currency.High / currency.USDToMXN
			currency.Last = currency.Last / currency.USDToMXN
			currency.Low = currency.Low / currency.USDToMXN
			currency.Vwap = currency.Vwap / currency.USDToMXN
			currency.Ask = currency.Ask / currency.USDToMXN
			currency.Bid = currency.Bid / currency.USDToMXN
			currency.Change_24 = currency.Change_24 / currency.USDToMXN
		} else if filter == "HDK" {
			currency.Volume = currency.Volume / currency.HKDToMXN
			currency.High = currency.High / currency.HKDToMXN
			currency.Last = currency.Last / currency.HKDToMXN
			currency.Low = currency.Low / currency.HKDToMXN
			currency.Vwap = currency.Vwap / currency.HKDToMXN
			currency.Ask = currency.Ask / currency.HKDToMXN
			currency.Bid = currency.Bid / currency.HKDToMXN
			currency.Change_24 = currency.Change_24 / currency.HKDToMXN
		}
		currencies = append(currencies, currency)
	}
	return currencies, err
}

func (ur *currencyRepository) GetCurrenciesByAllParams(dateIni string, dateEnd string, filter string) ([]models.Currency, error) {
	currencies := []models.Currency{}
	// close database
	defer ur.db.Close()
	// create the select sql query
	sqlStatement := `SELECT * FROM "practices"."Currency" WHERE "CreatedAt"  >= $1 AND "CreatedAt" < $2`
	// execute the sql statement
	rows, err := ur.db.Query(sqlStatement, dateIni, dateEnd)
	CheckError(err)
	// close the statement
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var currency models.Currency
		// unmarshal the row object to user
		err = rows.Scan(&currency.Id, &currency.IdCrypto, &currency.CreatedAt, &currency.Book, &currency.Volume, &currency.High, &currency.Last, &currency.Low, &currency.Vwap, &currency.Ask, &currency.Bid, &currency.Change_24, &currency.USDToMXN, &currency.HKDToMXN)
		CheckError(err)

		if filter == "USD" {
			currency.Volume = currency.Volume / currency.USDToMXN
			currency.High = currency.High / currency.USDToMXN
			currency.Last = currency.Last / currency.USDToMXN
			currency.Low = currency.Low / currency.USDToMXN
			currency.Vwap = currency.Vwap / currency.USDToMXN
			currency.Ask = currency.Ask / currency.USDToMXN
			currency.Bid = currency.Bid / currency.USDToMXN
			currency.Change_24 = currency.Change_24 / currency.USDToMXN
		} else if filter == "HDK" {
			currency.Volume = currency.Volume / currency.HKDToMXN
			currency.High = currency.High / currency.HKDToMXN
			currency.Last = currency.Last / currency.HKDToMXN
			currency.Low = currency.Low / currency.HKDToMXN
			currency.Vwap = currency.Vwap / currency.HKDToMXN
			currency.Ask = currency.Ask / currency.HKDToMXN
			currency.Bid = currency.Bid / currency.HKDToMXN
			currency.Change_24 = currency.Change_24 / currency.HKDToMXN
		}

		currencies = append(currencies, currency)
	}
	return currencies, err
}

func (ur *currencyRepository) GetAllCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency
	// close database
	defer ur.db.Close()
	// create the select sql query
	sqlStatement := `SELECT * FROM "practices"."Currency"`
	// execute the sql statement
	rows, err := ur.db.Query(sqlStatement)
	CheckError(err)
	// close the statement
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var currency models.Currency
		// unmarshal the row object to user
		err = rows.Scan(&currency.Id, &currency.IdCrypto, &currency.CreatedAt, &currency.Book, &currency.Volume, &currency.High, &currency.Last, &currency.Low, &currency.Vwap, &currency.Ask, &currency.Bid, &currency.Change_24, &currency.USDToMXN, &currency.HKDToMXN)
		CheckError(err)
		currencies = append(currencies, currency)
	}
	return currencies, err
}

func (ur *currencyRepository) CloseConnection() bool {
	err := ur.db.Close()
	return err == nil
}

func NewCurrencyRepository() CurrencyRepository {
	return &currencyRepository{db: CreateConnection()}
}
