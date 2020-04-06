FROM golang:alpine as builder

ENV GOROOT /usr/local/go

RUN apk -U --no-cache add git make

ADD . /src
WORKDIR /src

RUN make binary

# ---

FROM alpine

COPY --from=builder /src/build/entrypoint.sh /app/entrypoint.sh
COPY --from=builder /src/bin/claim-parser /app/claim-parser

RUN apk -U --no-cache add bash ca-certificates \
    && chmod +x /app/claim-parser \
    && chmod +x /app/entrypoint.sh

WORKDIR /app

VOLUME ["/app/config"]

ENTRYPOINT ["/app/entrypoint.sh"]
