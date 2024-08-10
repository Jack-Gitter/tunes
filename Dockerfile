# Copy all file and build
FROM golang as golang
WORKDIR /tunes
COPY . . 
COPY .env .env
RUN go build -o /bin/tunes ./main.go


# Copy only binary to 
FROM golang
WORKDIR /tunes
COPY --from=golang /bin/tunes /bin/
COPY --from=golang ./tunes/.env /bin/.env

CMD ["/bin/tunes"]
