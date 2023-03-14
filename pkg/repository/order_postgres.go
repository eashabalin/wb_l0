package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"wb_l0/pkg/model"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}

func (r *OrderPostgres) Create(order model.Order) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	fmt.Printf("%+v\n", order)

	// save to postgres
	id, err := r.createPayment(tx, order.Payment, order.OrderUID)
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}
	fmt.Println("id: ", id)

	return "", tx.Commit()
}

func (r *OrderPostgres) getDeliveryServiceIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where delivery_service = '%s'", deliveryServiceTable, name)
	row := r.db.QueryRow(query)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createDeliveryService(tx *sql.Tx, name string) (int, error) {
	createDeliveryServiceQuery := fmt.Sprintf("INSERT into %s (delivery_service) VALUES ($1) RETURNING id",
		deliveryServiceTable)
	row := tx.QueryRow(createDeliveryServiceQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getLocaleIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where locale = '%s'", localesTable, name)
	row := r.db.QueryRow(query)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createLocale(tx *sql.Tx, name string) (int, error) {
	createLocaleQuery := fmt.Sprintf("INSERT into %s (locale) VALUES ($1) RETURNING id",
		localesTable)
	row := tx.QueryRow(createLocaleQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getBankIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where bank = '%s'", banksTable, name)
	row := r.db.QueryRow(query)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createBank(tx *sql.Tx, name string) (int, error) {
	createBankQuery := fmt.Sprintf("INSERT into %s (bank) VALUES ($1) RETURNING id",
		banksTable)
	row := tx.QueryRow(createBankQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getRegionIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where region = $1", regionsTable)
	row := r.db.QueryRow(query, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createRegion(tx *sql.Tx, name string) (int, error) {
	createRegionQuery := fmt.Sprintf("INSERT into %s (region) VALUES ($1) RETURNING id",
		regionsTable)
	row := tx.QueryRow(createRegionQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getPaymentProviderIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where provider = $1", paymentProvidersTable)
	row := r.db.QueryRow(query, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createPaymentProvider(tx *sql.Tx, name string) (int, error) {
	createPaymentProviderQuery := fmt.Sprintf("INSERT into %s (provider) VALUES ($1) RETURNING id",
		paymentProvidersTable)
	row := tx.QueryRow(createPaymentProviderQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getBrandIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where brand = $1", brandsTable)
	row := r.db.QueryRow(query, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createBrand(tx *sql.Tx, name string) (int, error) {
	createBrandQuery := fmt.Sprintf("INSERT into %s (brand) VALUES ($1) RETURNING id",
		brandsTable)
	row := tx.QueryRow(createBrandQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getCurrencyIDByName(name string) (int, error) {
	query := fmt.Sprintf("select * from %s where currency = &1", currenciesTable)
	row := r.db.QueryRow(query, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createCurrency(tx *sql.Tx, name string) (int, error) {
	createCurrencyQuery := fmt.Sprintf("INSERT into %s (currency) VALUES ($1) RETURNING id",
		currenciesTable)
	row := tx.QueryRow(createCurrencyQuery, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) getCityID(city, region string) (int, error) {
	regionID, err := r.getRegionIDByName(region)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("select id from %s where city=$1 and region_id=$2", citiesTable)
	row := r.db.QueryRow(query, city, regionID)
	var id int
	err = row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createCity(tx *sql.Tx, city string, region string) (int, error) {
	regionID, err := r.getRegionIDByName(region)
	if err != nil {
		regionID, err = r.createRegion(tx, region)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	query := fmt.Sprintf("insert into %s (city, region_id) values ($1, $2) returning id", citiesTable)
	row := tx.QueryRow(query, city, regionID)
	var id int
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, nil
}

func (r *OrderPostgres) createDelivery(tx *sql.Tx, delivery model.Delivery, orderID string) (string, error) {
	query := fmt.Sprintf("insert into %s (order_uid, name, phone, zip, city_id, address, email) values ($1, $2, $3, $4, $5, $6, $7) returning order_uid", deliveriesTable)
	cityID, err := r.getCityID(delivery.City, delivery.Region)
	if err != nil {
		cityID, err = r.createCity(tx, delivery.City, delivery.Region)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	row := tx.QueryRow(query, orderID, delivery.Name, delivery.Phone, delivery.Zip, cityID, delivery.Address, delivery.Email)
	var id string
	err = row.Scan(&id)
	return id, err
}

func (r *OrderPostgres) createPayment(tx *sql.Tx, payment model.Payment, orderUID string) (uid string, err error) {
	query := fmt.Sprintf(`insert into %s (order_uid, transaction, request_id, currency_id, provider_id, amount,
												 payment_dt, bank_id, delivery_cost, goods_total, custom_fee)
								  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning order_uid`, paymentsTable)
	providerID, err := r.getPaymentProviderIDByName(payment.Provider)
	if err != nil {
		providerID, err = r.createPaymentProvider(tx, payment.Provider)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	currencyID, err := r.getCurrencyIDByName(payment.Currency)
	if err != nil {
		currencyID, err = r.createCurrency(tx, payment.Currency)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	bankID, err := r.getBankIDByName(payment.Bank)
	if err != nil {
		bankID, err = r.createBank(tx, payment.Bank)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	row := tx.QueryRow(query, orderUID, payment.Transaction, payment.RequestID, currencyID, providerID,
		payment.Amount, payment.PaymentDatetime, bankID, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)

	err = row.Scan(&uid)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return uid, err
}
