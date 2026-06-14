window.mtixBooking = function mtixBooking() {
    return {
        movies: [],
        schedules: [],
        seats: [],
        selectedMovie: null,
        selectedSchedule: null,
        selectedSeats: [],
        currentBooking: null,
        paymentAmount: 0,
        showConfirmModal: false,

        resetBooking() {
            this.movies = [];
            this.schedules = [];
            this.seats = [];
            this.selectedMovie = null;
            this.selectedSchedule = null;
            this.selectedSeats = [];
            this.currentBooking = null;
            this.paymentAmount = 0;
            this.showConfirmModal = false;
        },

        async loadMovies() {
            this.isLoading = true;
            this.error = '';
            try {
                const data = await apiFetch('/movies');
                this.movies = data.data;
                this.view = 'movies';
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        async selectMovie(movie) {
            this.selectedMovie = movie;
            this.isLoading = true;
            try {
                const data = await apiFetch('/schedules');
                this.schedules = data.data.filter(schedule => schedule.movie_id === movie.id);
                this.view = 'schedules';
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        async selectSchedule(schedule) {
            this.selectedSchedule = schedule;
            this.selectedSeats = [];
            this.isLoading = true;
            try {
                const data = await apiFetch(`/schedules/${schedule.id}/seats`);
                this.seats = data.data;
                this.view = 'seats';
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        toggleSeat(seat) {
            if (seat.status !== 'available') return;
            const index = this.selectedSeats.findIndex(selected => selected.id === seat.id);
            if (index >= 0) this.selectedSeats.splice(index, 1);
            else this.selectedSeats.push(seat);
        },

        totalPrice() {
            if (!this.selectedSchedule) return 0;
            const studentRate = this.user?.is_student === 2 ? 0.8 : 1;
            return this.selectedSeats.reduce((total, seat) => {
                return total + this.selectedSchedule.seat_price * (seat.is_vip ? 1.5 : 1) * studentRate;
            }, 0);
        },

        async bookTickets() {
            if (!this.selectedSeats.length) return;
            this.isLoading = true;
            this.error = '';
            try {
                const data = await apiFetch('/bookings', {
                    method: 'POST',
                    body: JSON.stringify({
                        schedule_id: this.selectedSchedule.id,
                        seat_ids: this.selectedSeats.map(seat => seat.id)
                    })
                });
                this.currentBooking = data.data;
                this.paymentAmount = data.data.total_price;
                this.showConfirmModal = false;
                this.view = 'checkout';
            } catch (error) {
                this.error = error.message;
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
                await apiFetch(`/bookings/${this.currentBooking.id}`, { method: 'DELETE' });
                this.currentBooking = null;
                await this.selectSchedule(this.selectedSchedule);
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        async pay() {
            this.isLoading = true;
            this.error = '';
            try {
                await apiFetch('/payments', {
                    method: 'POST',
                    body: JSON.stringify({
                        booking_id: this.currentBooking.id,
                        method: 'credit_card',
                        amount: Number(this.paymentAmount)
                    })
                });
                this.view = 'success';
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        }
    };
};
