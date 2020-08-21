# Sample Auth implementations 

## Javascript 
### Setup
This setup assumes that your machine has `nodejs` (with `npm`) installed. `cd` into the root directory, then run `npm i` to install the dependencies.

### Configuration
You'll need to set the environment variables `CLIENT_ID`, `CLIENT_SECRET` and `REDIRECT_URI` to match the application configuration in Atlas. You'll most likely want to point `ngrok`(https://ngrok.com/) at this local server, and use the HTTPS version as your redirect URI.

If using ngrok, you'll want to run `ngrok http 3000`, and set the `REDIRECT_URI` in both the Atlas application config and the environment variable to match the HTTPS address generated by ngrok.

### Running the application

After the dependencies have finished installing, run the below block to start.
```
DEBUG="myapp:*" CLIENT_ID = "..." CLIENT_SECRET = "..." REDIRECT_URI = "https://....ngrok.io/callback" npm start
```

To kick off a journey, go to http://localhost:3000/redirect