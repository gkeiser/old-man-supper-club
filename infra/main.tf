provider "aws" {
  region = var.aws_region
}

# 1. DynamoDB Table
resource "aws_dynamodb_table" "supper_club" {
  name         = var.dynamo_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PK"
  range_key    = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }
}

# 2. IAM Role for Lambda
resource "aws_iam_role" "lambda_exec" {
  name = "omsc_lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# 3. IAM Policy for CloudWatch and DynamoDB
resource "aws_iam_role_policy" "lambda_policy" {
  name = "omsc_lambda_policy"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:Scan",
          "dynamodb:Query",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem"
        ]
        Effect   = "Allow"
        Resource = aws_dynamodb_table.supper_club.arn
      }
    ]
  })
}

# 4. Lambda Function
resource "aws_lambda_function" "web_app" {
  function_name = "omsc-web-server"
  role          = aws_iam_role.lambda_exec.arn
  
  # We assume you will compile the Go binary as 'bootstrap' and zip it
  filename         = "../bootstrap.zip"
  source_code_hash = filebase64sha256("../bootstrap.zip")
  
  handler = "bootstrap"
  runtime = "provided.al2023"
  
  environment {
    variables = {
      DYNAMO_TABLE_NAME    = aws_dynamodb_table.supper_club.name
      GOOGLE_CLIENT_ID     = var.google_client_id
      GOOGLE_CLIENT_SECRET = var.google_client_secret
      GOOGLE_REDIRECT_URL  = var.google_redirect_url
      ALLOWED_EMAILS       = var.allowed_emails
      SESSION_SECRET       = var.session_secret
      ENV                  = "production"
    }
  }
}

# 5. Lambda Function URL (Public Endpoint)
resource "aws_lambda_function_url" "web_app_url" {
  function_name      = aws_lambda_function.web_app.function_name
  authorization_type = "NONE"
}
