FROM alpine
RUN apk -U add ca-certificates
ADD k8guard-action /
EXPOSE 3000
CMD ["/k8guard-action"]