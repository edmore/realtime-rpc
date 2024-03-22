resource "aws_ecr_repository" "realtime-rpc" {
  name                 = "realtime-rpc"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false # consider implications of setting to true
  }
}