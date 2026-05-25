provider "aws" {
  region = "ap-south-1"
}

# ─────────────────────────────────────────────────────────────────────────────
# Secrets Manager Data Sources
# ─────────────────────────────────────────────────────────────────────────────
data "aws_secretsmanager_secret" "db_password" {
  name = "pebble/prod/db_password"
}

data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = data.aws_secretsmanager_secret.db_password.id
}

# ─────────────────────────────────────────────────────────────────────────────
# RDS PostgreSQL (Dev Tier)
# ─────────────────────────────────────────────────────────────────────────────
resource "aws_db_instance" "pebble_postgres" {
  identifier           = "pebble-dev-db"
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "15.4"
  instance_class       = "db.t3.micro"
  username             = "pebble_admin"
  password             = data.aws_secretsmanager_secret_version.db_password.secret_string
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true
  publicly_accessible  = false
  
  vpc_security_group_ids = [aws_security_group.db_sg.id]

  tags = {
    Environment = "dev"
    Project     = "pebble"
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# ElastiCache Redis (Dev Tier)
# ─────────────────────────────────────────────────────────────────────────────
resource "aws_elasticache_cluster" "pebble_redis" {
  cluster_id           = "pebble-dev-redis"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  engine_version       = "7.0"
  port                 = 6379

  security_group_ids = [aws_security_group.redis_sg.id]

  tags = {
    Environment = "dev"
    Project     = "pebble"
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# Security Groups (Allowing ECS Fargate to access DB/Redis)
# ─────────────────────────────────────────────────────────────────────────────
resource "aws_security_group" "db_sg" {
  name        = "pebble-db-sg"
  description = "Allow inbound traffic from ECS"

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # Adjust to your VPC CIDR
  }
}

resource "aws_security_group" "redis_sg" {
  name        = "pebble-redis-sg"
  description = "Allow inbound traffic from ECS"

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # Adjust to your VPC CIDR
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# IAM Role for ECS Task Execution (Secrets Access)
# ─────────────────────────────────────────────────────────────────────────────
resource "aws_iam_role_policy" "ecs_secrets_access" {
  name = "ecs-secrets-access"
  role = aws_iam_role.ecs_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Effect   = "Allow"
        Resource = [
          data.aws_secretsmanager_secret.db_password.arn,
          "arn:aws:secretsmanager:ap-south-1:*:secret:pebble/prod/*"
        ]
      }
    ]
  })
}

# Assuming ecs_execution_role is defined elsewhere or we define a minimal one here
resource "aws_iam_role" "ecs_execution_role" {
  name = "pebble-ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}
