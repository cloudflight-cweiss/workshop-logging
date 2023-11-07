FROM golang:1.21 AS build

WORKDIR /workspace

COPY . /workspace

RUN go build -o workshop-load-bin workshop-load
RUN chmod +x ./workshop-load

FROM golang:1.21
COPY --from=build /workspace/workshop-load-bin /

ENTRYPOINT ["/workshop-load-bin"]