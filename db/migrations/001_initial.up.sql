create user app with encrypted password 'app_pwd';

alter default privileges in schema public
  grant select, insert, update, delete
on tables to app;

alter default privileges in schema public
  grant execute on functions to app;

alter default privileges in schema public
  grant select, update
on sequences to app;

grant usage on schema public to app;

create table authors
(
    id         bigserial primary key,
    name       text      not null,
    extra      jsonb null,
    created_at timestamp not null default now()
);
