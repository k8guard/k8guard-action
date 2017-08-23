FROM varikin/golang-glide-alpine AS build
WORKDIR /go/src/github.com/k8guard/k8guard-action
COPY ./ ./
RUN apk -U add make
RUN make deps build

FROM alpine
RUN apk -U add ca-certificates
COPY --from=build /go/src/github.com/k8guard/k8guard-action/k8guard-action /
EXPOSE 3000
CMD ["/k8guard-action"]