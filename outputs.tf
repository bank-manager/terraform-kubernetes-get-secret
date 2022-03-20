output "result" {
  value = data.external.secret.result
  sensitive = true
}
