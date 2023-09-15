package valid

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func CheckBody(Name string, Quantity int, Unit_coast int) bool { // true - ошибка, false - ошибок нет
	if CheckName(Name) || CheckQuantity(Quantity) || CheckUnitCoast(Unit_coast) {
		return true
	} else {
		return false
	}
}

func CheckID(ID int) bool {
	if ID == 0 {
		return true
	} else {
		return false
	}
}

func CheckName(Name string) bool {
	if len(Name) == 0 {
		return true
	} else {
		return false
	}
}

func CheckQuantity(Quantity int) bool {
	if Quantity == 0 {
		return true
	} else {
		return false
	}
}

func CheckUnitCoast(UnitCoast int) bool {
	if UnitCoast == 0 {
		return true
	} else {
		return false
	}
}

func getDBConnection() (*pgx.Conn, error) {
	config, err := pgx.ParseConfig("postgres://postgres:123@localhost/Web-Service")
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeDBConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}
