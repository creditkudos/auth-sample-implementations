# Sample Auth implementations 

## Setup
1. Setup and launch `ngrok` on port `3000` to create an HTTPS redirect URI, set this URI as `REDIRECT_URI` in the `Dockerfile` for the implementation you intend to run.
1. Input your application `CLIENT_ID` and `CLIENT_SECRET` into the same `Dockerfile`.
1. Launch the example implementation with `docker build --tag auth-sample ./[language] && docker run auth-sample`  
e.g: `docker build --tag auth-sample ./go && docker run -it auth-sample`
1. Connect to either your `REDIRECT_URI` or `localhost:3000` to start a journey.
