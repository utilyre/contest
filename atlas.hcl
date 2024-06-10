env "local" {
  src = "file://db/schemata"
  url = "postgres://admin:admin@localhost:5432/contest?sslmode=disable"
  dev = "docker://postgres/16.1-alpine3.19/dev"
}
