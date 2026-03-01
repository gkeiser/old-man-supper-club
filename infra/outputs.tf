output "lambda_function_url" {
  description = "The public URL for the Supper Club web app"
  value       = aws_lambda_function_url.web_app_url.function_url
}
