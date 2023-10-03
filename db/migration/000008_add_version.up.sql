alter table users add column if not exists version_id int not null default 0;
alter table projects add column if not exists version_id int not null default 0;
alter table tasks add column if not exists version_id int not null default 0;