# OTLP Log Viewer

## Introduction
This take-home assignment is designed to give you an opportunity to demonstrate your skills and experience in
building a web application. We expect you to spend 2-4 hours on this assignment. If you find yourself spending more time
than that, please stop and submit what you have. We are not looking for a complete solution, but rather a demonstration
of your skills and experience.

To submit your solution, please create a public GitHub repository and send us the link. Please include a README.md file
with instructions on how to run your application.

## Overview
The goal of this assignment is to build a simple web application that renders a list of [log records](https://opentelemetry.io/docs/concepts/signals/logs/) from an OTLP logs
endpoint. For an overview of the expected capabilities, see the [Expected Capabilities](#expected-capabilities) section
below.

## OTLP Logs HTTP Endpoint
The OTLP logs HTTP endpoint is available at https://otlp-logs-endpoint.herokuapp.com/v1/logs. The endpoint returns a
data structure compatible with 

## Expected Capabilities
Take this piece of OTLP logs JSON and render the logs in a list with columns
Severity
Time
Body
When clicking a log record, expand the row and show all attributes associated with the log record
Render a histogram visualizing the distribution of log records
X-Axis: Time
Y-Axis: Count
Technology Constraints
Use React and Next.js with the app router.
Use any additional libraries you want and need.
Notes
As this assignment is for the role of a Senior Product Engineer, we expect you to pay some attention to the experience and design of the solution. For example:
Structure consistent with established solutions in the observability domain
Consistent terminology usage
Some styling/visuals

## References

- [OpenTelemetry Logs](https://opentelemetry.io/docs/concepts/signals/logs/)
- [OpenTelemetry Protocol (OTLP)](https://github.com/open-telemetry/opentelemetry-proto)
- [OTLP Logs Examples](https://github.com/open-telemetry/opentelemetry-proto/blob/main/examples/logs.json)