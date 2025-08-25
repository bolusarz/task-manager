-- name: GetUserById :one
select * from  users where id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: CreateUser :one
insert into users (first_name, last_name, email, password_hash, profile_picture_url) values ($1, $2, $3, $4, $5) returning *;

-- name: UpdateUser :one
update users
	set first_name = $2,
		last_name = $3,
		email = $4,
		profile_picture_url = $5,
        is_email_verified = $6
where id = $1
returning *;

-- name: UpdateUserPassword :exec
update users 
	set password_hash = $2
where id = $1;

-- name: DeleteUser :exec
delete from users where id = $1;