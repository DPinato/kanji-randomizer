provider "aws" {
    region = "eu-central-1"
}

variable "function_name" {
    default = "kanji-randomizer"
}

variable "filename" {
    default = "kanji-randomizer.zip"
}

resource "aws_lambda_function" "golambda" {
    function_name = "${var.function_name}"
    filename      = "${var.filename}"
    role          = aws_iam_role.iam_for_kanji_randomizer.arn
    runtime       = "go1.x"
    handler       = "${var.function_name}"
    source_code_hash = filebase64sha256(var.filename)
}

resource "aws_iam_role" "iam_for_kanji_randomizer" {
    name = "iam_for_kanji_randomizer"
    assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
            "Service": "lambda.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

