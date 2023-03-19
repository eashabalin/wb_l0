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

	deliveryServiceID, err := r.getDeliveryServiceIDByName(order.DeliveryService)
	if err != nil {
		deliveryServiceID, err = r.createDeliveryService(tx, order.DeliveryService)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	localeID, err := r.getLocaleIDByName(order.Locale)
	if err != nil {
		localeID, err = r.createLocale(tx, order.Locale)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	query := fmt.Sprintf("insert into %s (order_uid, track_number, entry, locale_id, internal_signature, customer_id, delivery_service_id, shardkey, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", ordersTable)

	_, err = tx.Exec(
		query,
		order.UID,
		order.TrackNumber,
		order.Entry, localeID,
		order.InternalSignature,
		order.CustomerID, deliveryServiceID,
		order.ShardKey, order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = r.createDelivery(tx, order.Delivery, order.UID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = r.createPayment(tx, order.Payment, order.UID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	for _, item := range order.Items {
		_, _, err = r.createOrderItem(tx, item, order.UID)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	return order.UID, tx.Commit()
}

func (r *OrderPostgres) GetByUID(uid string) (*model.Order, bool) {
	var order model.Order

	query := fmt.Sprintf(`select o.order_uid, o.track_number, o.entry, l.locale, o.internal_signature, o.customer_id, ds.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
       								d.name d_name, d.phone d_phone, d.zip d_zip, ct.city d_city, d.address d_address, r.region d_region, d.email d_email,
       								p.transaction p_transaction, p.request_id p_request_id, c.currency p_currency, pp.provider p_provider, p.amount p_amount, p.payment_dt p_payment_dt, b.bank p_bank, p.delivery_cost p_delivery_cost, p.goods_total p_goods_total, p.custom_fee p_custom_fee
								from orders o
    								inner join delivery_services ds on ds.id = o.delivery_service_id
    								inner join locales l on l.id = o.locale_id
    								inner join deliveries d on o.order_uid = d.order_uid
    								inner join cities ct on ct.id = d.city_id
    								inner join regions r on r.id = ct.region_id
    								inner join payments p on o.order_uid = p.order_uid
    								inner join currencies c on c.id = p.currency_id
    								inner join payment_providers pp on p.provider_id = pp.id
    								inner join banks b on p.bank_id = b.id
								where o.order_uid = $1;`)
	row := r.db.QueryRow(query, uid)
	err := row.Scan(
		&(order.UID),
		&(order.TrackNumber),
		&(order.Entry),
		&(order.Locale),
		&(order.InternalSignature),
		&(order.CustomerID),
		&(order.DeliveryService),
		&(order.ShardKey),
		&(order.SmID),
		&(order.DateCreated),
		&(order.OofShard),
		&(order.Delivery.Name),
		&(order.Delivery.Phone),
		&(order.Delivery.Zip),
		&(order.Delivery.City),
		&(order.Delivery.Address),
		&(order.Delivery.Region),
		&(order.Delivery.Email),
		&(order.Payment.Transaction),
		&(order.Payment.RequestID),
		&(order.Payment.Currency),
		&(order.Payment.Provider),
		&(order.Payment.Amount),
		&(order.Payment.PaymentDatetime),
		&(order.Payment.Bank),
		&(order.Payment.DeliveryCost),
		&(order.Payment.GoodsTotal),
		&(order.Payment.CustomFee),
	)

	if err != nil {
		return nil, false
	}

	items := r.getOrderItems(uid)

	order.Items = items

	return &order, true
}

func (r *OrderPostgres) getOrderItems(uid string) []model.Item {

	query := fmt.Sprintf(`select i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, b.brand, i.status
								from items i
								inner join brands b on b.id = i.brand_id
								inner join orders_items oi on i.rid = oi.item_rid
								where oi.order_uid=$1;`)
	rows, err := r.db.Query(query, uid)
	if err != nil {
		return []model.Item{}
	}
	if err != nil {
		return []model.Item{}
	}
	items := make([]model.Item, 0)
	for rows.Next() {
		item := model.Item{}
		err = rows.Scan(
			&(item.ChrtID),
			&(item.TrackNumber),
			&(item.Price),
			&(item.RID),
			&(item.Name),
			&(item.Sale),
			&(item.Size),
			&(item.TotalPrice),
			&(item.NmID),
			&(item.Brand),
			&(item.Status),
		)
		if err != nil {
			return []model.Item{}
		}
		items = append(items, item)
	}
	return items
}

func (r *OrderPostgres) GetAll() []model.Order {
	var orders []model.Order

	query := fmt.Sprintf(`select o.order_uid, o.track_number, o.entry, l.locale, o.internal_signature, o.customer_id, ds.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
       								d.name d_name, d.phone d_phone, d.zip d_zip, ct.city d_city, d.address d_address, r.region d_region, d.email d_email,
       								p.transaction p_transaction, p.request_id p_request_id, c.currency p_currency, pp.provider p_provider, p.amount p_amount, p.payment_dt p_payment_dt, b.bank p_bank, p.delivery_cost p_delivery_cost, p.goods_total p_goods_total, p.custom_fee p_custom_fee
								from orders o
    								inner join delivery_services ds on ds.id = o.delivery_service_id
    								inner join locales l on l.id = o.locale_id
    								inner join deliveries d on o.order_uid = d.order_uid
    								inner join cities ct on ct.id = d.city_id
    								inner join regions r on r.id = ct.region_id
    								inner join payments p on o.order_uid = p.order_uid
    								inner join currencies c on c.id = p.currency_id
    								inner join payment_providers pp on p.provider_id = pp.id
    								inner join banks b on p.bank_id = b.id;`)
	rows, err := r.db.Query(query)

	for rows.Next() {
		var order model.Order

		err = rows.Scan(
			&(order.UID),
			&(order.TrackNumber),
			&(order.Entry),
			&(order.Locale),
			&(order.InternalSignature),
			&(order.CustomerID),
			&(order.DeliveryService),
			&(order.ShardKey),
			&(order.SmID),
			&(order.DateCreated),
			&(order.OofShard),
			&(order.Delivery.Name),
			&(order.Delivery.Phone),
			&(order.Delivery.Zip),
			&(order.Delivery.City),
			&(order.Delivery.Address),
			&(order.Delivery.Region),
			&(order.Delivery.Email),
			&(order.Payment.Transaction),
			&(order.Payment.RequestID),
			&(order.Payment.Currency),
			&(order.Payment.Provider),
			&(order.Payment.Amount),
			&(order.Payment.PaymentDatetime),
			&(order.Payment.Bank),
			&(order.Payment.DeliveryCost),
			&(order.Payment.GoodsTotal),
			&(order.Payment.CustomFee),
		)

		if err != nil {
			return []model.Order{}
		}

		items := r.getOrderItems(order.UID)

		order.Items = items

		orders = append(orders, order)
	}

	return orders
}

func (r *OrderPostgres) getDeliveryServiceIDByName(name string) (int, error) {
	query := fmt.Sprintf("select id from %s where delivery_service = $1", deliveryServiceTable)
	row := r.db.QueryRow(query, name)
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
	query := fmt.Sprintf("select id from %s where locale = $1", localesTable)
	row := r.db.QueryRow(query, name)
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
	query := fmt.Sprintf("select id from %s where bank = $1", banksTable)
	row := r.db.QueryRow(query, name)
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
	query := fmt.Sprintf("select id from %s where region = $1", regionsTable)
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
	query := fmt.Sprintf("select id from %s where provider = $1", paymentProvidersTable)
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
	query := fmt.Sprintf("select id from %s where brand = $1", brandsTable)
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
	query := fmt.Sprintf("select id from %s where currency = $1", currenciesTable)
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

func (r *OrderPostgres) itemExists(rid string) bool {
	query := fmt.Sprintf("select rid from %s where rid=$1", itemsTable)
	err := r.db.QueryRow(query, rid).Scan(&rid)

	if err != nil {
		return false
	}

	return true
}

// Adds item to order
func (r *OrderPostgres) createOrderItem(tx *sql.Tx, item model.Item, orderUID string) (uid, rid string, error error) {
	itemExists := r.itemExists(item.RID)
	if itemExists == false {
		query := fmt.Sprintf("insert into %s (rid, chrt_id, track_number, price, name, sale, size, total_price, nm_id, brand_id, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", itemsTable)

		brandID, err := r.getBrandIDByName(item.Brand)

		if err != nil {
			brandID, err = r.createBrand(tx, item.Brand)
			if err != nil {
				tx.Rollback()
				return "", "", err
			}
		}

		_, err = tx.Exec(
			query,
			item.RID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			brandID,
			item.Status,
		)

		if err != nil {
			tx.Rollback()
			return "", "", err
		}

	}

	query := fmt.Sprintf("insert into %s (order_uid, item_rid) values ($1, $2)", ordersItemsTable)
	_, err := tx.Exec(query, orderUID, item.RID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}

	return orderUID, rid, nil
}
