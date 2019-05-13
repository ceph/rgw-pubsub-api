#!/bin/bash
echo "// Code generated by $0; DO NOT EDIT." > secret_vars.go
echo "package main" >> secret_vars.go
echo "var (" >> secret_vars.go
echo "  googleVisionAPIKey  = \"$GOOGLE_VISION_API_KEY\"" >> secret_vars.go
echo "  s3AccessID          = \"$S3_ACCESS_KEY_ID\"" >> secret_vars.go
echo "  s3AccessKey         = \"$S3_SECRET_ACCESS_KEY\"" >> secret_vars.go
echo "  s3Endpoint          = \"$S3_HOSTNAME\"" >> secret_vars.go
echo ")" >> secret_vars.go
