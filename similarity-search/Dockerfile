ARG GOLANG_VERSION=1.20.2
FROM golang:${GOLANG_VERSION} as builder
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app
COPY . .
RUN go get -d &&\
    go build -o main .

FROM ubuntu:focal
ENV TZ=Asia/Calcutta
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive TZ=$TZ apt-get -y install tzdata
RUN apt-get install -y ca-certificates && update-ca-certificates
RUN ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
WORKDIR /app
COPY --from=builder /app/.env ./
COPY --from=builder /app/DigiCertGlobalRootCA.crt.pem ./
COPY --from=builder /app/main ./ 
CMD ["./main"]