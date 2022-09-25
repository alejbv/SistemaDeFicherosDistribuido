FROM golang:alpine

WORKDIR /usr/app

#ADD . ./usr/app
COPY . .

RUN go build -o /usr/app/tagFile .
EXPOSE 8080-50051
#CMD ["/bin/sh"]   
ENTRYPOINT ["/usr/app/tagFile"]      
