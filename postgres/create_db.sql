drop table if exists deliveries, payments, items, order_info, order_delivery, invalid_data;

create table deliveries(
    id          serial primary key,
    name        varchar(30),
    phone       varchar(20),
    zip         varchar(20),
    city        varchar(30),
    address     varchar(30),
    region      varchar(30),
    email       varchar(30)
);

create table payments(
    transaction   varchar(50) primary key,
    request_id    varchar(50),
    currency      varchar(10),
    provider      varchar(30),
    amount        integer,
    payment_dt    bigint,
    bank          varchar(30),
    delivery_cost integer,
    goods_totals  integer,
    custom_fee    integer
);

create table order_info(
    order_uid           varchar(50) primary key references payments(transaction),
    track_number        varchar(50) unique,
    entry               varchar(15),
    locale              varchar(15),
    internal_signature  varchar(15),
    customer_id         varchar(15),
    delivery_service    varchar(15),
    shardkey            varchar(15),
    sm_id               integer,
    date_created        timestamp,
    oof_shard           varchar(15)
);

create table items(
    chrt_id       bigint primary key,
    track_number  varchar(50) references order_info(track_number),
    price         integer,
    rid           varchar(50),
    name          varchar(30),
    sale          integer,
    size          varchar(30),
    total_price   integer,
    nm_id         bigint,
    brand         varchar(30),
    status        integer
);

create table order_delivery(
    order_uid           varchar(50) references order_info(order_uid),
    delivery_id         integer references deliveries(id)
);

create table invalid_data
(
    id                  serial primary key,
    data                varchar,
    timestamp           timestamp default now()
);