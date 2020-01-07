log_level = "aa"

entrypoint "http" {
  bind = ["127.0.0.1:80"]
}

entrypoint "https" {
  bind = ["*:443"]
}
