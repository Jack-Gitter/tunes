# Copy all file and build
FROM golang as golang
WORKDIR /tunes
COPY . . 
COPY .env .env
COPY .entrypoint.sh .entrypoint.sh
RUN go build


# Copy only binary and .env to other container
FROM golang
COPY --from=golang ./tunes/tunes /bin/tunes
COPY --from=golang ./tunes/.env /bin/.env
COPY --from=golang ./tunes/.entrypoint.sh /bin/.entrypoint.sh
COPY --from=golang ./tunes/db/migrations /bin/db/migrations
RUN chmod 755 /bin/.entrypoint.sh

ENTRYPOINT ["/bin/.entrypoint.sh"]
