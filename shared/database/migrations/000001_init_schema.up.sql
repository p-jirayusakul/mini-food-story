CREATE TYPE "table_session_status" AS ENUM (
    'active',
    'closed',
    'expired'
    );

CREATE TABLE "md_table_statuses" (
                                     "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                     "code" varchar(15) UNIQUE NOT NULL,
                                     "name" varchar(100) UNIQUE NOT NULL,
                                     "name_en" varchar(100) UNIQUE NOT NULL,
                                     "created_at" timestamptz NOT NULL DEFAULT NOW(),
                                     "updated_at" timestamptz
);

CREATE TABLE "md_categories" (
                                 "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                 "name" varchar(100) UNIQUE NOT NULL,
                                 "name_en" varchar(100) UNIQUE NOT NULL,
                                 icon_name  varchar(50),
                                 "sort_order" int UNIQUE NOT NULL,
                                 "created_at" timestamptz NOT NULL DEFAULT NOW(),
                                 "updated_at" timestamptz
);

CREATE TABLE "md_order_statuses" (
                                     "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                     "code" varchar(15) UNIQUE NOT NULL,
                                     "name" varchar(100) UNIQUE NOT NULL,
                                     "name_en" varchar(100) UNIQUE NOT NULL,
                                     "sort_order" int UNIQUE NOT NULL,
                                     "is_final" bool NOT NULL DEFAULT false,
                                     "created_at" timestamptz NOT NULL DEFAULT NOW(),
                                     "updated_at" timestamptz
);

CREATE TABLE "md_payment_methods" (
                                      "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                      "code" varchar(15) UNIQUE NOT NULL,
                                      "name" varchar(100) UNIQUE NOT NULL,
                                      "name_en" varchar(100) UNIQUE NOT NULL,
                                      "enable" bool NOT NULL DEFAULT false,
                                      "created_at" timestamptz NOT NULL DEFAULT NOW(),
                                      "updated_at" timestamptz
);

CREATE TABLE "md_payment_statuses" (
                                       "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                       "code" varchar(15) UNIQUE NOT NULL,
                                       "name" varchar(100) UNIQUE NOT NULL,
                                       "name_en" varchar(100) UNIQUE NOT NULL,
                                       "is_final" bool NOT NULL DEFAULT false,
                                       "created_at" timestamptz NOT NULL DEFAULT NOW(),
                                       "updated_at" timestamptz
);

CREATE TABLE "tables" (
                          "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                          "table_number" int UNIQUE NOT NULL,
                          "status_id" bigint NOT NULL,
                          "seats" int NOT NULL DEFAULT 0,
                          "created_at" timestamptz NOT NULL DEFAULT NOW(),
                          "updated_at" timestamptz
);

CREATE TABLE "table_session" (
                                 "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                 "table_id" bigint NOT NULL,
                                 "session_id" uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
                                 "number_of_people" int NOT NULL DEFAULT 1,
                                 "status" table_session_status,
                                 "started_at" timestamptz NOT NULL DEFAULT NOW(),
                                 "expire_at" timestamptz NOT NULL,
                                 "ended_at" timestamptz
);

CREATE TABLE "products" (
                            "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                            "name" varchar(255) UNIQUE NOT NULL,
                            "name_en" varchar(255) UNIQUE NOT NULL,
                            "categories" bigint NOT NULL,
                            "description" text,
                            "price" numeric(10,2) NOT NULL DEFAULT 0,
                            "is_available" bool NOT NULL DEFAULT false,
                            "image_url" text,
                            "created_at" timestamptz NOT NULL DEFAULT NOW(),
                            "updated_at" timestamptz
);

CREATE TABLE "orders" (
                          "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                          "order_number" varchar(50) UNIQUE NOT NULL,
                          "session_id" uuid,
                          "table_id" bigint NOT NULL,
                          "status_id" bigint NOT NULL,
                          "total_amount" numeric(10,2) NOT NULL DEFAULT 0,
                          "created_at" timestamptz NOT NULL DEFAULT NOW(),
                          "updated_at" timestamptz
);

