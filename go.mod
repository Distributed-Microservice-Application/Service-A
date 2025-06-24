module service-a

go 1.23.0

toolchain go1.24.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.20.4
	github.com/segmentio/kafka-go v0.4.48
	google.golang.org/grpc v1.58.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
)

replace golang.org/x/net => golang.org/x/net v0.22.0

replace golang.org/x/text => golang.org/x/text v0.14.0

replace golang.org/x/sys => golang.org/x/sys v0.15.0
