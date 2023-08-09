FROM alpine:3.18.3

ENV TZ=Asia/Jakarta

COPY ./config.yaml /
COPY ./bin/ /

EXPOSE 8080

ENTRYPOINT [ "./player-be" ]