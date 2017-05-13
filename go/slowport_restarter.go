package main

import (
  "os"
  "flag"
  "fmt"
  "time"
  "net/http"
  "io/ioutil"
  "regexp"
  "crypto/sha256"
  "encoding/hex"
  "bytes"
  "net/url"
)

var client = &http.Client{
  Timeout: time.Second * 10,
}

func main() {
  password, host := parseArgs()

  challenge := getChallengeToken(host)

  hashedPassword := hashPassword(challenge, password)

  sessionId := login(host, hashedPassword, challenge)

  csrfToken := getCsrfToken(host, sessionId)

  rebootRouter(host, sessionId, csrfToken)
}

func parseArgs() (password string, host string) {
  passwordPtr := flag.String("password", "", "the router password")
  hostPtr := flag.String("host", "speedport.ip", "the router address")
  flag.Parse()

  return *passwordPtr, *hostPtr
}

func getChallengeToken(host string) string {
  req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/", host), nil)
  resp, err := client.Do(req)
  if err != nil {
    os.Stderr.WriteString("slowport-restarter --- Error getting challenge token\n")
    os.Exit(1)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  var html = string(body);

  var challengeExp = regexp.MustCompile(`var challenge = \"(?P<challenge>[^\"]+)\"`)
  challenge := (*findPattern(html, challengeExp))["challenge"]
  return challenge
}

func hashPassword(challenge string, password string) string {
  hash := sha256.New()
  hash.Write([]byte(fmt.Sprintf("%s:%s", challenge, password)))
  hashedBytes := hash.Sum(nil)
  hashedPassword := hex.EncodeToString(hashedBytes)
  return hashedPassword
}

func login(host string, hashedPassword string, challenge string) string {
  data := url.Values{}
  data.Set("csrf_token", "nulltoken")
  data.Add("password", hashedPassword)
  data.Add("showpw", "0")
  data.Add("challengev", challenge)
  req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/data/Login.json", host), bytes.NewBufferString(data.Encode()))
  resp, err := client.Do(req)
  if err != nil {
  	os.Stderr.WriteString("slowport-restarter --- Error logging in\n")
    os.Exit(2)
  }

  // Extract session id
  var sessionIdExp = regexp.MustCompile(`SessionID_R3=(?P<sessionid>[^;]+);`)
  sessionId := (*findPattern(resp.Header.Get("set-cookie"), sessionIdExp))["sessionid"]
  if sessionId == "" {
    os.Stderr.WriteString("slowport-restarter --- Error logging in, sessionId not found\n")
    os.Exit(3)
  }
  return sessionId
}

func getCsrfToken(host string, sessionId string) string {
  req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/html/content/config/problem_handling.html", host), nil)
  req.Header.Add("Cookie", fmt.Sprintf("SessionID_R3=%s", sessionId))
  resp, err := client.Do(req)
  if err != nil {
  	os.Stderr.WriteString("slowport-restarter --- Error getting CSRF token for reboot\n")
    os.Exit(4)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  var html = string(body);

  var csrfTokenExp = regexp.MustCompile(`var csrf_token = \"(?P<token>[^\"]+)\"`)
  csrfToken := (*findPattern(html, csrfTokenExp))["token"]
  if csrfToken == "" {
    os.Stderr.WriteString("slowport-restarter --- Error getting CSRF token for reboot, not found in HTML\n")
    os.Exit(5)
  }
  return csrfToken
}

func rebootRouter(host string, sessionId string, csrfToken string) {
  fmt.Println("Rebooting")
  fmt.Println("host:", host)
  fmt.Println("sessionId:", sessionId)
  fmt.Println("csrfToken:", csrfToken)

  data := url.Values{}
  data.Set("csrf_token", csrfToken)
  data.Add("reboot_device", "true")
  req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/data/Reboot.json", host), bytes.NewBufferString(data.Encode()))
  req.Header.Add("Cookie", fmt.Sprintf("SessionID_R3=%s", sessionId))
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("error:", err)
  	os.Stderr.WriteString("slowport-restarter --- Error sending reboot request\n")
    os.Exit(6)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  var html = string(body);

  var rebootStatusExp = regexp.MustCompile(`\"varvalue\":\"(?P<status>[^\"]+)\"`)
  rebootStatus := (*findPattern(html, rebootStatusExp))["status"]
  if rebootStatus != "ok" {
    os.Stderr.WriteString("slowport-restarter --- Unexpected reboot response from router\n")
    os.Exit(7)
  }
}

func findPattern(text string, regex *regexp.Regexp) *map[string]string {
  match  := regex.FindStringSubmatch(text)
  result := make(map[string]string)
  if match != nil {
    for i, name := range regex.SubexpNames() {
      if i != 0 { result[name] = match[i] }
    }
    return &result
  } else  {
    return &result
  }
}
