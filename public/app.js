function mtixApp() {
    return {
        view: 'login',
        user: null,
        error: '',
        isLoading: false,
        ...window.mtixAuth(),
        ...window.mtixBooking(),
        ...window.mtixAdmin(),

        init() {
            this.restoreSession();
        },

        async goHome() {
            if (this.user?.is_admin) {
                this.view = 'admin';
                return this.loadAdmin();
            }
            return this.loadMovies();
        },

        formatDate(dateStr) {
            return new Date(dateStr).toLocaleString('id-ID', {
                weekday: 'short',
                day: 'numeric',
                month: 'short',
                hour: '2-digit',
                minute: '2-digit'
            });
        },

        formatRupiah(number) {
            return new Intl.NumberFormat('id-ID', {
                style: 'currency',
                currency: 'IDR',
                maximumFractionDigits: 0
            }).format(number || 0);
        }
    };
}
