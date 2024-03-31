# Weather Station Prometheus Exporter

This is a simple web server written in Go that exposes weather station metrics in Prometheus format. It can be called on a path similar to `/weatherstation/updateweatherstation.php` to update the metrics, and it exposes these metrics in Prometheus format on the `/metrics` path. Each call to the `/weatherstation/updateweatherstation.php` endpoint updates the metrics exposed on the `/metrics` path.

## Features

- Updates weather station metrics via HTTP endpoint.
- Exposes metrics in Prometheus format.
- Converts all numerical values to the metric system.

## Installation

### Docker

To run the server using Docker, you can use the provided Dockerfile. Build the Docker image and run a container as follows:

```bash
docker build -t weatherstation .
docker run -p 8080:8080 weatherstation
```

### Local Build

To build and run the server locally, you need to have Go installed on your system. Clone the repository and run the following commands:

```bash
go mod download
go build
./weatherstation
```

## Usage

Once the server is running, you can update the weather station metrics by making a GET request to the `/weatherstation/updateweatherstation.php` endpoint with appropriate query parameters. The server will then expose these metrics in Prometheus format on the `/metrics` endpoint.

Example usage:

```bash
curl "http://localhost:8080/weatherstation/updateweatherstation.php?ID=myId&PASSWORD=123456&action=updateraww&realtime=1&rtfreq=5&dateutc=now&baromin=29.79&tempf=50.7&dewptf=49.4&humidity=95&windspeedmph=0.0&windgustmph=0.0&winddir=45&rainin=0.0&dailyrainin=0.0&indoortempf=69.2&indoorhumidity=60"
```

You can then access the metrics in Prometheus format by navigating to `http://localhost:8080/metrics` in your web browser or using tools like `curl`.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or features you would like to see added.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
