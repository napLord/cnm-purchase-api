-- +goose Up
create sequence purchases_id_seq;

create table purchases
(
    id       bigint primary key default nextval('purchases_id_seq'),
    total_sum bigint    not null,
    removed  bool      not null,
    created  timestamp not null,
    updated  timestamp
);

create index on purchases using hash(id);

-- +goose Down
drop table purchases;

drop sequence  purchases_id_seq;
