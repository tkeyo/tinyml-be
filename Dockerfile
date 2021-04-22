## STAGE 1
FROM golang:alpine AS build

RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY . /app/

RUN CGO_ENABLED=0 go build -o /bin/tinyml-be

## STAGE 2
FROM scratch

ARG API_AUTH_KEY
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG GIN_MODE

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/tinyml-be /bin/tinyml-be

EXPOSE 8081

CMD ["./bin/tinyml-be"]

# Run image with
# docker run -it -p 8081:8081 -e API_AUTH_KEY=123 -e GIN_MODE=release -e AWS_ACCESS_KEY_ID=xyz AWS_SECRET_ACCESS_KEY=xyz tinymlbe:latest