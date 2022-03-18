var express = require("express");
var router = express.Router();

const axios = require("axios");

const OAuth2Client = require("client-oauth2");

const CLIENT_ID = process.env.CLIENT_ID;
const CLIENT_SECRET = process.env.CLIENT_SECRET;

const REDIRECT_URI = process.env.REDIRECT_URI;

const ACCESS_TOKEN_URI = process.env.ACCESS_TOKEN_URI || "https://api.creditkudos.com/v2/oauth/token";
const AUTHORIZATION_URI = process.env.AUTHORIZATION_URI || "https://auth.creditkudos.com/oauth/authorize";

const CK_CUSTOMER_TOKEN_ENDPOINT = process.env.AUTHORIZATION_URI ||"https://api.creditkudos.com/v3/customer_token"

/* GET home page. */
router.get("/", function (req, res, next) {
  res.render("index", { title: "Express" });
  console.log(req.url);
});

// example implementation of using Customer Token endpoint
router.get("/redirect", (req, res) => {
  let token;
  axios
    .post(
      CK_CUSTOMER_TOKEN_ENDPOINT,
      {},
      {
        auth: {
          username: CLIENT_ID,
          password: CLIENT_SECRET,
        },
        params: {
          email: "sam.pull@example.com",
          first_name: "Samuel",
          last_name: "Pull",
          custom_reference: "SP-123",
          date_of_birth: "1985-10-25",
          postcode: "XY12AB",
        },
      }
    )
    .then(function (response) {
      token = response.data.data.customerToken;
      var uri = client.code.getUri({
        query: {
          customer_token: token,
          // This state variable is the value that'll be returned as a query parameter in the /callback redirect
          // This can be used to tie a given callback to a specific session
          state: "state_value",
          debug: true,
          redirect_uri: REDIRECT_URI,
        },
      });

      res.redirect(uri);
    })
    .catch(function (error) {
      console.log(error);
    });
});

const client = new OAuth2Client({
  clientId: CLIENT_ID,
  clientSecret: CLIENT_SECRET,
  accessTokenUri: ACCESS_TOKEN_URI,
  authorizationUri: AUTHORIZATION_URI,
  redirectUri: REDIRECT_URI,
  scopes: ["read"],
});

router.get("/callback", (req, res) => {
  client.code
    .getToken(req.originalUrl)
    .then((user) => {
      console.log(user);

      user.refresh().then((updatedUser) => {
        console.log(updatedUser !== user); //=> true
        console.log(updatedUser.accessToken);
      });

      return res.redirect("/complete");
    })
    .catch((e) => {
      console.log("Something went wrong");
      console.log(e);
    });
});

router.get("/complete", (req, res) => {
  res.render("index", { title: "Complete!" });
});

module.exports = router;
