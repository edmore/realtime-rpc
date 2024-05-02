// ECS Cluster
resource "aws_kms_key" "realtime_cluster" {
  description             = "ecs_cluster_kms_key"
  deletion_window_in_days = 7
}

resource "aws_cloudwatch_log_group" "realtime_cluster" {
  name = "ecs-cluster-log"
}


resource "aws_service_discovery_http_namespace" "realtime_service_namespace" {
  name        = "realtime-namespace"
  description = "service namespace"
}

resource "aws_ecs_cluster" "realtime_cluster" {
  name = "realtime_cluster"

  configuration {
    execute_command_configuration {
      kms_key_id = aws_kms_key.realtime_cluster.arn
      logging    = "OVERRIDE"

      log_configuration {
        cloud_watch_encryption_enabled = true
        cloud_watch_log_group_name     = aws_cloudwatch_log_group.realtime_cluster.name
      }
    }
  }

  service_connect_defaults {
    namespace = aws_service_discovery_http_namespace.realtime_service_namespace.arn
  }
}

// ECS Task definition - realtime-rpc
resource "aws_ecs_task_definition" "realtime-rpc" {
  family                = "realtime-rpc"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 2048
  memory                   = 4096
  task_role_arn      = aws_iam_role.task_role_for_ecs_task.arn
  execution_role_arn = aws_iam_role.execution_role_for_ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "realtime-rpc"
      image     = aws_ecr_repository.realtime-rpc.repository_url
      environment: [
      ],
      essential = true
      portMappings = [
        {
          name = "grpc-port"
          containerPort = 50051
          hostPort      = 50051
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group = "/ecs/realtime-rpc/logs"
          awslogs-region = var.region
          awslogs-stream-prefix = "ecs"
          awslogs-create-group = "true"
        }
      }
    }
  ])
}

resource "aws_ecs_service" "realtime-rpc" {
  name            = "realtime-rpc"
  cluster         = aws_ecs_cluster.realtime_cluster.id
  task_definition = aws_ecs_task_definition.realtime-rpc.arn
  launch_type = "FARGATE"
  desired_count = 1

  network_configuration {
    subnets = local.subnet_ids_list
    assign_public_ip = true
    security_groups = [aws_default_security_group.default.id]
  }

  service_connect_configuration {
    enabled = true
    namespace = aws_service_discovery_http_namespace.realtime_service_namespace.arn
    service {
      client_alias {
        dns_name = "grpc-service"
        port = "50051"
      }
      discovery_name = "realtime-grpc-service"
      port_name = "grpc-port"
    }
  }
}


// ECS Task definition - realtime-rpc-client
resource "aws_ecs_task_definition" "realtime-rpc-client" {
  family                = "realtime-rpc-client"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 2048
  memory                   = 4096
  task_role_arn      = aws_iam_role.task_role_for_ecs_task.arn
  execution_role_arn = aws_iam_role.execution_role_for_ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "realtime-rpc-client"
      image     = aws_ecr_repository.realtime-rpc-client.repository_url
      environment: [
        {"name": "SERVER_ENDPOINT", "value": "${var.service_endpoint}"}
      ],
      essential = true
      portMappings = [
        {
          containerPort = 8081
          hostPort      = 8081
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group = "/ecs/realtime-rpc-client/logs"
          awslogs-region = var.region
          awslogs-stream-prefix = "ecs"
          awslogs-create-group = "true"
        }
      }
    }
  ])
}