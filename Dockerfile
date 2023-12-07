# use latest Go image
FROM golang:latest

# CHANGE THESE USERNAME AND PASSWORDS
ENV simap_poc_username=defaultUsername
ENV simap_poc_password=defaultPassword

# set working directory in container
WORKDIR /app

# copy dir into app
COPY . .

# copy the start script & maek executable
COPY start.sh .
RUN chmod +x start.sh

# download dependancies if required
# RUN go mod download

# compile
RUN go build -o main .

# start the IMAP server via the  shellscript
CMD ["./start.sh"]