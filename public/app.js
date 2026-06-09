function mtixApp() {
    return {
        view: 'login', // login, movies, schedules, seats, checkout, success
        user: null,
        
        // Data
        movies: [],
        schedules: [],
        seats: [],
        
        // Selection
        selectedMovie: null,
        selectedSchedule: null,
        selectedSeats: [],
        ticketType: 'regular',
        
        // Booking
        currentBooking: null,
        paymentAmount: 0,
        showConfirmModal: false,
        
        // Error handling
        error: '',
        isLoading: false,
        
        // Images (Placeholders for demo)
        posters: {
            1: 'https://images.unsplash.com/photo-1534447677768-be436bb09401?q=80&w=1000&auto=format&fit=crop', // Space
            2: 'https://images.unsplash.com/photo-1478479405421-ce83c92fb3ba?q=80&w=1000&auto=format&fit=crop'  // City/Batman vibe
        },

        init() {
            // Check if logged in (mock)
            const savedUser = localStorage.getItem('mtix_user');
            if (savedUser) {
                this.user = JSON.parse(savedUser);
                this.loadMovies();
            }
        },

        async login(username, password) {
            this.isLoading = true;
            this.error = '';
            try {
                const res = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });
                
                if (!res.ok) throw new Error('Login failed. Use admin/admin');
                
                const data = await res.json();
                this.user = { id: 1, username: username, token: data.token };
                localStorage.setItem('mtix_user', JSON.stringify(this.user));
                this.loadMovies();
            } catch (err) {
                this.error = err.message;
            } finally {
                this.isLoading = false;
            }
        },

        async logout() {
            if (this.user) {
                try {
                    await fetch('/logout', { method: 'POST' });
                } catch (e) {
                    console.error('Logout failed', e);
                }
            }
            
            // Reset all state
            this.user = null;
            this.movies = [];
            this.schedules = [];
            this.seats = [];
            this.selectedMovie = null;
            this.selectedSchedule = null;
            this.selectedSeats = [];
            this.ticketType = 'regular';
            this.currentBooking = null;
            this.paymentAmount = 0;
            this.showConfirmModal = false;
            this.error = '';
            
            localStorage.removeItem('mtix_user');
            this.view = 'login';
        },

        async loadMovies() {
            this.isLoading = true;
            try {
                const res = await fetch('/movies');
                const data = await res.json();
                this.movies = data.data;
                this.view = 'movies';
            } catch (err) {
                console.error(err);
            } finally {
                this.isLoading = false;
            }
        },

        async selectMovie(movie) {
            this.selectedMovie = movie;
            this.isLoading = true;
            try {
                const res = await fetch('/schedules');
                const data = await res.json();
                // Filter schedules for this movie
                this.schedules = data.data.filter(s => s.movie_id === movie.id);
                this.view = 'schedules';
            } catch (err) {
                console.error(err);
            } finally {
                this.isLoading = false;
            }
        },

        async selectSchedule(schedule) {
            this.selectedSchedule = schedule;
            this.selectedSeats = [];
            this.isLoading = true;
            try {
                const res = await fetch(`/schedules/${schedule.id}/seats`);
                const data = await res.json();
                this.seats = data.data;
                this.view = 'seats';
            } catch (err) {
                console.error(err);
            } finally {
                this.isLoading = false;
            }
        },

        toggleSeat(seat) {
            if (seat.status !== 'available') return;
            
            const index = this.selectedSeats.findIndex(s => s.id === seat.id);
            if (index > -1) {
                this.selectedSeats.splice(index, 1);
            } else {
                this.selectedSeats.push(seat);
            }
        },

        get totalPrice() {
            if (!this.selectedSchedule) return 0;
            let base = this.selectedSchedule.seat_price * this.selectedSeats.length;
            if (this.ticketType === 'vip') return base * 1.5;
            if (this.ticketType === 'student') return base * 0.8;
            return base;
        },

        async bookTickets() {
            if (this.selectedSeats.length === 0) return;
            
            this.isLoading = true;
            this.error = '';
            
            try {
                const req = {
                    user_id: this.user.id,
                    schedule_id: this.selectedSchedule.id,
                    seat_ids: this.selectedSeats.map(s => s.id),
                    ticket_type: this.ticketType
                };
                
                const res = await fetch('/bookings', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(req)
                });
                
                if (!res.ok) throw new Error('Booking failed');
                
                const data = await res.json();
                this.currentBooking = data.data;
                this.paymentAmount = this.currentBooking.total_price;
                this.view = 'checkout';
                this.showConfirmModal = false;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.isLoading = false;
            }
        },

        async cancelBooking() {
            if (!this.currentBooking) {
                this.view = 'seats';
                return;
            }
            
            this.isLoading = true;
            try {
                await fetch(`/bookings/${this.currentBooking.id}`, {
                    method: 'DELETE'
                });
                this.currentBooking = null;
                this.selectedSeats = [];
                // Refresh seats
                this.selectSchedule(this.selectedSchedule);
            } catch (err) {
                console.error(err);
                this.view = 'seats';
            } finally {
                this.isLoading = false;
            }
        },

        async pay() {
            this.isLoading = true;
            this.error = '';
            
            try {
                const res = await fetch('/payments', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        booking_id: this.currentBooking.id,
                        method: 'credit_card',
                        amount: parseFloat(this.paymentAmount)
                    })
                });
                
                if (!res.ok) throw new Error('Payment failed. Check amount.');
                
                const data = await res.json();
                this.view = 'success';
            } catch (err) {
                this.error = err.message;
            } finally {
                this.isLoading = false;
            }
        },
        
        formatDate(dateStr) {
            const d = new Date(dateStr);
            return d.toLocaleString('id-ID', { weekday: 'short', day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
        },
        
        formatRupiah(number) {
            return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(number);
        }
    }
}
