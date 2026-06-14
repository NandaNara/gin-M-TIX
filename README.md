# Gin M-TIX - Cinema Ticket Booking

Aplikasi booking bioskop sederhana dengan Gin, penyimpanan in-memory, Alpine.js,
dan Tailwind CSS. Proyek ini mendemonstrasikan Factory, Strategy, dan Facade
Pattern tanpa database atau frontend build tool.

- **Factory Pattern** untuk membuat jenis tiket.
- **Strategy Pattern** untuk menentukan strategi harga.
- **Facade Pattern** untuk menyederhanakan alur pemesanan dan pembayaran.

Data disimpan secara in-memory melalui `config.Database`, sehingga aplikasi bisa
langsung dijalankan tanpa instalasi database tambahan.

## Tujuan Program

Program ini dibuat untuk menunjukkan bagaimana design pattern dapat dipakai dalam
studi kasus nyata, yaitu sistem booking tiket bioskop. Aplikasi menyediakan
endpoint untuk:

- Login dan logout sederhana.
- Mengelola data film.
- Membuat dan melihat jadwal tayang.
- Melihat kursi pada jadwal tertentu.
- Membuat booking tiket.
- Melihat detail booking.
- Melakukan pembayaran.

Dengan struktur ini, setiap bagian program punya tanggung jawab yang jelas.
Controller menerima request HTTP, service menjalankan business logic, repository
mengelola data, dan folder `patterns` berisi implementasi design pattern.

## Fitur

- Register dan login dengan Bearer token in-memory serta role user/admin.
- Profil user, riwayat booking, serta pengajuan status student melalui JPG,
  PNG, atau PDF maksimal 5 MB.
- Dashboard admin untuk menambah movie, schedule, studio, dan memproses
  pengajuan student.
- Studio otomatis membuat kursi; dua baris terakhir menjadi VIP.
- Harga tiket ditentukan otomatis dari jadwal, kursi VIP, dan status student.
- Ownership check untuk melihat, membatalkan, dan membayar booking.
- Admin ditolak oleh API booking dan payment.

## Akun Demo

| Username | Password | Role |
| --- | --- | --- |
| `admin` | `admin` | Admin |
| `andi` | `andi` | Verified student |
| `budi` | `budi` | User |
| `cici` | `cici` | User |

## Struktur Program

```text
.
├── main.go
├── config/
│   └── database.go
├── controllers/
│   ├── auth_controller.go
│   ├── booking_controller.go
│   ├── movie_controller.go
│   ├── schedule_controller.go
│   └── studio_controller.go
├── models/
│   ├── booking.go
│   ├── movie.go
│   ├── payment.go
│   ├── schedule.go
│   ├── seat.go
│   ├── studio.go
│   ├── ticket.go
│   └── user.go
├── repositories/
│   ├── booking_repository.go
│   ├── movie_repository.go
│   ├── schedule_repository.go
│   ├── studio_repository.go
│   └── user_repository.go
├── services/
│   ├── booking_service.go
│   ├── payment_service.go
│   └── pricing_service.go
├── patterns/
│   ├── factory/
│   │   └── ticket_factory.go
│   ├── strategy/
│   │   ├── pricing_strategy.go
│   │   ├── weekday_pricing.go
│   │   ├── weekend_pricing.go
│   │   ├── holiday_pricing.go
│   │   └── midnight_pricing.go
│   └── facade/
│       └── booking_facade.go
├── public/
│   ├── app.js
│   ├── index.html
│   ├── js/
│   │   ├── admin.js
│   │   ├── api.js
│   │   ├── auth.js
│   │   └── booking.js
│   └── poster/
│       ├── interstellar.jpeg
│       └── the_dark_knight.jpeg
├── routes/
│   ├── routes.go
│   └── routes_test.go
├── Dockerfile
├── docker-compose.yml
└── gin-M-TIX.postman_collection.json
```

## Alur Program

Saat aplikasi dijalankan, `main.go` membuat database in-memory melalui
`config.NewDatabase()`. Database ini berisi seed awal berupa user, movie, studio,
seat, dan schedule.

Setelah itu, `routes.SetupRouter()` membuat repository, service, facade, dan
controller. Semua endpoint didaftarkan di `routes/routes.go`.

Alur utama booking tiket:

