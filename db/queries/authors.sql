-- name: GetAuthor :one
select * from authors
where id = $1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY id;
