window.mtixAuth = function mtixAuth() {
    return {
        bookings: [],

        async restoreSession() {
            if (!localStorage.getItem('mtix_token')) return;
            try {
                const data = await apiFetch('/users/me');
                this.user = data.data;
                await this.goHome();
            } catch {
                localStorage.removeItem('mtix_token');
                this.view = 'login';
            }
        },

        async login(username, password) {
            this.isLoading = true;
            this.error = '';
            try {
                const data = await apiFetch('/login', {
                    method: 'POST',
                    body: JSON.stringify({ username, password })
                });
                localStorage.setItem('mtix_token', data.token);
                this.user = data.user;
                await this.goHome();
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        async register(username, password) {
            this.isLoading = true;
            this.error = '';
            try {
                await apiFetch('/register', {
                    method: 'POST',
                    body: JSON.stringify({ username, password })
                });
                await this.login(username, password);
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        async logout() {
            try {
                await apiFetch('/logout', { method: 'POST' });
            } catch {
                // A missing server session should not trap the user in the UI.
            }
            localStorage.removeItem('mtix_token');
            this.user = null;
            this.error = '';
            this.resetBooking();
            this.view = 'login';
        },

        async showProfile() {
            this.error = '';
            try {
                const [profile, bookings] = await Promise.all([
                    apiFetch('/users/me'),
                    apiFetch('/users/me/bookings')
                ]);
                this.user = profile.data;
                this.bookings = bookings.data;
                this.view = 'profile';
            } catch (error) {
                this.error = error.message;
            }
        },

        async submitStudentApplication(file) {
            if (!file) {
                this.error = 'Choose a JPG, PNG, or PDF file first.';
                return;
            }
            this.isLoading = true;
            this.error = '';
            try {
                const body = new FormData();
                body.append('evidence', file);
                const data = await apiFetch('/users/me/student-application', { method: 'POST', body });
                this.user = data.data;
            } catch (error) {
                this.error = error.message;
            } finally {
                this.isLoading = false;
            }
        },

        studentStatusLabel(status) {
            return ['Non-student', 'Pending', 'Verified'][status] || 'Unknown';
        }
    };
};
