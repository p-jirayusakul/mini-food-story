CREATE TABLE "md_table_statuses" (
                                     "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                     "code" varchar(15) UNIQUE NOT NULL,
                                     "name" varchar(100) UNIQUE NOT NULL,
                                     "name_en" varchar(100) UNIQUE NOT NULL,
                                     "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                                     "updated_at" timestamp
);

CREATE TABLE "md_categories" (
                                 "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                 "name" varchar(100) UNIQUE NOT NULL,
                                 "name_en" varchar(100) UNIQUE NOT NULL,
                                 "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                                 "updated_at" timestamp
);

CREATE TABLE "md_order_statuses" (
                                     "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                                     "code" varchar(15) UNIQUE NOT NULL,
                                     "name" varchar(100) UNIQUE NOT NULL,
                                     "name_en" varchar(100) UNIQUE NOT NULL,
                                     "sort_order" int UNIQUE NOT NULL,
                                     "is_final" bool NOT NULL DEFAULT false,
                                     "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                                     "updated_at" timestamp
);

CREATE TABLE "tables" (
                          "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                          "table_number" int UNIQUE NOT NULL,
                          "status_id" bigint NOT NULL,
                          "seats" int NOT NULL DEFAULT 0,
                          "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                          "updated_at" timestamp
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
                            "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                            "updated_at" timestamp
);

CREATE TABLE "orders" (
                          "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                          "session_id" uuid,
                          "table_id" bigint NOT NULL,
                          "status_id" bigint NOT NULL,
                          "total_amount" numeric(10,2) NOT NULL DEFAULT 0,
                          "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                          "updated_at" timestamp
);

CREATE TABLE "order_items" (
                               "id" bigint UNIQUE PRIMARY KEY NOT NULL,
                               "order_id" bigint NOT NULL,
                               "product_id" bigint NOT NULL,
                               "status_id" bigint NOT NULL,
                               "product_name" varchar(255) NOT NULL,
                               "product_name_en" varchar(255) NOT NULL,
                               "price" numeric(10,2) NOT NULL DEFAULT 0,
                               "quantity" int NOT NULL DEFAULT 1,
                               "note" text,
                               "prepared_at" timestamp,
                               "created_at" timestamp NOT NULL DEFAULT 'NOW()',
                               "updated_at" timestamp
);

CREATE INDEX ON "md_table_statuses" ("id");

CREATE INDEX ON "md_table_statuses" ("code");

CREATE INDEX ON "md_categories" ("id");

CREATE INDEX ON "md_order_statuses" ("id");

CREATE INDEX ON "md_order_statuses" ("code");

CREATE INDEX ON "tables" ("id");

CREATE INDEX ON "tables" ("table_number");

CREATE INDEX ON "tables" ("status_id");

CREATE INDEX ON "products" ("id");

CREATE INDEX ON "products" ("name");

CREATE INDEX ON "products" ("name_en");

CREATE UNIQUE INDEX ON "orders" ("id");

CREATE INDEX "orders_table_status" ON "orders" ("table_id", "status_id");

CREATE INDEX ON "orders" ("status_id");

CREATE INDEX ON "orders" ("created_at");

CREATE INDEX ON "order_items" ("id");

CREATE INDEX "order_items_order_status" ON "order_items" ("order_id", "status_id");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("status_id");

COMMENT ON COLUMN "order_items"."prepared_at" IS 'เวลาที่ทำอาหารเสร็จ';

ALTER TABLE "tables" ADD FOREIGN KEY ("status_id") REFERENCES "md_table_statuses" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("categories") REFERENCES "md_categories" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("status_id") REFERENCES "md_order_statuses" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("table_id") REFERENCES "tables" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("status_id") REFERENCES "md_order_statuses" ("id");
