# Project Overview

## Description
โปรเจคนี้เป็นระบบที่ประกอบด้วยหลายบริการ (Microservices) ที่ทำงานร่วมกัน เหมาะสำหรับระบบจัดการร้านอาหาร ระบบนี้ถูกพัฒนาขึ้นมาด้วยภาษา Go และใช้ Docker ในการจัดการบริการต่าง ๆ

## โครงสร้างของโปรเจค
โปรเจคนี้มีโครงสร้างดังต่อไปนี้:
- `menu-service/`: บริการจัดการเมนูอาหาร
- `order-service/`: บริการจัดการคำสั่งซื้อ
- `table-service/`: บริการจัดการโต๊ะในร้านอาหาร
- `kitchen-service/`: บริการจัดการครัว
- `payment-service/`: บริการจัดการการชำระเงิน
- `shared/`: ไฟล์หรือโค้ดที่สามารถใช้งานร่วมกันได้ในหลายบริการ
- `pkg/`: ไฟล์ utils ต่าง ๆ ไม่ว่าจะเป็น middleware converter ต่าง ๆ

## การติดตั้งและเริ่มต้นใช้งาน
1. **Clone Repository**
   ```bash
   git clone https://github.com/p-jirayusakul/mini-food-sotry.git
   cd mini-food-sotry
   ```

2. **ตั้งค่าตัวแปรสภาพแวดล้อม**
   ตรวจสอบไฟล์ `.env` และตั้งค่าตามที่ต้องการ

3. **รันด้วย Docker Compose**
   หากคุณมี Docker และ Docker Compose ติดตั้งอยู่ ให้รันคำสั่งดังนี้:
    ```bash
   docker compose build
   ```
   ```bash
   docker compose up
   ```
   
4. **การรันคำสั่งอื่นๆ**
   ใช้ `Makefile` สำหรับการรันคำสั่งพิเศษ เช่น การสร้างหรือทดสอบ
    
   ###### generate sql with sqlc
   ```bash
   make sqlc
   ```
   ###### generate mockup for unit test
   ```bash
   make mock
   ```

## การพัฒนาและการมีส่วนร่วม
1. **Branching Model**
   - ใช้ `main` สำหรับ version ที่สามารถใช้งานได้
   - ใช้ branch feature เช่น `feature/<ชื่อฟีเจอร์>` สำหรับการพัฒนา

2. **เปิด Pull Request**
   - ตรวจสอบว่าฟีเจอร์ครบถ้วนและผ่านการทดสอบก่อนเปิด PR
   - เพิ่มคำอธิบายให้ชัดเจนถึงการเปลี่ยนแปลงที่ทำ

> **หมายเหตุ**: หากมีบริการที่จะเพิ่มในอนาคต คุณสามารถขยายโฟลเดอร์ในโครงสร้างโปรเจคตามความเหมาะสม