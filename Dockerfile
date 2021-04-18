FROM golang:alpine AS build

WORKDIR /app/
COPY . /app/

RUN CGO_ENABLED=0 go build -o /bin/demo

FROM scratch
ARG API_AUTH_KEY

COPY --from=build /bin/demo /bin/demo
CMD ["./bin/demo"]