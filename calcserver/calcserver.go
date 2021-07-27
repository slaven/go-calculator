package calcserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type (
	// CalcServer contains logic for the server
	CalcServer struct {
		mux    *http.ServeMux
		routes []string
		cache  Caching
	}

	// Calculation operation
	CalcOperation string

	// Result of calculation
	CalcResult struct {
		Action CalcOperation `json:"action"`
		Answer float64       `json:"answer"`
		X      float64       `json:"x"`
		Y      float64       `json:"y"`
		Cached bool          `json:"cached"`
	}
)

// Define calculations
const (
	addCalc      CalcOperation = "add"
	subtractCalc CalcOperation = "subtract"
	multiplyCalc CalcOperation = "multiply"
	divideCalc   CalcOperation = "divide"
)

// Create returns new CalcServer instance
func Create() *CalcServer {
	mux := http.NewServeMux()
	srv := &CalcServer{
		mux:    mux,
		routes: []string{},
		cache:  NewCache(),
	}

	// Register routes
	for _, route := range getSupportedRoutes() {
		r1 := "/" + string(route)
		// handle trailing slash
		r2 := "/" + string(route) + "/"
		srv.routes = append(srv.routes, r1, r2)
	}

	return srv
}

// Handler interface for HTTP server
func (srv *CalcServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle only GET requests
	if req.Method != http.MethodGet {
		log.Printf("[Request error]: Method Not Allowed | path: %v; method: %v", req.URL.Path, req.Method)
		respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("HTTP %v Not Allowed", req.Method))
		return
	}

	// Check calculation routes for Calculation
	selectedOperation := CalcOperation("")
	for _, route := range srv.routes {
		if route == req.URL.Path {
			selectedOperation = CalcOperation(strings.Trim(route, "/"))
			break
		}
	}

	// Route not found
	if string(selectedOperation) == "" {
		log.Printf("[Request error]: Calculation Not Found | path: %v; method: %v", req.URL.Path, req.Method)
		respondWithError(w, http.StatusBadRequest, "Calculation Not Found")
		return
	}

	// Process query params
	x, y, err := getQueryValues(req.URL.Query())
	if err != nil {
		log.Printf("[Request error]: Missing Query Params: %v | path: %v; method: %v; query: %v", err.Error(), req.URL.Path, req.Method, req.URL.Query())
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Calculate result, handle error
	result, err := srv.calculateResult(selectedOperation, x, y)
	if err != nil {
		log.Printf("[Request error]: Calculation failed: %v | path: %v; method: %v; query: %v", err.Error(), req.URL.Path, req.Method, req.URL.Query())
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return result
	log.Printf("[Request ok]: path: %v; method: %v; query: %v; Response: %v", req.URL.Path, req.Method, req.URL.Query(), result)
	respondWithJSON(w, http.StatusOK, result)
}

// Calculate result
func (srv *CalcServer) calculateResult(calculation CalcOperation, x, y float64) (CalcResult, error) {
	// initial CalcResult data
	returnResult := CalcResult{
		Action: calculation,
		X:      x,
		Y:      y,
	}

	// build cacheKey
	cacheKey := buildCacheKey(calculation, x, y)

	// found data in Cache
	if val, found := srv.cache.Get(cacheKey); found {
		log.Printf("[Cache hit]: Get for cacheKey: %v; value %v", cacheKey, val)
		returnResult.Answer = val
		returnResult.Cached = found
		return returnResult, nil
	}

	// calculate result
	calcVal, err := calculate(calculation, x, y)
	if err != nil {
		log.Printf("[Calculate error]: %v", err.Error())
		return returnResult, err
	}
	returnResult.Answer = calcVal
	returnResult.Cached = false

	// cache result
	ok, err := srv.cache.SetOrUpdate(cacheKey, calcVal)
	if err != nil {
		log.Printf("[Cache error]: SetOrUpdate failed for cacheKey: %v; value %v; Error: %v", cacheKey, calcVal, err.Error())
	}
	if !ok {
		log.Printf("[Cache error]: SetOrUpdate failed to save cacheKey: %v; value %v", cacheKey, calcVal)
	}
	log.Printf("[Cache set]: SetOrUpdate saved cacheKey: %v; value %v", cacheKey, calcVal)

	// return non-cached result
	return returnResult, nil
}

// Calculate logic
func calculate(calculation CalcOperation, x, y float64) (float64, error) {
	switch calculation {
	case addCalc:
		return x + y, nil
	case subtractCalc:
		return x - y, nil
	case multiplyCalc:
		return x * y, nil
	case divideCalc:
		if y == 0 {
			return 0, errors.New("division by zero: Y param cannot be 0")
		}
		return x / y, nil
	default:
		return 0, errors.New("unknown calculation operation")
	}
}

// getSupportedRoutes returns supported calculations/routes
func getSupportedRoutes() []CalcOperation {
	return []CalcOperation{addCalc, subtractCalc, multiplyCalc, divideCalc}
}

// Get x and y from query parameters
func getQueryValues(vals url.Values) (float64, float64, error) {
	valX := vals.Get("x")
	if valX == "" {
		return 0, 0, errors.New("value for X parameter is missing")
	}
	x, err := strconv.ParseFloat(valX, 64)
	if err != nil {
		return 0, 0, err
	}

	valY := vals.Get("y")
	if valY == "" {
		return 0, 0, errors.New("value for Y parameter is missing")
	}
	y, err := strconv.ParseFloat(valY, 64)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

// respondWithJSON returns JSON encoded response
func respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError returns JSON encoded error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}
