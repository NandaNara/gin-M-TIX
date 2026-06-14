package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gin-M-TIX/config"
	"gin-M-TIX/models"

	"github.com/gin-gonic/gin"
)

func TestStudentApplicationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = os.RemoveAll("uploads")
	t.Cleanup(func() { _ = os.RemoveAll("uploads") })
	db := config.NewDatabase()
	router := SetupRouter(db)
	userToken := login(t, router, "budi", "budi")
	adminToken := login(t, router, "admin", "admin")

	response := request(t, router, http.MethodGet, "/admin/student-applications", nil, userToken, "")
	if response.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d, want 403", response.Code)
	}

	response = uploadEvidence(t, router, userToken, "student.txt", "not an image")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("invalid upload status = %d, want 400", response.Code)
	}

	response = uploadEvidence(t, router, userToken, "student.png", "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	if response.Code != http.StatusCreated {
		t.Fatalf("valid upload status = %d, body = %s", response.Code, response.Body.String())
	}
	evidencePath := db.StudentApplications[3].EvidencePath
	if _, err := os.Stat(evidencePath); err != nil {
		t.Fatalf("evidence was not stored on disk: %v", err)
	}

	response = request(t, router, http.MethodGet, "/users/me", nil, userToken, "")
	if strings.Contains(response.Body.String(), "password") {
		t.Fatal("profile response exposed password")
	}
	var profile struct {
		Data models.User `json:"data"`
	}
	decode(t, response, &profile)
	if profile.Data.IsStudent != models.StudentPending {
		t.Fatalf("student status = %d, want pending", profile.Data.IsStudent)
	}

	response = request(t, router, http.MethodGet, "/admin/student-applications", nil, adminToken, "")
	var applications struct {
		Data []models.StudentApplication `json:"data"`
	}
	decode(t, response, &applications)
	if len(applications.Data) != 1 || applications.Data[0].UserID != profile.Data.ID {
		t.Fatalf("pending applications = %#v", applications.Data)
	}

	response = request(t, router, http.MethodGet, fmt.Sprintf("/admin/student-applications/%d/evidence", profile.Data.ID), nil, adminToken, "")
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("evidence status/content-type = %d/%q", response.Code, response.Header().Get("Content-Type"))
	}

	response = request(t, router, http.MethodPost, fmt.Sprintf("/admin/student-applications/%d/resolve", profile.Data.ID), map[string]bool{
		"approved": true,
	}, adminToken, "application/json")
	if response.Code != http.StatusOK {
		t.Fatalf("approve status = %d, body = %s", response.Code, response.Body.String())
	}
	response = request(t, router, http.MethodGet, "/users/me", nil, userToken, "")
	decode(t, response, &profile)
	if profile.Data.IsStudent != models.StudentTrue {
		t.Fatalf("student status = %d, want verified", profile.Data.IsStudent)
	}
	if _, err := os.Stat(evidencePath); !os.IsNotExist(err) {
		t.Fatalf("approved evidence still exists: %v", err)
	}

	rejectedToken := login(t, router, "cici", "cici")
	response = uploadEvidence(t, router, rejectedToken, "student.pdf", "%PDF-1.4\n")
	if response.Code != http.StatusCreated {
		t.Fatalf("PDF upload status = %d, body = %s", response.Code, response.Body.String())
	}
	response = request(t, router, http.MethodPost, "/admin/student-applications/4/resolve", map[string]bool{
		"approved": false,
	}, adminToken, "application/json")
	if response.Code != http.StatusOK {
		t.Fatalf("reject status = %d", response.Code)
	}
	response = request(t, router, http.MethodGet, "/users/me", nil, rejectedToken, "")
	decode(t, response, &profile)
	if profile.Data.IsStudent != models.StudentFalse {
		t.Fatalf("rejected student status = %d, want non-student", profile.Data.IsStudent)
	}
	response = request(t, router, http.MethodGet, "/admin/student-applications/4/evidence", nil, adminToken, "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("resolved evidence status = %d, want 404", response.Code)
	}
	response = uploadEvidence(t, router, rejectedToken, "student.pdf", "%PDF-1.4\n")
	if response.Code != http.StatusCreated {
		t.Fatalf("reapply status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestRegisterAndAdminTransactionGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter(config.NewDatabase())

	response := request(t, router, http.MethodPost, "/register", map[string]string{
		"username": "newuser",
		"password": "secret",
	}, "", "application/json")
	if response.Code != http.StatusCreated || strings.Contains(response.Body.String(), "password") {
		t.Fatalf("register status/body = %d/%s", response.Code, response.Body.String())
	}
	response = request(t, router, http.MethodPost, "/register", map[string]string{
		"username": "NEWUSER",
		"password": "secret",
	}, "", "application/json")
	if response.Code != http.StatusConflict {
		t.Fatalf("duplicate register status = %d, want 409", response.Code)
	}
	_ = login(t, router, "newuser", "secret")

	adminToken := login(t, router, "admin", "admin")
	response = request(t, router, http.MethodPost, "/bookings", map[string]any{
		"schedule_id": 1,
		"seat_ids":    []int{1},
	}, adminToken, "application/json")
	if response.Code != http.StatusForbidden {
		t.Fatalf("admin booking status = %d, want 403", response.Code)
	}
	response = request(t, router, http.MethodPost, "/payments", map[string]any{
		"booking_id": 1,
		"method":     "credit_card",
		"amount":     50000,
	}, adminToken, "application/json")
	if response.Code != http.StatusForbidden {
		t.Fatalf("admin payment status = %d, want 403", response.Code)
	}
}

func TestStudioCreationGeneratesVIPRows(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := config.NewDatabase()
	router := SetupRouter(db)
	adminToken := login(t, router, "admin", "admin")

	response := request(t, router, http.MethodPost, "/studios", map[string]any{
		"name": "Studio Test", "seat_rows": 3, "seat_columns": 3,
	}, adminToken, "application/json")
	if response.Code != http.StatusCreated {
		t.Fatalf("create studio status = %d, body = %s", response.Code, response.Body.String())
	}
	var created struct {
		Data models.Studio `json:"data"`
	}
	decode(t, response, &created)

	count, vipCount := 0, 0
	for _, seat := range db.Seats {
		if seat.StudioID != created.Data.ID {
			continue
		}
		count++
		if bool(seat.IsVIP) {
			vipCount++
			if seat.Row == "A" {
				t.Fatalf("VIP seat found outside last two rows: %#v", seat)
			}
		}
	}
	if count != 9 || vipCount != 6 {
		t.Fatalf("seat counts = %d total, %d VIP; want 9 and 6", count, vipCount)
	}
}

func TestAdminCRUDAndUserAccessGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter(config.NewDatabase())
	userToken := login(t, router, "budi", "budi")
	adminToken := login(t, router, "admin", "admin")

	for _, endpoint := range []struct {
		method string
		path   string
	}{
		{http.MethodPut, "/movies/1"},
		{http.MethodDelete, "/movies/1"},
		{http.MethodPut, "/schedules/1"},
		{http.MethodDelete, "/schedules/1"},
		{http.MethodPut, "/studios/1"},
		{http.MethodDelete, "/studios/1"},
	} {
		response := request(t, router, endpoint.method, endpoint.path, map[string]any{}, userToken, "application/json")
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s %s status = %d, want 403", endpoint.method, endpoint.path, response.Code)
		}
	}

	response := uploadMoviePoster(t, router, http.MethodPost, "/movies", adminToken, map[string]string{
		"title": "Editable", "genre": "Drama", "duration_minutes": "100",
	}, "poster.png", "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	var movie struct {
		Data models.Movie `json:"data"`
	}
	decode(t, response, &movie)
	oldPoster := filepath.Join("public", strings.TrimPrefix(movie.Data.PosterURL, "/ui/"))
	if _, err := os.Stat(oldPoster); err != nil {
		t.Fatalf("poster was not stored: %v", err)
	}
	response = uploadMoviePoster(t, router, http.MethodPut, fmt.Sprintf("/movies/%d", movie.Data.ID), adminToken, map[string]string{
		"title": "Edited", "genre": "Drama", "duration_minutes": "110",
	}, "replacement.jpg", "\xff\xd8\xff\xe0\x00\x10JFIF\x00")
	decode(t, response, &movie)
	if movie.Data.Title != "Edited" {
		t.Fatalf("updated movie title = %q", movie.Data.Title)
	}
	if _, err := os.Stat(oldPoster); !os.IsNotExist(err) {
		t.Fatalf("old poster still exists: %v", err)
	}
	newPoster := filepath.Join("public", strings.TrimPrefix(movie.Data.PosterURL, "/ui/"))

	response = request(t, router, http.MethodPost, "/studios", map[string]any{
		"name": "Editable Studio", "seat_rows": 2, "seat_columns": 3,
	}, adminToken, "application/json")
	var studio struct {
		Data models.Studio `json:"data"`
	}
	decode(t, response, &studio)
	response = request(t, router, http.MethodPut, fmt.Sprintf("/studios/%d", studio.Data.ID), map[string]any{
		"name": "Edited Studio", "seat_rows": 3, "seat_columns": 4,
	}, adminToken, "application/json")
	decode(t, response, &studio)
	if studio.Data.Name != "Edited Studio" || studio.Data.SeatRows != 3 {
		t.Fatalf("updated studio = %#v", studio.Data)
	}

	startTime := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	response = request(t, router, http.MethodPost, "/schedules", map[string]any{
		"movie_id": movie.Data.ID, "studio_id": studio.Data.ID, "start_time": startTime, "base_price": 50000,
	}, adminToken, "application/json")
	var schedule struct {
		Data models.Schedule `json:"data"`
	}
	decode(t, response, &schedule)
	response = request(t, router, http.MethodPut, fmt.Sprintf("/schedules/%d", schedule.Data.ID), map[string]any{
		"movie_id": movie.Data.ID, "studio_id": studio.Data.ID, "start_time": startTime, "base_price": 60000,
	}, adminToken, "application/json")
	decode(t, response, &schedule)
	if schedule.Data.BasePrice != 60000 {
		t.Fatalf("updated schedule price = %.0f", schedule.Data.BasePrice)
	}

	for _, endpoint := range []string{
		fmt.Sprintf("/schedules/%d", schedule.Data.ID),
		fmt.Sprintf("/movies/%d", movie.Data.ID),
		fmt.Sprintf("/studios/%d", studio.Data.ID),
	} {
		response = request(t, router, http.MethodDelete, endpoint, nil, adminToken, "")
		if response.Code != http.StatusOK {
			t.Fatalf("DELETE %s status = %d, body = %s", endpoint, response.Code, response.Body.String())
		}
	}
	if _, err := os.Stat(newPoster); !os.IsNotExist(err) {
		t.Fatalf("deleted movie poster still exists: %v", err)
	}
}