CREATE TABLE "order_items" (
                               "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                               "order_id" bigint NOT NULL,
                               "product_id" bigint NOT NULL,
                               "status_id" bigint NOT NULL,
                               "product_name" varchar(255) NOT NULL,
                               "product_name_en" varchar(255) NOT NULL,
                               "product_image_url" text,
                               "price" numeric(10,2) NOT NULL DEFAULT 0,
                               "quantity" int NOT NULL DEFAULT 1,
                               "note" text,
                               "prepared_at" timestamptz,
                               "created_at" timestamptz NOT NULL DEFAULT NOW(),
                               "updated_at" timestamptz
);

CREATE TABLE "order_sequences" (
                                   "order_date" date PRIMARY KEY NOT NULL,
                                   "current_number" int NOT NULL
);

CREATE TABLE "payments" (
                            "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                            "order_id" bigint NOT NULL,
                            "amount" numeric(10,2) NOT NULL DEFAULT 0,
                            "method" bigint NOT NULL,
                            "status" bigint NOT NULL,
                            "paid_at" timestamptz,
                            "transaction_id" text UNIQUE NOT NULL,
                            "ref_code" varchar(150) UNIQUE NOT NULL,
                            "note" text,
                            "created_at" timestamptz NOT NULL DEFAULT NOW(),
                            "updated_at" timestamptz
);

CREATE INDEX ON "md_table_statuses" ("id");

CREATE INDEX ON "md_table_statuses" ("code");

CREATE INDEX ON "md_categories" ("id");

CREATE INDEX ON "md_order_statuses" ("id");

CREATE INDEX ON "md_order_statuses" ("code");

CREATE INDEX ON "md_payment_methods" ("id");

CREATE INDEX ON "md_payment_methods" ("code");

CREATE INDEX ON "md_payment_statuses" ("id");

CREATE INDEX ON "md_payment_statuses" ("code");

CREATE INDEX ON "tables" ("id");

CREATE INDEX ON "tables" ("table_number");

CREATE INDEX ON "tables" ("status_id");

CREATE INDEX ON "table_session" ("id");

CREATE INDEX ON "table_session" ("table_id");

CREATE INDEX ON "table_session" ("session_id");

CREATE INDEX ON "products" ("id");

CREATE INDEX ON "products" ("name");

CREATE INDEX ON "products" ("name_en");

CREATE INDEX ON "orders" ("id");

CREATE INDEX ON "orders" ("order_number");

CREATE INDEX ON "orders" ("status_id");

CREATE INDEX ON "orders" ("created_at");

CREATE INDEX ON "order_items" ("id");

CREATE INDEX "order_items_order_id" ON "order_items" ("order_id", "id");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "payments" ("id");

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("transaction_id");

CREATE INDEX "transaction_id_status" ON "payments" ("transaction_id", "status");

COMMENT ON COLUMN "order_items"."prepared_at" IS 'เวลาที่ทำอาหารเสร็จ';

ALTER TABLE "tables" ADD FOREIGN KEY ("status_id") REFERENCES "md_table_statuses" ("id");

ALTER TABLE "table_session" ADD FOREIGN KEY ("table_id") REFERENCES "tables" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("categories") REFERENCES "md_categories" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("status_id") REFERENCES "md_order_statuses" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("table_id") REFERENCES "tables" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("status_id") REFERENCES "md_order_statuses" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("method") REFERENCES "md_payment_methods" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("status") REFERENCES "md_payment_statuses" ("id");

-- Insert Master Data

insert into public.md_categories (id, name, name_en, icon_name, sort_order, created_at, updated_at)
values  (1921143886227443712, 'อาหาร', 'Food', 'soup', 1, '2025-10-10 05:34:36.643495 +00:00', null),
        (1921144050476388352, 'เครื่องดื่ม', 'Drink', 'glass-water', 2, '2025-10-10 05:34:36.643495 +00:00', null),
        (1921144250070732800, 'ขนม', 'Dessert', 'cake-slice', 3, '2025-10-10 05:34:36.643495 +00:00', null);

