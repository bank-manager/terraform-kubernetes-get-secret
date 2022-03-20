data "external" "secret" {
  program = ["${path.module}/bin/get-secret-value"]

  query = {
    namespace = var.namespace
    name      = var.name
    context   = var.context
  }
}
