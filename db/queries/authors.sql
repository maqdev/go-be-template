-- name: GetAuthor :one
select id, name, extra, created_at from authors
where id = $1;

-- name: ListAuthors :many
select id, name, extra, created_at from authors
where (sqlc.narg(prevID)::bigint is null or id > sqlc.narg(prevID)::bigint)
order by id
limit @lim;

-- name: CreateAuthor :one
insert into authors (
  name, extra
) values (
  $1, $2
)
returning id;