FROM golang:alpine as builder

ENV GOROOT /usr/local/go

RUN apk -U --no-cache add git make

ADD . /src/claim-parser
WORKDIR /src/claim-parser

RUN make binary

# ---

FROM alpine

COPY --from=builder /src/claim-parser/bin/claim-parser /app/claim-parser
COPY build/entrypoint.sh /app/entrypoint.sh

RUN apk -U --no-cache add bash ca-certificates \
    && chmod +x /app/claim-parser \
    && chmod +x /app/entrypoint.sh

WORKDIR /app

VOLUME ["/app/config"]

ENTRYPOINT ["/app/entrypoint.sh"]
