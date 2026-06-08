# Gin M-TIX - Cinema Ticket Booking

Proyek ini adalah contoh sederhana sistem pemesanan tiket bioskop menggunakan Golang Gin. Fokus utama dari program ini ada pada penerapan tiga design pattern:

- Factory Pattern untuk membuat jenis tiket.
- Strategy Pattern untuk menentukan strategi harga.
- Facade Pattern untuk menyederhanakan alur pemesanan dan pembayaran.

Data disimpan secara in-memory melalui `config.Database`, sehingga aplikasi bisa langsung dijalankan tanpa instalasi database tambahan. Ketika server dimatikan, data booking dan payment yang dibuat saat runtime akan hilang.

## Tujuan Program

Program ini dibuat untuk menunjukkan bagaimana design pattern dapat dipakai dalam studi kasus nyata, yaitu sistem booking tiket bioskop. Aplikasi menyediakan endpoint untuk:

- Login sederhana.
- Mengelola data film.
- Membuat dan melihat jadwal tayang.
- Melihat kursi pada jadwal tertentu.
- Membuat booking tiket.
- Melihat detail booking.
- Melakukan pembayaran.

Dengan struktur ini, setiap bagian program punya tanggung jawab yang jelas. Controller menerima request HTTP, service menjalankan business logic, repository mengelola data, dan folder `patterns` berisi implementasi design pattern.

## Struktur Program

```text
.
├── main.go
├── config/
│   └── database.go
├── controllers/
│   ├── auth_controller.go
│   ├── movie_controller.go
│   ├── schedule_controller.go
│   └── booking_controller.go
├── models/
│   ├── movie.go
│   ├── studio.go
│   ├── seat.go
│   ├── schedule.go
│   ├── booking.go
│   ├── ticket.go
│   └── payment.go
├── repositories/
│   ├── movie_repository.go
│   ├── schedule_repository.go
│   └── booking_repository.go
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
│   └── index.html
└── routes/
    └── routes.go
```

## Alur Program

Saat aplikasi dijalankan, `main.go` membuat database in-memory melalui `config.NewDatabase()`. Database ini berisi seed awal berupa movie, studio, seat, dan schedule.

Setelah itu, `routes.SetupRouter()` membuat repository, service, facade, dan controller. Semua endpoint didaftarkan di file `routes/routes.go`.

Alur utama booking tiket:

1. User melihat daftar film melalui `GET /movies`.
2. User melihat jadwal tayang melalui `GET /schedules`.
3. User melihat daftar kursi pada jadwal tertentu melalui `GET /schedules/:id/seats`.
4. User membuat booking melalui `POST /bookings`.
5. `BookingController` mengirim request ke `BookingFacade`.
6. `BookingFacade` meneruskan proses ke `BookingService`.
7. `BookingService` memvalidasi jadwal, kursi, dan jenis tiket.
8. `PricingService` memilih strategi harga weekday atau weekend.
9. `TicketFactory` membuat tiket sesuai jenis tiket, yaitu `regular` atau `vip`.
10. `BookingRepository` menyimpan booking dan tiket ke database in-memory.
11. User membayar melalui `POST /payments`.
12. `PaymentService` memvalidasi nominal pembayaran dan mengubah status booking menjadi `paid` jika pembayaran berhasil.

## Penerapan Design Pattern

### 1. Factory Pattern

Lokasi: `patterns/factory/ticket_factory.go`

Factory Pattern digunakan untuk membuat tiket berdasarkan jenis tiket. Saat ini tersedia:

- `regular` — harga normal.
- `vip` — harga naik 50%.
- `student` — diskon 20%.

Kode pemanggil tidak perlu tahu detail pembuatan masing-masing tiket. Cukup memanggil `NewTicketFactory(ticketType)`, lalu factory akan mengembalikan pembuat tiket yang sesuai.

Contoh penerapan:

```go
factory, err := ticketfactory.NewTicketFactory(request.TicketType)
ticket := factory.CreateTicket(schedule.ID, seat, baseSeatPrice)
```

### 2. Strategy Pattern

Lokasi: `patterns/strategy`

Strategy Pattern digunakan untuk menentukan harga berdasarkan waktu tayang. `PricingService` memilih strategi berdasarkan `StartTime` dari schedule dengan urutan prioritas:

1. `HolidayPricing`: harga naik 50% (tanggal 1 Januari dan 25 Desember).
2. `MidnightPricing`: harga naik 20% (jam 22:00–02:00).
3. `WeekendPricing`: harga naik 25% (Sabtu dan Minggu).
4. `WeekdayPricing`: harga normal.

Dengan pattern ini, aturan harga baru bisa ditambahkan tanpa mengubah alur booking utama.

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

Facade Pattern digunakan untuk menyederhanakan akses dari controller ke proses booking dan payment. Controller cukup berhubungan dengan `BookingFacade`, tanpa perlu tahu detail service apa saja yang dipakai di belakangnya.

Contoh:

```go
booking, err := ctrl.facade.CreateBooking(request)
payment, booking, err := ctrl.facade.Pay(request)
```

## Endpoint

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | `/login` | Login demo |
| GET | `/movies` | Melihat daftar film |
| POST | `/movies` | Menambah film |
| PUT | `/movies/:id` | Mengubah film |
| DELETE | `/movies/:id` | Menghapus film |
| GET | `/schedules` | Melihat jadwal tayang |
| POST | `/schedules` | Menambah jadwal tayang |
| GET | `/schedules/:id/seats` | Melihat kursi pada jadwal tertentu |
| POST | `/bookings` | Membuat booking |
| GET | `/bookings/:id` | Melihat detail booking |
| GET | `/users/:id/bookings` | Melihat booking milik user |
| POST | `/payments` | Melakukan pembayaran |

## Antarmuka Web (Frontend)

Proyek ini dilengkapi dengan antarmuka web bertema **"Midnight Premiere"** yang mewah dan elegan. Frontend dibangun dengan stack super ringan tanpa *build tools*:
- **Alpine.js** untuk reaktivitas dan state management.
- **Tailwind CSS** untuk *styling* antarmuka (Glassmorphism, Dark mode).
- **Lucide Icons** untuk ikon minimalis.

Seluruh file frontend berada di dalam direktori `public/`.

## Cara Menjalankan

Jalankan aplikasi:

```bash
go run .
```

Server akan berjalan pada berikut namun juga bisa diakses melalui Antarmuka Web (UI) melalui browser di:

```text
http://localhost:8080
```
*(Root URL otomatis akan mengarahkan Anda ke `/ui/`)*

Jalankan pengecekan compile & testing:

```bash
go test ./...
```

## Contoh Request

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

### Melihat Film

```bash
curl http://localhost:8080/movies
```

### Membuat Jadwal

```bash
curl -X POST http://localhost:8080/schedules \
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
curl http://localhost:8080/schedules/1/seats
```

### Membuat Booking

```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "schedule_id": 1,
    "seat_ids": [1, 2],
    "ticket_type": "regular"
  }'
```

Untuk tiket VIP, gunakan:

```json
{
  "ticket_type": "vip"
}
```

### Melakukan Pembayaran

```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{
    "booking_id": 1,
    "method": "bank_transfer",
    "amount": 90000
  }'
```

## Catatan

- Aplikasi ini menggunakan data in-memory, bukan database permanen.
- Login hanya demo dan belum menggunakan JWT asli.
- Kursi dianggap ter-booking berdasarkan `schedule_id`, sehingga kursi yang sama bisa dipakai lagi pada jadwal berbeda.
- Tujuan utama proyek adalah demonstrasi design pattern pada aplikasi REST API sederhana.
