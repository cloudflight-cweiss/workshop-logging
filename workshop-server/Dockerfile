FROM golang:1.21 AS build

WORKDIR /workspace

COPY go.mod /workspace
COPY go.sum /workspace

RUN go mod download

COPY . /workspace

RUN go build .
RUN chmod +x ./workshop-server

FROM golang:1.21
COPY --from=build /workspace/workshop-server /

ENTRYPOINT ["/workshop-server"]