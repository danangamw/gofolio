data "external_schema" "gorm" {
  # Atlas calls the loader that outputs DDL SQL from GORM structs to stdout.
  program = [
    "go", "run", "-mod=mod",
    "./cmd/atlas-loader/main.go",
  ]
}

env "local" {
  src = data.external_schema.gorm.url

  # Dev database: Docker Postgres ephemeral container yang dipakai Atlas
  # untuk menghitung diff antara skema saat ini vs skema baru
  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir    = "file://migrations"
    format = atlas
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")

  migration {
    dir    = "file://migrations"
    format = atlas
  }
}
