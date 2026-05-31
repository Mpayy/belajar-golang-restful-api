## cara install golang-migrate di $GOPATH/bin
```bash
go install -tags 'database1 database2' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## cara membuat database migration
```bash
migrate create -ext sql -dir db/migration nama_tabel
```

## cara menjalankan database migration
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" up
```

## cara mengembalikan database migration
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" down
```

## url untuk melihat cara migrate database lainya : https://github.com/golang-migrate/migrate/databases#

## cara migrate hanya beberapa step keatas
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" up 1
```

## cara migrate hanya beberapa step kebawah
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" down 1
```

## cara mengatasi dirty state, harus manual

1. drop tabel yang bermasalah
```sql
DROP TABLE IF EXISTS nama_tabel;
```
2. lalu perbaiki/mengubah versi migrationnya ke versi sebelumnya
# melihat versi migration sebelumnya
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" version
```

# mengubah versi migration ke versi sebelumnya
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" force versi_sebelumnya
```

3. perbaiki code pada migration yang bermasalah, lalu jalankan migrationnya kembali
```bash
migrate -path db/migration -database "mysql://root:password@tcp(localhost:3306)/nama_database" up
```