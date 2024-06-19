`youless-observer` Helm Charts
==============================

[![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat) ][release-url]
[![Artifact Hub][artifact-hub-img]][artifact-hub-url]

[release-url]: https://github.com/roeldev/youless-observer/releases/tag/v0.2.0

[artifact-hub-img]: https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/youless-observer

[artifact-hub-url]: https://artifacthub.io/packages/search?repo=youless-observer

Real-time usage observer service for YouLess energy meter

## Usage

[Helm](https://helm.sh) must be installed to use the charts.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```sh
helm repo add youless-observer https://roeldev.github.io/youless-observer
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy for the container. |
| image.pullSecrets | list | `[]` | Specify the imagePullSecrets on the pod. |
| service | object | `{"annotations":{},"enabled":true,"port":80,"type":"LoadBalancer"}` | Expose the deployment using the service. |
| networkPolicy.create | bool | `true` |  |
| observer.debug | bool | `false` | When true, overrides and sets: `log.level: "debug"`, `log.accessLog: true`, `tls.insecureSkipVerify: true`. |
| observer.image.repository | string | `"roeldev/youless-observer"` | The image repository and name. |
| observer.image.tag | string | `"0.2.0"` | Tag of the image to deploy. |
| observer.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy for the container. |
| observer.log.level | string | `"warn"` | Active log level, must be one of: debug, info, warn, error, fatal, panic, disabled. |
| observer.log.timestamp | bool | `false` | Add timestamp to log entries when true. |
| observer.log.accessLog | bool | `false` | Enable server access logging when true. |
| observer.server.port | int | `2512` | Port on which the server inside the container listens to. |
| observer.server.tls.enabled | bool | `false` | Enables TLS on the server when true. |
| observer.server.tls.secretName | string | `""` | Name of secret containing TLS certificate and key. |
| observer.server.tls.caSecretName | string | `""` | Name of secret containing a CA certificate (ca.crt). |
| observer.server.tls.certFile | string | `""` | Path to the TLS certificate file. Is ignored when `secretName` is set. |
| observer.server.tls.keyFile | string | `""` | Path to the TLS key file. Is ignored when `secretName` is set. |
| observer.server.tls.caCertFile | string | `""` | Path to the CA certificate file. Is ignored when `caSecretName` is set. |
| observer.server.tls.verifyClient | bool | `true` | Enables mTLS when set to true. |
| observer.server.tls.insecureSkipVerify | bool | `false` |  |
| observer.otel.serviceName | string | `""` | Optional [custom name](https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_service_name) for the service. |
| observer.otel.resourceAttributes | string | `""` | Optional addition [resource attributes](https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_resource_attributes). |
| observer.otel.metric.exportInterval | int | `10000` |  |
| observer.otel.metric.exportTimeout | int | `10000` |  |
| observer.otel.traces.enabled | bool | `true` | Enables tracing when true. |
| observer.otel.traces.sampler | string | `"parentbased_traceidratio"` | Specifies the [sampler](https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_traces_sampler) used to sample traces. |
| observer.otel.traces.samplerArg | float | `0.5` | Specifies [arguments](https://opentelemetry.io/docs/languages/sdk-configuration/general/#otel_traces_sampler_arg), if applicable, to the sampler defined by `sampler`. |
| observer.otel.exporter.otlp.enabled | bool | `true` |  |
| observer.otel.exporter.otlp.endpoint | string | `""` | A base [endpoint URL](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/#otel_exporter_otlp_endpoint) for any signal type, with an optionally-specified port number. |
| observer.otel.exporter.otlp.headers | string | `""` | A list of [headers](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/#otel_exporter_otlp_headers) to apply to all outgoing data |
| observer.otel.exporter.otlp.protocol | string | `"grpc"` | Specifies the [OTLP transport protocol](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/#otel_exporter_otlp_protocol) to be used for all telemetry data. |
| observer.otel.exporter.otlp.timeout | int | `5000` | The [timeout](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/#otel_exporter_otlp_timeout) value for all outgoing data, in milliseconds. |
| observer.otel.exporter.otlp.tls.enabled | bool | `false` |  |
| observer.otel.exporter.otlp.tls.secretName | string | `""` | Name of secret containing TLS certificate and key. |
| observer.otel.exporter.otlp.tls.caSecretName | string | `""` | Name of secret containing a CA certificate (`ca.crt` key). |
| observer.otel.exporter.otlp.tls.certFile | string | `""` | Path to the TLS certificate file. Is ignored when `secretName` is set. |
| observer.otel.exporter.otlp.tls.keyFile | string | `""` | Path to the TLS key file. Is ignored when `secretName` is set. |
| observer.otel.exporter.otlp.tls.caCertFile | string | `""` | Path to the CA certificate file. Is ignored when `caSecretName` is set. |
| observer.otel.exporter.otlp.tls.insecureSkipVerify | bool | `false` |  |
| observer.youless.configMapName | string | `""` | When a configMapName is provided, the config map is used to populate the `YOULESS_BASE_URL`, `YOULESS_NAME`, `YOULESS_TIMEOUT` and `YOULESS_PASSWORD` env variables. The config map must have the `url` and `name` keys. The `timeout` and `password` keys are optional. |
| observer.youless.secretName | string | `""` | When provided, the secret is mounted and used as the value for `YOULESS_PASSWORD_FILE`. |
| observer.youless.url | string | `"http://youless"` | Url of the YouLess device to connect to. |
| observer.youless.name | string | `"YouLess"` | Name of the YouLess device for identification. |
| observer.youless.timeout | string | `"5s"` | Connection timeout for the YouLess device. |
| observer.youless.password | string | `""` | Password used to connect with the YouLess device. Is ignored when `secretName` is set. |
| observer.excludePower | bool | `false` |  |
| observer.excludeS0 | bool | `false` |  |
| observer.excludeGas | bool | `false` |  |
| observer.excludeWater | bool | `false` |  |
| observer.singlePhase | bool | `true` |  |
| readinessProbe.periodSeconds | int | `10` |  |
| readinessProbe.timeoutSeconds | int | `1` |  |
| readinessProbe.failureThreshold | int | `3` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.readOnlyRootFilesystem | bool | `true` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `1000` |  |
| resources.limits.cpu | string | `"10m"` |  |
| resources.limits.memory | string | `"32Mi"` |  |
| resources.requests.cpu | string | `"10m"` |  |
| resources.requests.memory | string | `"32Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |

## License

Copyright © 2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in
the [LICENSE](https://github.com/roeldev/youless-observer/blob/main/LICENSE) file.
