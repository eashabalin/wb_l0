create table if not exists locales (
    id serial primary key,
    locale varchar not null unique
);

create table if not exists delivery_services (
    id serial primary key,
    delivery_service varchar not null unique
);

create table if not exists regions (
    id serial primary key,
    region varchar not null unique
);

create table if not exists cities (
    id serial primary key,
    city varchar not null unique,
    region_id serial not null constraint fk_city_region references regions
);

create table if not exists banks (
    id serial primary key,
    bank varchar not null unique
);

create table if not exists currencies (
    id serial primary key,
    currency varchar not null unique
);

create table if not exists payment_providers (
    id serial primary key,
    provider varchar not null unique
);

create table if not exists orders (
    order_uid varchar primary key,
    track_number varchar,
    entry varchar,
    locale_id int references locales,
    internal_signature varchar,
    customer_id varchar,
    delivery_service_id int references delivery_services,
    shardkey varchar,
    sm_id int,
    date_created date,
    oof_shard varchar
);

create table if not exists payments (
    order_uid varchar primary key,
    transaction varchar not null,
    request_id varchar not null,
    currency_id int not null references currencies,
    provider_id int not null references payment_providers,
    amount int not null,
    payment_dt int not null,
    bank_id int not null references banks,
    delivery_cost int not null,
    goods_total int not null,
    custom_fee int not null
);

create table if not exists deliveries (
    order_uid varchar primary key references orders,
    name varchar not null,
    phone varchar(20) not null,
    zip varchar not null,
    city_id int not null references cities,
    address varchar not null,
    email varchar not null
);

create table if not exists brands (
    id serial primary key,
    brand varchar not null unique
);

create table if not exists items (
    chrt_id int not null,
    track_number varchar not null,
    price int not null,
    rid varchar primary key,
    name varchar not null,
    sale int not null,
    size varchar not null,
    total_price int not null,
    nm_id int not null,
    brand_id int not null references brands,
    status int not null
);

create table if not exists orders_items (
    order_uid varchar not null references orders,
    item_rid varchar not null references items,
    primary key (order_uid, item_rid)
);
