#!/bin/bash
set -e

YAML_FILES=$(find ./api/http/oms -name "*.yaml")

INTERFACE_NAME="interfaces"
PROJECT_DIR="./api/http/oms"

for file in $YAML_FILES; do
  filename=$(basename -- "$file")
  filename="${filename%.*}"

  echo "make dir for $filename"
  if [ ! -d "./$INTERFACE_NAME/$filename" ]; then
    mkdir -p  "./$INTERFACE_NAME/$filename"
    cp "./$INTERFACE_NAME/http.go" "./$INTERFACE_NAME/$filename"
  fi

  echo "generate http_types $filename"
  oapi-codegen -generate types -o "$INTERFACE_NAME/$filename/http_types.gen.go" -package "$INTERFACE_NAME" "$PROJECT_DIR/$filename.yaml"
  echo "generate http_server $filename"
  oapi-codegen -generate chi-server -o "$INTERFACE_NAME/$filename/http_server.gen.go" -package "$INTERFACE_NAME" "$PROJECT_DIR/$filename.yaml"
  echo "generate http_client $filename"
  oapi-codegen -generate client -o "$INTERFACE_NAME/$filename/http_client.gen.go" -package "$INTERFACE_NAME" "$PROJECT_DIR/$filename.yaml"

done
