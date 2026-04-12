package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/keren/pao-event-panel/config"
	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/router"
)

var (
	testRouter     *gin.Engine
	superAdminToken string
	adminToken      string
	createdAdminID  float64
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Init config with test values
	config.AppConfig = &config.Config{
		Port:           "8080",
		DBDriver:       "sqlite",
		DBDSN:          ":memory:",
		JWTSecret:      "test-secret",
		SuperAdminUser: "admin",
		SuperAdminPass: "admin123",
		GinMode:        "test",
	}

	if err := database.InitDB("sqlite", ":memory:"); err != nil {
		panic(err)
	}
	if err := database.Migrate(); err != nil {
		panic(err)
	}
	database.Seed("admin", "admin123")

	testRouter = gin.New()
	router.Setup(testRouter)

	// Get super admin token
	superAdminToken = loginAs("admin", "admin123")

	os.Exit(m.Run())
}

// --- Helpers ---

func doRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func loginAs(username, password string) string {
	w := doRequest("POST", "/api/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, "")
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"].(string)
}

func parseJSON(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func parseJSONArray(w *httptest.ResponseRecorder) []interface{} {
	var result []interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

// --- Auth Tests ---

func TestLogin_Success(t *testing.T) {
	w := doRequest("POST", "/api/auth/login", map[string]string{
		"username": "admin",
		"password": "admin123",
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["token"] == nil {
		t.Fatal("expected token in response")
	}
	user := resp["user"].(map[string]interface{})
	if user["username"] != "admin" {
		t.Fatalf("expected username admin, got %v", user["username"])
	}
	if user["role"] != "super_admin" {
		t.Fatalf("expected role super_admin, got %v", user["role"])
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	w := doRequest("POST", "/api/auth/login", map[string]string{
		"username": "admin",
		"password": "wrong",
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	w := doRequest("POST", "/api/auth/login", map[string]string{}, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestMe(t *testing.T) {
	w := doRequest("GET", "/api/auth/me", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["username"] != "admin" {
		t.Fatalf("expected username admin, got %v", resp["username"])
	}
}

func TestMe_Unauthorized(t *testing.T) {
	w := doRequest("GET", "/api/auth/me", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- Users Tests ---

func TestCreateUser(t *testing.T) {
	w := doRequest("POST", "/api/users", map[string]string{
		"username": "testadmin",
		"password": "testpass123",
		"role":     "admin",
	}, superAdminToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["username"] != "testadmin" {
		t.Fatalf("expected username testadmin, got %v", resp["username"])
	}
	createdAdminID = resp["id"].(float64)

	// Login as the new admin to get a token for later tests
	adminToken = loginAs("testadmin", "testpass123")
}

func TestListUsers(t *testing.T) {
	w := doRequest("GET", "/api/users", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	users := parseJSONArray(w)
	if len(users) < 2 {
		t.Fatalf("expected at least 2 users, got %d", len(users))
	}
}

func TestDeleteUser_Self(t *testing.T) {
	// Super admin tries to delete themselves
	w := doRequest("GET", "/api/auth/me", nil, superAdminToken)
	me := parseJSON(w)
	myID := me["id"].(float64)

	w = doRequest("DELETE", fmt.Sprintf("/api/users/%d", int(myID)), nil, superAdminToken)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	if createdAdminID == 0 {
		t.Skip("no admin created to delete")
	}
	w := doRequest("DELETE", fmt.Sprintf("/api/users/%d", int(createdAdminID)), nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Recreate the admin for subsequent tests
	doRequest("POST", "/api/users", map[string]string{
		"username": "testadmin2",
		"password": "testpass123",
		"role":     "admin",
	}, superAdminToken)
	adminToken = loginAs("testadmin2", "testpass123")
}

func TestCreateUser_Forbidden(t *testing.T) {
	if adminToken == "" {
		t.Skip("no admin token")
	}
	w := doRequest("POST", "/api/users", map[string]string{
		"username": "shouldfail",
		"password": "testpass123",
		"role":     "admin",
	}, adminToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// --- Students Tests ---

var testStudentID float64

func TestCreateStudent(t *testing.T) {
	w := doRequest("POST", "/api/students", map[string]interface{}{
		"name":       "Test Student",
		"student_id": "s12345",
		"house_id":   1,
		"lap_count":  5,
	}, superAdminToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	testStudentID = resp["id"].(float64)
	if resp["name"] != "Test Student" {
		t.Fatalf("expected name Test Student, got %v", resp["name"])
	}
	if resp["student_id"] != "s12345" {
		t.Fatalf("expected student_id s12345, got %v", resp["student_id"])
	}
	if resp["lap_count"].(float64) != 5 {
		t.Fatalf("expected lap_count 5, got %v", resp["lap_count"])
	}
}

func TestCreateStudent_Unauthorized(t *testing.T) {
	w := doRequest("POST", "/api/students", map[string]interface{}{
		"name":     "Fail",
		"house_id": 1,
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestListStudents(t *testing.T) {
	w := doRequest("GET", "/api/students", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	students := parseJSONArray(w)
	if len(students) < 1 {
		t.Fatalf("expected at least 1 student, got %d", len(students))
	}
}

func TestUpdateStudent(t *testing.T) {
	if testStudentID == 0 {
		t.Skip("no student created")
	}
	w := doRequest("PUT", fmt.Sprintf("/api/students/%d", int(testStudentID)), map[string]interface{}{
		"name":       "Updated Student",
		"student_id": "s99999",
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["name"] != "Updated Student" {
		t.Fatalf("expected name Updated Student, got %v", resp["name"])
	}
}

func TestAddLaps(t *testing.T) {
	if testStudentID == 0 {
		t.Skip("no student created")
	}
	w := doRequest("POST", fmt.Sprintf("/api/students/%d/laps", int(testStudentID)), map[string]interface{}{
		"delta": 3,
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["lap_count"].(float64) != 8 {
		t.Fatalf("expected lap_count 8, got %v", resp["lap_count"])
	}
}

func TestAddLaps_Negative(t *testing.T) {
	if testStudentID == 0 {
		t.Skip("no student created")
	}
	w := doRequest("POST", fmt.Sprintf("/api/students/%d/laps", int(testStudentID)), map[string]interface{}{
		"delta": -2,
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["lap_count"].(float64) != 6 {
		t.Fatalf("expected lap_count 6, got %v", resp["lap_count"])
	}
}

func TestSetLaps(t *testing.T) {
	if testStudentID == 0 {
		t.Skip("no student created")
	}
	w := doRequest("PUT", fmt.Sprintf("/api/students/%d/laps", int(testStudentID)), map[string]interface{}{
		"lap_count": 10,
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["lap_count"].(float64) != 10 {
		t.Fatalf("expected lap_count 10, got %v", resp["lap_count"])
	}
}

func TestDeleteStudent(t *testing.T) {
	// Create a student to delete
	w := doRequest("POST", "/api/students", map[string]interface{}{
		"name":     "To Delete",
		"house_id": 2,
	}, superAdminToken)
	resp := parseJSON(w)
	id := resp["id"].(float64)

	w = doRequest("DELETE", fmt.Sprintf("/api/students/%d", int(id)), nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Rankings Tests ---

func TestHouseRankings(t *testing.T) {
	w := doRequest("GET", "/api/rankings/houses", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	houses := parseJSONArray(w)
	if len(houses) != 4 {
		t.Fatalf("expected 4 houses, got %d", len(houses))
	}
	first := houses[0].(map[string]interface{})
	if first["house_id"] == nil {
		t.Fatal("expected house_id in response")
	}
	if first["total_laps"] == nil {
		t.Fatal("expected total_laps in response")
	}
}

func TestStudentRankings(t *testing.T) {
	w := doRequest("GET", "/api/rankings/students", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	students := parseJSONArray(w)
	if len(students) < 1 {
		t.Fatalf("expected at least 1 student, got %d", len(students))
	}
	first := students[0].(map[string]interface{})
	if first["rank"] == nil {
		t.Fatal("expected rank in response")
	}
	if first["id"] == nil {
		t.Fatal("expected id in response")
	}
}

func TestRankingStats(t *testing.T) {
	w := doRequest("GET", "/api/rankings/stats", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := parseJSON(w)
	if resp["total_laps"] == nil {
		t.Fatal("expected total_laps in response")
	}
	if resp["total_participants"] == nil {
		t.Fatal("expected total_participants in response")
	}
}

// --- Activity Tests ---

func TestGetActivity(t *testing.T) {
	w := doRequest("GET", "/api/activity", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := parseJSON(w)
	if resp["activity_name"] == nil {
		t.Fatal("expected activity_name in response")
	}
}

func TestUpdateActivity(t *testing.T) {
	w := doRequest("PUT", "/api/activity", map[string]interface{}{
		"activity_name": "Test Event",
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["activity_name"] != "Test Event" {
		t.Fatalf("expected activity_name Test Event, got %v", resp["activity_name"])
	}
}

func TestUpdateActivity_WeatherLocation(t *testing.T) {
	w := doRequest("PUT", "/api/activity", map[string]interface{}{
		"weather_location": "Shanghai",
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["weather_location"] != "Shanghai" {
		t.Fatalf("expected weather_location Shanghai, got %v", resp["weather_location"])
	}
}

func TestUpdateActivity_Forbidden(t *testing.T) {
	if adminToken == "" {
		t.Skip("no admin token")
	}
	w := doRequest("PUT", "/api/activity", map[string]interface{}{
		"activity_name": "Hacked",
	}, adminToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// --- Houses Tests ---

func TestListHouses(t *testing.T) {
	w := doRequest("GET", "/api/houses", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	houses := parseJSONArray(w)
	if len(houses) != 4 {
		t.Fatalf("expected 4 houses, got %d", len(houses))
	}
}

func TestUpdateHouseColor(t *testing.T) {
	w := doRequest("PUT", "/api/houses/1/color", map[string]string{
		"color": "#FF0000",
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["color"] != "#FF0000" {
		t.Fatalf("expected color #FF0000, got %v", resp["color"])
	}
}

// --- Theme Presets Tests ---

var testPresetID float64

func TestApplyBuiltinPreset(t *testing.T) {
	w := doRequest("POST", "/api/theme-presets/builtin/apply", map[string]string{
		"name": "ocean",
	}, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["primary_color"] == nil {
		t.Fatal("expected primary_color in response")
	}
}

func TestCreateThemePreset(t *testing.T) {
	w := doRequest("POST", "/api/theme-presets", map[string]string{
		"name": "My Custom Theme",
	}, superAdminToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	testPresetID = resp["id"].(float64)
	if resp["name"] != "My Custom Theme" {
		t.Fatalf("expected name My Custom Theme, got %v", resp["name"])
	}
}

func TestListThemePresets(t *testing.T) {
	w := doRequest("GET", "/api/theme-presets", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	presets := parseJSONArray(w)
	if len(presets) < 1 {
		t.Fatalf("expected at least 1 preset, got %d", len(presets))
	}
}

func TestApplyThemePreset(t *testing.T) {
	if testPresetID == 0 {
		t.Skip("no preset created")
	}
	w := doRequest("POST", fmt.Sprintf("/api/theme-presets/%d/apply", int(testPresetID)), nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteThemePreset(t *testing.T) {
	if testPresetID == 0 {
		t.Skip("no preset created")
	}
	w := doRequest("DELETE", fmt.Sprintf("/api/theme-presets/%d", int(testPresetID)), nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Upload Tests ---

func TestUploadLogo(t *testing.T) {
	// Create a tiny PNG in memory
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	var imgBuf bytes.Buffer
	png.Encode(&imgBuf, img)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.png"`)
	h.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(imgBuf.Bytes())
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/upload/logo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+superAdminToken)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseJSON(w)
	if resp["path"] == nil {
		t.Fatal("expected path in response")
	}
}

// --- Export Tests ---

func TestExportRankings(t *testing.T) {
	w := doRequest("GET", "/api/export/rankings", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	contentType := w.Header().Get("Content-Type")
	expected := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if contentType != expected {
		t.Fatalf("expected content-type %s, got %s", expected, contentType)
	}
	if w.Body.Len() == 0 {
		t.Fatal("expected non-empty body")
	}
}

// --- Autocomplete Tests ---

func TestAutocompleteSearch_NoDB(t *testing.T) {
	w := doRequest("GET", "/api/autocomplete/search?q=test", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	results := parseJSONArray(w)
	if len(results) != 0 {
		t.Fatalf("expected 0 results without DB, got %d", len(results))
	}
}

func TestAutocompleteStatus_NoDB(t *testing.T) {
	w := doRequest("GET", "/api/autocomplete/status", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := parseJSON(w)
	if resp["loaded"] != false {
		t.Fatalf("expected loaded=false, got %v", resp["loaded"])
	}
}

func TestAutocompleteSearch_EmptyQuery(t *testing.T) {
	w := doRequest("GET", "/api/autocomplete/search?q=", nil, superAdminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	results := parseJSONArray(w)
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty query, got %d", len(results))
	}
}
