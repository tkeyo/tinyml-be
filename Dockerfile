FROM golang:alpine AS build

WORKDIR /app/
COPY . /app/

RUN CGO_ENABLED=0 go build -o /bin/demo

FROM scratch
ARG API_AUTH_KEY

COPY --from=build /bin/demo /bin/demo
CMD ["./bin/demo"]

# Run image with
# docker run -it -p 8081:8081 -e API_AUTH_KEY=123 -e GIN_MODE=release -e AWS_ACCESS_KEY_ID=xyz AWS_SECRET_ACCESS_KEY=xyz tinymlbe:latest