func TestBookingPricingAndOwnership(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter(config.NewDatabase())
	normalToken := login(t, router, "budi", "budi")
	studentToken := login(t, router, "andi", "andi")

	response := request(t, router, http.MethodGet, "/schedules", nil, "", "")
	var schedules struct {
		Data []struct {
			models.Schedule
			SeatPrice float64 `json:"seat_price"`
		} `json:"data"`
	}
	decode(t, response, &schedules)
	priceByID := map[int]float64{}
	for _, schedule := range schedules.Data {
		priceByID[schedule.ID] = schedule.SeatPrice
	}

	response = request(t, router, http.MethodPost, "/bookings", map[string]any{
		"schedule_id": 1, "seat_ids": []int{1, 57},
	}, normalToken, "application/json")
	var normalBooking struct {
		Data models.Booking `json:"data"`
	}
	decode(t, response, &normalBooking)
	assertPrice(t, normalBooking.Data.TotalPrice, priceByID[1]*2.5)

	response = request(t, router, http.MethodPost, "/bookings", map[string]any{
		"schedule_id": 2, "seat_ids": []int{65, 95},
	}, studentToken, "application/json")
	var studentBooking struct {
		Data models.Booking `json:"data"`
	}
	decode(t, response, &studentBooking)
	assertPrice(t, studentBooking.Data.TotalPrice, priceByID[2]*2.5*0.8)

	response = request(t, router, http.MethodGet, fmt.Sprintf("/bookings/%d", studentBooking.Data.ID), nil, normalToken, "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("foreign booking read status = %d, want 404", response.Code)
	}
	response = request(t, router, http.MethodPost, "/payments", map[string]any{
		"booking_id": studentBooking.Data.ID,
		"method":     "credit_card",
		"amount":     studentBooking.Data.TotalPrice,
	}, normalToken, "application/json")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("foreign booking payment status = %d, want 400", response.Code)
	}
}

