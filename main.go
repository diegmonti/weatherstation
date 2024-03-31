package main

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "sync"
	"time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    mutex           sync.Mutex
    updateCounter   prometheus.Counter
    baromin         *prometheus.GaugeVec
    temperature     *prometheus.GaugeVec
    dewptf          *prometheus.GaugeVec
    humidity        *prometheus.GaugeVec
    windSpeed       *prometheus.GaugeVec
    windGust        *prometheus.GaugeVec
    windDir         *prometheus.GaugeVec
    rainfall        *prometheus.GaugeVec
    dailyRainfall   *prometheus.GaugeVec
    indoorTemp      *prometheus.GaugeVec
    indoorHumidity  *prometheus.GaugeVec
)

func main() {
	registry := prometheus.NewRegistry()

    // Initialize Prometheus metrics
    updateCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "weatherstation_updates_total",
        Help: "Total number of weather station updates",
    })
    registry.MustRegister(updateCounter)

    baromin = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_barometric_pressure",
        Help: "Barometric pressure in millibar",
    }, []string{"id"})
    registry.MustRegister(baromin)

    temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_temperature",
        Help: "Temperature in Celsius",
    }, []string{"id"})
    registry.MustRegister(temperature)

    dewptf = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_dew_point",
        Help: "Dew point temperature in Celsius",
    }, []string{"id"})
    registry.MustRegister(dewptf)

    humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_humidity",
        Help: "Humidity percentage",
    }, []string{"id"})
    registry.MustRegister(humidity)

    windSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_wind_speed",
        Help: "Wind speed in meters per second",
    }, []string{"id"})
    registry.MustRegister(windSpeed)

    windGust = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_wind_gust",
        Help: "Wind gust speed in meters per second",
    }, []string{"id"})
    registry.MustRegister(windGust)

    windDir = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_wind_direction",
        Help: "Wind direction in degrees",
    }, []string{"id"})
    registry.MustRegister(windDir)

    rainfall = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_rainfall",
        Help: "Rainfall in millimeters",
    }, []string{"id"})
    registry.MustRegister(rainfall)

    dailyRainfall = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_daily_rainfall",
        Help: "Daily rainfall in millimeters",
    }, []string{"id"})
    registry.MustRegister(dailyRainfall)

    indoorTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_indoor_temperature",
        Help: "Indoor temperature in Celsius",
    }, []string{"id"})
    registry.MustRegister(indoorTemp)

    indoorHumidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "weatherstation_indoor_humidity",
        Help: "Indoor humidity percentage",
    }, []string{"id"})
    registry.MustRegister(indoorHumidity)

    // Define HTTP endpoints
    http.HandleFunc("/weatherstation/updateweatherstation.php", updateWeatherStation)
    http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

    // Start HTTP server
    fmt.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Failed to start server: %s", err)
    }
}

func updateWeatherStation(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    query := r.URL.Query()

    // Extract ID from query parameters
    ID := query.Get("ID")

    // Update metrics
    mutex.Lock()
    defer mutex.Unlock()

    // Increment the update counter
    updateCounter.Inc()

    barominValue := parseFloat(query.Get("baromin"))
    temperatureValue := parseFloat(query.Get("tempf"))
    dewptfValue := parseFloat(query.Get("dewptf"))
    humidityValue := parseFloat(query.Get("humidity"))
    windSpeedValue := parseFloat(query.Get("windspeedmph"))
    windGustValue := parseFloat(query.Get("windgustmph"))
    windDirValue := parseFloat(query.Get("winddir"))
    rainfallValue := parseFloat(query.Get("rainin"))
    dailyRainfallValue := parseFloat(query.Get("dailyrainin"))
    indoorTempValue := parseFloat(query.Get("indoortempf"))
    indoorHumidityValue := parseFloat(query.Get("indoorhumidity"))

    baromin.WithLabelValues(ID).Set(convertInchesOfMercuryToMillibar(barominValue))
    temperature.WithLabelValues(ID).Set(convertFahrenheitToCelsius(temperatureValue))
    dewptf.WithLabelValues(ID).Set(convertFahrenheitToCelsius(dewptfValue))
    humidity.WithLabelValues(ID).Set(humidityValue)
    windSpeed.WithLabelValues(ID).Set(convertMphToMps(windSpeedValue))
    windGust.WithLabelValues(ID).Set(convertMphToMps(windGustValue))
    windDir.WithLabelValues(ID).Set(windDirValue)
    rainfall.WithLabelValues(ID).Set(convertInchesToMillimeters(rainfallValue))
    dailyRainfall.WithLabelValues(ID).Set(convertInchesToMillimeters(dailyRainfallValue))
    indoorTemp.WithLabelValues(ID).Set(convertFahrenheitToCelsius(indoorTempValue))
    indoorHumidity.WithLabelValues(ID).Set(indoorHumidityValue)

    fmt.Fprintf(w, "success")
	logMessage("Metrics updated successfully")
}

func logMessage(message string) {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] %s\n", currentTime, message)
}

func parseFloat(value string) float64 {
    parsedValue := strings.TrimSpace(value)
    if parsedValue == "" {
        return 0.0
    }
    floatValue, err := strconv.ParseFloat(parsedValue, 64)
    if err != nil {
        fmt.Printf("Failed to parse float: %s", err)
        return 0.0
    }
    return floatValue
}

func convertFahrenheitToCelsius(fahrenheit float64) float64 {
    return (fahrenheit - 32) * 5 / 9
}

func convertMphToMps(mph float64) float64 {
    return mph * 0.44704
}

func convertInchesToMillimeters(inches float64) float64 {
    return inches * 25.4
}

func convertInchesOfMercuryToMillibar(inches float64) float64 {
    return inches * 33.8639
}
