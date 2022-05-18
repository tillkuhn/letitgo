-- thanks https://antonz.org/json-virtual-columns/
create table events (value TEXT);
alter table events add column object_id integer as (json_extract(value, '$.object_id'));
alter table events add column object text as (json_extract(value, '$.object'));
alter table events add column action text as (json_extract(value, '$.action'));
create index events_object_id on events(object_id);
insert into events values ('{"timestamp":"2022-05-15T09:31:00Z","object":"user","object_id":11,"action":"login","details":{"ip":"192.168.0.1"}}');
insert into events values ('{"timestamp":"2022-05-15T09:32:00Z","object":"account","object_id":12,"action":"deposit","details":{"amount":"1000","currency":"USD"}}');
insert into events values ('{"timestamp":"2022-05-15T09:33:00Z","object":"company","object_id":13,"action":"edit","details":{"fields":["address","phone"]}}');
select object, action from events where object_id = 11;
select object, action from events;

