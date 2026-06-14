window.mtixAdmin = function mtixAdmin() {
    return {
        adminTab: 'movie',
        studios: [],
        studentApplications: [],
        adminMessage: '',
        editingMovieId: null,
        editingScheduleId: null,
        editingStudioId: null,

        async loadAdmin() {
            if (!this.user?.is_admin) {
                await this.loadMovies();
                this.error = 'Admin access required.';
                return false;
            }
            this.isLoading = true;
            this.error = '';
            try {
                const [movies, schedules, studios, applications] = await Promise.all([
                    apiFetch('/movies'),
                    apiFetch('/schedules'),
                    apiFetch('/studios'),
                    apiFetch('/admin/student-applications')
                ]);
                this.movies = movies.data;
                this.schedules = schedules.data;
                this.studios = studios.data;
                this.studentApplications = applications.data;
                this.view = 'admin';
                return true;
            } catch (error) {
                this.error = error.message;
                return false;
            } finally {
                this.isLoading = false;
            }
        },

        async createMovie(form) {
            const values = new FormData(form);
            const saved = await this.adminRequest(
                this.editingMovieId ? `/movies/${this.editingMovieId}` : '/movies',
                values,
                this.editingMovieId ? 'PUT' : 'POST'
            );
            if (saved) {
                this.editingMovieId = null;
                form.reset();
            }
        },

        async createSchedule(form) {
            const values = new FormData(form);
            const saved = await this.adminRequest(this.editingScheduleId ? `/schedules/${this.editingScheduleId}` : '/schedules', {
                movie_id: Number(values.get('movie_id')),
                studio_id: Number(values.get('studio_id')),
                start_time: new Date(values.get('start_time')).toISOString(),
                base_price: Number(values.get('base_price'))
            }, this.editingScheduleId ? 'PUT' : 'POST');
            if (saved) {
                this.editingScheduleId = null;
                form.reset();
            }
        },

        async createStudio(form) {
            const values = new FormData(form);
            const saved = await this.adminRequest(this.editingStudioId ? `/studios/${this.editingStudioId}` : '/studios', {
                name: values.get('name'),
                seat_rows: Number(values.get('seat_rows')),
                seat_columns: Number(values.get('seat_columns'))
            }, this.editingStudioId ? 'PUT' : 'POST');
            if (saved) {
                this.editingStudioId = null;
                form.reset();
            }
        },

        async adminRequest(path, body, method = 'POST') {
            this.isLoading = true;
            this.error = '';
            this.adminMessage = '';
            try {
                await apiFetch(path, {
                    method,
                    body: body instanceof FormData ? body : JSON.stringify(body)
                });
                this.adminMessage = 'Saved successfully.';
                await this.loadAdmin();
                return true;
            } catch (error) {
                this.error = error.message;
                return false;
            } finally {
                this.isLoading = false;
            }
        },

        editMovie(movie, form) {
            this.editingMovieId = movie.id;
            form.elements.title.value = movie.title;
            form.elements.genre.value = movie.genre;
            form.elements.duration_minutes.value = movie.duration_minutes;
            form.elements.rating.value = movie.rating || '';
            form.elements.poster.value = '';
            form.scrollIntoView({ behavior: 'smooth' });
        },

        editSchedule(schedule, form) {
            this.editingScheduleId = schedule.id;
            form.elements.movie_id.value = schedule.movie_id;
            form.elements.studio_id.value = schedule.studio_id;
            const date = new Date(schedule.start_time);
            form.elements.start_time.value = new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
            form.elements.base_price.value = schedule.base_price;
            form.scrollIntoView({ behavior: 'smooth' });
        },

        editStudio(studio, form) {
            this.editingStudioId = studio.id;
            form.elements.name.value = studio.name;
            form.elements.seat_rows.value = studio.seat_rows;
            form.elements.seat_columns.value = studio.seat_columns;
            form.scrollIntoView({ behavior: 'smooth' });
        },

        cancelEdit(type, form) {
            this[`editing${type}Id`] = null;
            form.reset();
        },

        async deleteAdmin(path) {
            if (!window.confirm('Delete this item?')) return;
            this.error = '';
            try {
                await apiFetch(path, { method: 'DELETE' });
                this.adminMessage = 'Deleted successfully.';
                this.editingMovieId = null;
                this.editingScheduleId = null;
                this.editingStudioId = null;
                await this.loadAdmin();
            } catch (error) {
                this.error = error.message;
            }
        },

        async resolveStudent(userID, action) {
            this.error = '';
            try {
                await apiFetch(`/admin/student-applications/${userID}/resolve`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ approved: action === 'approve' })
                });
                await this.loadAdmin();
            } catch (error) {
                this.error = error.message;
            }
        },

        async viewStudentEvidence(userID) {
            this.error = '';
            const popup = window.open('', '_blank');
            try {
                const blob = await apiFetch(`/admin/student-applications/${userID}/evidence`);
                const url = URL.createObjectURL(blob);
                if (popup) popup.location = url;
                setTimeout(() => URL.revokeObjectURL(url), 60000);
            } catch (error) {
                if (popup) popup.close();
                this.error = error.message;
            }
        }
    };
};
