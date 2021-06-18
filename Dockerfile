
FROM tnsmith/yottadb-golang:r1.30.0-1.15.6 as builder

#SET WORKING DIRECTORY
WORKDIR /go/src/

#COPY CODE INTO WORKSPACE
COPY . .

#BUILD APP
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/app main.go

FROM tnsmith/yottadb:r1.30.0

WORKDIR /

COPY --chown=ydbadm:ydbadm --from=builder /go/src/startService.sh /application/startService.sh
COPY --chown=ydbadm:ydbadm --from=builder /go/app /application/app

# ENTRYPOINT
RUN chmod +x -R /application
USER ydbadm
ENTRYPOINT ["/application/startService.sh"]

