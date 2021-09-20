FROM golang:1.17-buster
  WORKDIR /home/app/repo
  COPY . /home/app/repo

  ENV CLIENT_ID=
  ENV CLIENT_SECRET=
  ENV REDIRECT_URL=
  
  CMD go run main.go

  EXPOSE 3000:3000