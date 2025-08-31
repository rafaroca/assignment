# OTLP Log Processor (Go)

## Running

Start the log processor with defaults by issuing `go run .`.

See the possible CLI arguments with `go run . --help`. Most importantly you can specify the attribute key with `-attributeKey <key>` and the duration window for the log output with `-duration <duration>`.

Install the [telemetrygen tool](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/cmd/telemetrygen/README.md) with `go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest`.

Export 1 mio sample logs with `telemetrygen logs --otlp-insecure --logs 100000 --body '{"foo":"bar"}' --telemetry-attributes service.name=\"logattribute\" --workers 10` and check the output for "telemetrygen" and "logattribute".
For ease of templating and viewing in the command line, the output is done via `fmt.Println`.

## Tests

The test suite runs with `go test`.

The tests check the output of the processor, the log export, the counter increase and the string conversion.

## Structure

I expanded the scaffolding of the `logs_service.go` to extract the content of the log export and search for the attribute key inside the resource, scope and log attributes.
The log structure of `ExportLogsServiceRequest` contains state which makes it necessary to convert the attribute value into a string before sending it over a channel.
The usage of the channel makes the export non-blocking which is crucial for the gRPC call.
The `-bufferSize <size>` CLI argument controls the channel buffer size.
Increasing the buffer size allows for peaks in the number of messages to not block the gRPC call.

Since the `processor.go` is only adding to a map, this should be faster than even the parsing of the logs inside `logs_service.go`.
Thus multiple workers for `processor.go` would not increase throughput.
Profiling could help prove this hypothesis.
The `processor.go` avoids locking of the `logStats` map by avoiding concurrent access.
By listening to the `ticker` and the `logIntake` in a single `select` statement, the `logStats` map is never read and written to simultaneously.

The `log_service.go` uses metrics counters to keep track of the different sources of attributes.
