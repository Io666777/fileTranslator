make reup

migrate -path ./schema -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up

start "" http://localhost:5000

pause
