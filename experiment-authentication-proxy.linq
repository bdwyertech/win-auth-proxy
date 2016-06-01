<Query Kind="Expression">
  <Namespace>System.Net</Namespace>
</Query>

new WebClient{ Proxy = new WebProxy("http://127.0.0.1:8080") }.DownloadString("http://www.google.com")
// not working yet!
// new WebClient{ Proxy = new WebProxy("http://127.0.0.1:8080") }.DownloadString("https://www.google.com")