1. User melihat daftar film melalui `GET /movies`.
2. User melihat jadwal tayang melalui `GET /schedules`.
3. User melihat daftar kursi pada jadwal tertentu melalui `GET /schedules/:id/seats`.
4. Pada antarmuka web, user melihat total harga dan melakukan konfirmasi terlebih
   dahulu sebelum booking dibuat.
5. User membuat booking melalui `POST /bookings`.
6. `BookingController` mengirim request ke `BookingFacade`.
7. `BookingFacade` meneruskan proses ke `BookingService`.
8. `BookingService` memvalidasi jadwal, kursi, dan jenis tiket.
9. `PricingService` memilih strategi harga berdasarkan waktu tayang.
10. `TicketFactory` membuat tiket sesuai jenis: `regular`, `vip`,
    `regular_student`, atau `vip_student`.
11. `BookingRepository` menyimpan booking dan tiket ke database in-memory.
12. User membayar melalui `POST /payments`.
13. `PaymentService` memvalidasi nominal pembayaran dan mengubah status booking
    menjadi `paid` jika pembayaran berhasil.

## Penerapan Design Pattern

### 1. Factory Pattern

Lokasi: `patterns/factory/ticket_factory.go`

Factory Pattern digunakan untuk membuat tiket berdasarkan jenis kursi dan status
student. Tersedia empat kombinasi:

- `regular` — harga normal.
- `vip` — harga naik 50%.
- `regular_student` — diskon 20%.
- `vip_student` — harga naik 50%, lalu diskon 20%.

Kode pemanggil tidak perlu tahu detail pembuatan masing-masing tiket. Cukup
memanggil `NewTicketFactory(isVIP, isStudent)`, lalu factory akan
mengembalikan pembuat tiket yang sesuai.

Contoh penerapan:

```go
factory := ticketfactory.NewTicketFactory(bool(seat.IsVIP), isStudent)
ticket := factory.CreateTicket(schedule.ID, seat, baseSeatPrice)
```

### 2. Strategy Pattern

Lokasi: `patterns/strategy`

Strategy Pattern digunakan untuk menentukan harga berdasarkan waktu tayang.
`PricingService` memilih strategi berdasarkan `StartTime` dari schedule dengan
urutan prioritas:

1. `HolidayPricing`: harga naik 50% (tanggal 1 Januari dan 25 Desember).
2. `MidnightPricing`: harga naik 20% (jam 22:00–02:00).
3. `WeekendPricing`: harga naik 25% (Sabtu dan Minggu).
4. `WeekdayPricing`: harga normal.

Dengan pattern ini, aturan harga baru bisa ditambahkan tanpa mengubah alur
booking utama.

Contoh:

```go
if strategy.IsHoliday(schedule.StartTime) {
	return strategy.HolidayPricing{}
}
if strategy.IsMidnight(schedule.StartTime) {
	return strategy.MidnightPricing{}
}
if strategy.IsWeekend(schedule.StartTime) {
	return strategy.WeekendPricing{}
}
return strategy.WeekdayPricing{}
```

### 3. Facade Pattern

Lokasi: `patterns/facade/booking_facade.go`

Facade Pattern digunakan untuk menyederhanakan akses dari controller ke proses
booking dan payment. Controller cukup berhubungan dengan `BookingFacade`, tanpa
perlu tahu detail service apa saja yang dipakai di belakangnya.

Contoh:

```go
booking, err := ctrl.facade.CreateBooking(request)
payment, booking, err := ctrl.facade.Pay(userID, request)
```

## Endpoint

| Method | Endpoint | Akses |
| --- | --- | --- |
| `POST` | `/register` | Publik |
| `POST` | `/login` | Publik |
| `POST` | `/logout` | User |
| `GET` | `/movies` | Publik |
| `POST` | `/movies` | Admin |
| `PUT` | `/movies/:id` | Admin |
| `DELETE` | `/movies/:id` | Admin |
| `GET` | `/schedules` | Publik |
| `POST` | `/schedules` | Admin |
| `PUT` | `/schedules/:id` | Admin |
| `DELETE` | `/schedules/:id` | Admin |
| `GET` | `/schedules/:id/seats` | Publik |
| `GET` | `/studios` | Publik |
| `POST` | `/studios` | Admin |
| `PUT` | `/studios/:id` | Admin |
| `DELETE` | `/studios/:id` | Admin |
| `GET` | `/users/me` | User |
| `GET` | `/users/me/bookings` | User |
| `POST` | `/users/me/student-application` | User |
| `POST` | `/bookings` | User (non-admin) |
| `GET` | `/bookings/:id` | Pemilik booking |
| `DELETE` | `/bookings/:id` | Pemilik booking |
| `POST` | `/payments` | User (non-admin) |
| `GET` | `/admin/student-applications` | Admin |
| `GET` | `/admin/student-applications/:id/evidence` | Admin |
| `POST` | `/admin/student-applications/:id/resolve` | Admin |

