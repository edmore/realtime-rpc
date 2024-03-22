output "realtime-rpc_ecr_repository" {
  description = "realtime-rpc ECR repository"

  value = aws_ecr_repository.realtime-rpc.repository_url
}