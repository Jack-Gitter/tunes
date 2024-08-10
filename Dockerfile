# Copy all file and build
FROM golang
WORKDIR /tunes
COPY . . 
RUN go build -o /bin/tunes ./main.go


# Copy only binary to 
FROM golang
COPY --from=0 /bin/tunes /bin/tunes
CMD ["/bin/tunes"]
