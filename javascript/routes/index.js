var express = require('express');
var router = express.Router();

const jwt = require('jsonwebtoken');
const OAuth2Client = require('client-oauth2')

const CLIENT_ID = "..."
const CLIENT_SECRET = "..."
const REDIRECT_URI = "https://....ngrok.io/callback"

const client = new OAuth2Client({
  clientId: CLIENT_ID,
  clientSecret: CLIENT_SECRET,
  accessTokenUri: "https://api.creditkudos.com/v2/oauth/token",
  authorizationUri: "https://auth.creditkudos.com/oauth/authorize",
  redirectUri: REDIRECT_URI,
  scopes: ['read']
})

/* GET home page. */
router.get('/', function(req, res, next) {
  res.render('index', { title: 'Express' });
  console.log(req.url)
});

router.get('/redirect', (req, res) => {
  const payload = {
    iss: CLIENT_ID,
    sub: 'customer',
    iat: parseInt(new Date() / 1000),
    jti: '234342',
    email: 'sam.pull@example.com',
    first_name: 'Samuel',
    last_name: 'Pull',
    custom_reference: 'SP-123',
    date_of_birth: '1985-10-25',
    postcode: 'XY12AB',
  }

  const token = jwt.sign(payload, CLIENT_SECRET)

  var uri = client.code.getUri({
    query: {
      customer_token: token,
      // This state variable is the value that'll be returned as a query parameter in the /callback redirect
      // This can be used to tie a given callback to a specific session
      state: "state_value",
      debug: true
    }
  })

  console.log(uri)

  res.redirect(uri)
})

router.get('/callback', (req, res) => {
  client.code.getToken(req.originalUrl)
    .then((user) => {
      console.log(user)

      user.refresh().then((updatedUser) => {
        console.log(updatedUser !== user) //=> true
        console.log(updatedUser.accessToken)
      })

      return res.redirect('/complete')
    })
    .catch(e => {
      console.log("Something went wrong")
      console.log(e)
    })
})

router.get('/complete', (req, res) => {
  res.render('index', { title: 'Complete!' });
})

module.exports = router;