insert into public.md_order_statuses (id, code, name, name_en, sort_order, is_final, created_at, updated_at)
values  (1921868485739155456, 'PENDING', 'รอยืนยันออเดอร์', 'Pending', 1, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155457, 'CONFIRMED', 'ยืนยันออเดอร์', 'Confirmed', 2, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155458, 'PREPARING', 'กำลังเตรียมอาหาร', 'Preparing', 3, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155459, 'SERVED', 'เสิร์ฟอาหารแล้ว', 'Served', 4, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155460, 'WAITING_PAYMENT', 'รอชำระเงิน', 'Waiting for Payment', 5, false, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155461, 'COMPLETED', 'เสร็จสิ้น', 'Completed', 6, true, '2025-09-18 11:43:59.552632 +00:00', null),
        (1921868485739155462, 'CANCELLED', 'ยกเลิก', 'Cancelled', 7, true, '2025-09-18 11:43:59.552632 +00:00', null);

insert into public.md_table_statuses (id, code, name, name_en, created_at, updated_at)
values  (1919968486671519744, 'AVAILABLE', 'ว่าง', 'Available', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486843486208, 'RESERVED', 'ถูกจองล่วงหน้า', 'Reserved', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680514, 'WAITING_PAYMENT', 'รอชำระเงิน', 'Waiting for Payment', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680515, 'CLEANING', 'รอทำความสะอาด', 'Cleaning', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680516, 'DISABLED', 'ปิดการใช้งานชั่วคราว', 'Disabled', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680512, 'WAIT_ORDER', 'รอสั่ง', 'Waiting to Order', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680513, 'WAIT_SERVE', 'รอเสิร์ฟ', 'Waiting to be Served', '2025-09-18 11:43:59.552632 +00:00', null),
        (1919968486847680517, 'FOOD_SERVED', 'อาหารครบแล้ว', 'Food Served', '2025-09-18 16:10:29.071839 +00:00', null),
        (1919968486847680518, 'CALL_WAITER', 'เรียกพนักงาน', 'Call Waiter', '2025-09-18 16:10:29.071839 +00:00', null);


insert into public.products (id, name, name_en, categories, description, price, is_available, image_url, created_at, updated_at)
values  (1921822053405560832, 'ข้าวผัด', 'Fried rice', 1921143886227443712, 'lorem ipso', 60.00, true, 'https://images.unsplash.com/photo-1603133872878-684f208fb84b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921828287366041600, 'ข้าวผัดกระเพรา', 'Phat kaphrao', 1921143886227443712, 'lorem ipso', 70.50, true, 'https://images.unsplash.com/photo-1694499792070-48e64a00cf0a?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921821817723424768, 'ข้าวมันไก่', 'Chicken rice', 1921143886227443712, 'lorem ipso', 80.00, true, 'https://images.unsplash.com/photo-1749640566096-5d8098d452b4?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921822608437809152, 'แป๊ปซี่', 'Pepsi', 1921144050476388352, 'lorem ipso', 30.00, true, 'https://images.unsplash.com/photo-1651000877733-fe2c0a70b3cd?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null),
        (1921822481287483392, 'เค้กแครอท', 'Carrot cake', 1921144250070732800, 'lorem ipso', 120.00, true, 'https://images.unsplash.com/photo-1622926421334-6829deee4b4b?ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&q=80&w=500', '2025-10-10 05:34:36.692235 +00:00', null);

INSERT INTO public."tables" (id,table_number,status_id,seats,created_at,updated_at) VALUES
                                                                                        (1920153361642950656,5,1919968486671519744,4,NOW(),NULL),
                                                                                        (1919972141986484224,3,1919968486671519744,4,NOW(),NULL),
                                                                                        (1919996486741921792,4,1919968486671519744,5,NOW(),NULL),
                                                                                        (1919968785813475328,1,1919968486671519744,5,NOW(),NULL),
                                                                                        (1919971956241731584,2,1919968486671519744,3,NOW(),NULL);


INSERT INTO public.md_payment_methods (id,code,"name",name_en,"enable",created_at,updated_at) VALUES
                                                                                               (1923732004537372672,'CASH','เงินสด','Cash',true,NOW(),NULL),
                                                                                               (1923732004537372675,'PROMPTPAY','พร้อมเพย์','PromptPay',true,NOW(),NULL);
