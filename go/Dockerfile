FROM golang:1.14
  WORKDIR /home/app/repo
  COPY . /home/app/repo

  ENV CLIENTID=
  ENV CLIENTSECRET=
  ENV REDIRECTURL=

  CMD go run main.go

  EXPOSE 3000/tcp