## Contoh Request

### Login

```bash
curl -X POST http://localhost:8999/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

### Register

```bash
curl -X POST http://localhost:8999/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"userbaru","password":"rahasia"}'
```

### Logout

```bash
curl -X POST http://localhost:8999/logout \
  -H 'Authorization: Bearer TOKEN'
```

### Melihat Film

```bash
curl http://localhost:8999/movies
```

### Membuat Jadwal

```bash
curl -X POST http://localhost:8999/schedules \
  -H 'Authorization: Bearer TOKEN' \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "studio_id": 1,
    "start_time": "2026-05-26T19:00:00+07:00",
    "base_price": 50000
  }'
```

### Melihat Kursi Jadwal

```bash
curl http://localhost:8999/schedules/1/seats
```

### Membuat Booking

`user_id` dan `ticket_type` tidak dikirim oleh client; server mengisinya
otomatis berdasarkan sesi login dan status kursi.

```bash
curl -X POST http://localhost:8999/bookings \
  -H 'Authorization: Bearer TOKEN' \
  -H "Content-Type: application/json" \
  -d '{"schedule_id":1,"seat_ids":[1,57]}'
```

### Melakukan Pembayaran

```bash
curl -X POST http://localhost:8999/payments \
  -H 'Authorization: Bearer TOKEN' \
  -H "Content-Type: application/json" \
  -d '{
    "booking_id": 1,
    "method": "bank_transfer",
    "amount": 90000
  }'
```

### Pengajuan Student

```bash
curl -X POST http://localhost:8999/users/me/student-application \
  -H 'Authorization: Bearer TOKEN' \
  -F 'evidence=@student-card.png'
```

### Movie Admin

Poster dikirim sebagai form data; wajib saat membuat movie dan opsional saat
edit.

```bash
curl -X POST http://localhost:8999/movies \
  -H 'Authorization: Bearer TOKEN' \
  -F 'title=Dune' -F 'genre=Sci-Fi' -F 'duration_minutes=166' \
  -F 'rating=PG-13' -F 'poster=@dune.jpg'
```

## Menjalankan

### Lokal

```bash
go run .
```

Buka `http://localhost:8999`. Root URL mengarahkan ke `/ui/`.

Jalankan pengecekan compile & testing:

```bash
go test ./...
```

### Docker

**Docker Compose (Disarankan):**

```bash
docker compose up --build
```

**Docker CLI:**

```bash
docker build -t gin-mtix-app .
docker run -d -p 8999:8999 --name mtix-app gin-mtix-app
```

## Antarmuka Web

Proyek ini dilengkapi dengan antarmuka web bertema **"Midnight Premiere"** yang
mewah dan elegan. Frontend dibangun dengan stack ringan tanpa build tools:

- **Alpine.js** untuk reaktivitas dan state management.
- **Tailwind CSS** untuk styling (Glassmorphism, Dark mode).
- **Lucide Icons** untuk ikon minimalis.

Seluruh file frontend berada di dalam direktori `public/`.

Fitur utama antarmuka web:

- Login dan logout langsung dari browser.
- Modal konfirmasi sebelum booking dibuat.
- Halaman checkout menampilkan nominal pembayaran yang harus dibayar.

## Catatan

- Data dan token berada di memori dan hilang saat server restart.
- Bukti student disimpan di `uploads/student-applications/` dan dihapus
  setelah approve/reject.
- Kursi dianggap ter-booking berdasarkan `schedule_id`, sehingga kursi yang
  sama bisa dipakai lagi pada jadwal berbeda.
- Password demo masih plain text karena proyek ini bukan sistem produksi.
- Tujuan utama proyek adalah demonstrasi design pattern pada aplikasi REST API
  sederhana.
