#!/bin/bash

ENV_FILE=".env"

# รับ key และ value จาก args
KEY=$1
VALUE=$2

# เช็คว่ามี key นี้อยู่ในไฟล์แล้วหรือยัง
if grep -q "^${KEY}=" "$ENV_FILE"; then
  # ถ้ามีแล้ว: แทนที่ค่าด้วยค่าใหม่
  sed -i "s|^${KEY}=.*|${KEY}=${VALUE}|" "$ENV_FILE"
else
  # ถ้ายังไม่มี: เพิ่มบรรทัดใหม่
  echo "${KEY}=${VALUE}" >> "$ENV_FILE"
fi