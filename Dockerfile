# Copy all file and build
FROM golang as golang
WORKDIR /tunes
COPY . . 
RUN go build


# Copy only binary and .env to other container
FROM golang
WORKDIR /tunes
COPY --from=golang ./tunes/tunes ./tunes
COPY --from=golang ./tunes/.env ./.env
COPY --from=golang ./tunes/.entrypoint.sh ./.entrypoint.sh
COPY --from=golang ./tunes/db/migrations ./db/migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest 

ENTRYPOINT ["./.entrypoint.sh"]
