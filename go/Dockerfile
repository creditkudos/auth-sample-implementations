FROM golang:1.17-buster
  WORKDIR /home/app/repo
  COPY . /home/app/repo

  ENV CLIENT_ID=
  ENV CLIENT_SECRET=
  ENV REDIRECT_URL=
  
  EXPOSE 3000
  CMD go run main.go