func login(t *testing.T, router http.Handler, username, password string) string {
	t.Helper()
	response := request(t, router, http.MethodPost, "/login", map[string]string{
		"username": username,
		"password": password,
	}, "", "application/json")
	if response.Code != http.StatusOK {
		t.Fatalf("login %s status = %d, body = %s", username, response.Code, response.Body.String())
	}
	var body struct {
		Token string `json:"token"`
	}
	decode(t, response, &body)
	return body.Token
}

func uploadEvidence(t *testing.T, router http.Handler, token, filename, contents string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("evidence", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte(contents)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return request(t, router, http.MethodPost, "/users/me/student-application", body.Bytes(), token, writer.FormDataContentType())
}

func uploadMoviePoster(t *testing.T, router http.Handler, method, path, token string, fields map[string]string, filename, contents string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("poster", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte(contents)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return request(t, router, method, path, body.Bytes(), token, writer.FormDataContentType())
}

func request(t *testing.T, router http.Handler, method, path string, body any, token, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	var payload []byte
	switch value := body.(type) {
	case nil:
	case []byte:
		payload = value
	default:
		var err error
		payload, err = json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func decode(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if response.Code < 200 || response.Code >= 300 {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response: %v, body = %s", err, response.Body.String())
	}
}

func assertPrice(t *testing.T, actual, expected float64) {
	t.Helper()
	if math.Abs(actual-expected) > 0.01 {
		t.Fatalf("price = %.2f, want %.2f", actual, expected)
	}
}
