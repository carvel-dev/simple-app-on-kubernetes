FROM golang:1.12
WORKDIR ~/k8s-simple-app-example/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

FROM scratch
COPY --from=0 ~/k8s-simple-app-example/app .
EXPOSE 80
ENTRYPOINT ["/app"]
