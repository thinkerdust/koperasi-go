koperasi-go/
├── api/                   # Tempat controller/handler HTTP
│   └── auth.go
│   └── user.go
├── assets/                # File statis (misalnya template, export file, dll.)
├── db/                    # Inisialisasi DB & migration
│   └── database.go
├── model/                 # Struct untuk tabel DB (User, Anggota, Pinjaman, dsb.)
│   └── user.go
│   └── anggota.go
├── repository/            # Akses query DB (UserRepository, PinjamanRepository, dll.)
│   └── user_repository.go
│   └── pinjaman_repository.go
├── middleware/            # Middleware (auth JWT, logging, error handler)
│   └── auth_middleware.go
├── helpers/               # Fungsi bantu (hash password, response formatter, dll.)
│   └── response.go
│   └── hash.go
├── routes/                # Routing API
│   └── routes.go
├── go.mod
├── go.sum
├── main.go                # Entry point aplikasi
├── README.md
├── .env                   # Konfigurasi environment
└── .gitignore
