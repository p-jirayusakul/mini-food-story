package mockshared

import shareModel "food-story/shared/model"

func MockOrderItems() []shareModel.OrderItems {
	return []shareModel.OrderItems{
		{
			ID:            int64(1),
			OrderID:       int64(1),
			OrderNumber:   "FS-20250520-0001",
			ProductID:     int64(1),
			StatusID:      int64(1),
			TableNumber:   1,
			StatusName:    "กำลังเตรียมอาหาร",
			StatusNameEN:  "Pending",
			StatusCode:    "PENDING",
			ProductName:   "Test Product",
			ProductNameEN: "Test Product",
			Price:         100.00,
			Quantity:      1,
			Note:          nil,
			CreatedAt:     "2025-05-20T15:13:15+07:00",
		},
		{
			ID:            int64(2),
			OrderID:       int64(1),
			OrderNumber:   "FS-20250520-0001",
			ProductID:     int64(1),
			StatusID:      int64(1),
			TableNumber:   1,
			StatusName:    "กำลังเตรียมอาหาร",
			StatusNameEN:  "Pending",
			StatusCode:    "PENDING",
			ProductName:   "Test Product 2",
			ProductNameEN: "Test Product 2",
			Price:         200.00,
			Quantity:      1,
			Note:          nil,
			CreatedAt:     "2025-05-20T16:00:00+07:00",
		},
	}
}

func MockOrderItemsPt() []*shareModel.OrderItems {
	var result []*shareModel.OrderItems
	for _, item := range MockOrderItems() {
		result = append(result, &item)
	}
	return result
}
