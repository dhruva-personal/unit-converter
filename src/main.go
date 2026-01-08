package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ConversionResult struct {
	Result float64 `json:"result"`
}

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "frontend/length.html")
}

func convert(convType string, amount float64, from string, to string) (float64, error) {
	switch convType {
	case "length":
		var meters float64

		// Convert from source unit to meters
		switch from {
		case "meters":
			meters = amount
		case "feet":
			meters = amount * 0.3048
		default:
			return 0, errors.New("invalid 'from' unit for length")
		}

		// Convert meters to target unit
		switch to {
		case "meters":
			return meters, nil
		case "feet":
			return meters / 0.3048, nil
		default:
			return 0, errors.New("invalid 'to' unit for length")
		}

	case "weight":
		var kg float64

		// Convert from source unit to kilograms
		switch from {
		case "kilograms":
			kg = amount
		case "pounds":
			kg = amount * 0.453592
		default:
			return 0, errors.New("invalid 'from' unit for weight")
		}

		// Convert kilograms to target unit
		switch to {
		case "kilograms":
			return kg, nil
		case "pounds":
			return kg / 0.453592, nil
		default:
			return 0, errors.New("invalid 'to' unit for weight")
		}

	case "temperature":
		var c float64

		// Convert from source unit to Celsius
		switch from {
		case "celsius":
			c = amount
		case "fahrenheit":
			c = (amount - 32) * 5 / 9
		default:
			return 0, errors.New("invalid 'from' unit for temperature")
		}

		// Convert Celsius to target unit
		switch to {
		case "celsius":
			return c, nil
		case "fahrenheit":
			return c*9/5 + 32, nil
		default:
			return 0, errors.New("invalid 'to' unit for temperature")
		}

	default:
		return 0, errors.New("invalid conversion type")
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", lengthHandler)
	r.Get("/length.html", lengthHandler)
	r.Get("/weight.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "frontend/weight.html")
	})
	r.Get("/temperature.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "frontend/temperature.html")
	})

	//conversion API endpoint
	r.Get("/convert", func(w http.ResponseWriter, r *http.Request) {
		convType := r.URL.Query().Get("type")
		amount := r.URL.Query().Get("amount")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		//convert amt to int
		amt, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fmt.Println("Conversion error:", err)
			return
		}

		//compute result
		result, err := convert(convType, amt, from, to)
		if err != nil {
			fmt.Println(w, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//send JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ConversionResult{Result: result})
	})

	http.ListenAndServe(":8080", r)
}
