FROM golang:1.26-alpine AS base
RUN mkdir -p /opt/app
WORKDIR /opt/app
COPY . .
RUN go mod download

FROM base AS build
RUN apk add build-base
RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark-service cmd/api/main.go

FROM base AS test-exec
ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE
RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    if [ -n "${COVERAGE_EXCLUDE}" ]; then \
        grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out; \
    else \
        cp coverage.tmp ${_outputdir}/coverage.out; \
    fi && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

FROM alpine:3.24.1 AS final
WORKDIR /app
COPY --from=build /opt/app/bookmark-service /app/bookmark-service
COPY --from=build /opt/app/docs /app/docs
CMD ["/app/bookmark-service"]