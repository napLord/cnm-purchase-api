-- +goose Up
create type lockStatus as enum ('locked', 'unlocked');

create sequence purchases_events_id_seq;

create table purchases_events
(
    id          bigint primary key default nextval('purchases_events_id_seq'),
    purchase_id bigint     not null references purchases (id),
    type        text       not null,
    status      lockStatus not null,
    payload     jsonb,
    updated     timestamp
);

create index on purchases_events using hash(id);

-- +goose Down
drop table purchases_events;

drop sequence purchases_events_id_seq;

drop type lockStatus;
