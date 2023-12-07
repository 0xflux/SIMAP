# use latest Go image
FROM golang:latest

# set working directory in container
WORKDIR /app

# copy dir into app
COPY . .

# download dependancies if required
RUN go mod download

# compile
RUN go build -o main .

# run app when container starts
CMD ["./main"]