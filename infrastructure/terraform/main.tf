provider "aws" {
  region = "ap-south-1"
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
  password             = "CHANGE_ME_IN_PROD"
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
