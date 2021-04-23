

PROTO_SRCS = common.proto resource.proto metrics.proto trace.proto logs.proto trace_service.proto metrics_service.proto logs_service.proto
PROTO_GO_SRC = $(PROTO_SRCS:.proto=.pb.go)
PROTO_VERISON=v0.7.0
PROTO_PACKAGE=otlp_7

.PHONY: gen-proto massage-protos download-protos clean


clean:
	rm $(PROTO_SRCS) $(PROTO_GO_SRC)

download-protos:
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/common/v1/common.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/resource/v1/resource.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/metrics/v1/metrics.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/logs/v1/logs.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/trace/v1/trace.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/collector/trace/v1/trace_service.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/collector/metrics/v1/metrics_service.proto ;\
	wget https://raw.githubusercontent.com/open-telemetry/opentelemetry-proto/$(PROTO_VERISON)/opentelemetry/proto/collector/logs/v1/logs_service.proto


massage-protos: $(PROTO_SRCS)
	for src in $(PROTO_SRCS) ; do \
		sed -i "s/option go_package.*;//g" $$src ;\
		sed -i "s/option java_package.*;//g" $$src ;\
		sed -i "s/package .*;/package $(PROTO_PACKAGE);/g" $$src ;\
		sed -i "s/opentelemetry\/proto\/common\/v1\///g" $$src ;\
		sed -i "s/opentelemetry\.proto\.common\.v1\.//g" $$src ;\
		sed -i "s/opentelemetry\/proto\/resource\/v1\///g" $$src ;\
		sed -i "s/opentelemetry\.proto\.resource\.v1\.//g" $$src ;\
		sed -i "s/opentelemetry\/proto\/trace\/v1\///g" $$src ;\
		sed -i "s/opentelemetry\.proto\.trace\.v1\.//g" $$src ;\
		sed -i "s/opentelemetry\/proto\/metrics\/v1\///g" $$src ;\
		sed -i "s/opentelemetry\.proto\.metrics\.v1\.//g" $$src ;\
		sed -i "s/opentelemetry\/proto\/logs\/v1\///g" $$src ;\
		sed -i "s/opentelemetry\.proto\.logs\.v1\.//g" $$src ;\
	done

gen-proto: massage-protos $(PROTO_GO_SRC)

%.pb.go: %.proto 
	protoc -I/usr/local/include -I. $< --gogofaster_out=plugins=grpc:.