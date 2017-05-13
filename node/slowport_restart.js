const request = require('request');
const crypto = require('crypto')

const sjcl = require('./sjcl');

const argv = require('optimist')
  .usage("Restarts your Speedport W 724V router.")
  .demand('p')
  .alias('p', 'password')
  .describe('p', 'router password')
  .alias('a', 'address')
  .describe('a', 'router address')
  .default('a', 'speedport.ip')
  .argv


console.log('pw: ', argv.password);
const password = argv.password;
console.log('address: ', argv.address);

request('http://speedport.ip', function(error, response, body) {
  const challenge = body.match("var challenge = \"([a-zA-Z0-9]+)\"")[1]
  console.log("challenge: ", challenge);

  // const hashedPassword = fromBits(sjcl.hash.sha256.hash(challenge + ":" + password));
  const hashedPassword = crypto.createHash('sha256').update(challenge + ":" + password).digest('hex');
  console.log("hashedPassword:", hashedPassword);

  request.post({
    url: 'http://speedport.ip/data/Login.json',
    form: {
      csrf_token: 'nulltoken',
      password: hashedPassword,
      showpw: 0,
      challengev: challenge
    }
  }, function(error, response, body) {
    const responseCookies = response.headers['set-cookie'][0];
    const sessionId = responseCookies.match("SessionID_R3=([^;]+); ")[1]
    console.log("sessionId: ", sessionId);

    request({
      url: 'http://speedport.ip/html/content/config/problem_handling.html',
      method: 'GET',
      headers: {
        'Cookie': `SessionID_R3=${sessionId}`
      }
    }, function(error, response, body) {

      const csrfToken = body.match("var csrf_token = \"([^\"]+)\"")[1]
      console.log("csrfToken: ", csrfToken);

      request.post({
        url: 'http://speedport.ip/data/Reboot.json',
        headers: {
          'Cookie': `SessionID_R3=${sessionId};`
        },
        form: {
          reboot_device: 'true',
          csrf_token: csrfToken
        }
      }, function(error, response, body) {
        console.log(error);
        console.log(body);
      })
    });
  });

});
