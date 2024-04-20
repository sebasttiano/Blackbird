package service

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
)

func ExampleService_GetValue() {
	repo := repository.NewMemStorage()
	service := NewService(&ServiceSettings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	service.SetValue(ctx, "test_gauge", "gauge", "3.33")
	value1, _ := service.GetValue(ctx, "test_gauge", "gauge")
	fmt.Println(value1)

	service.SetValue(ctx, "test_counter", "counter", "50")
	value2, _ := service.GetValue(ctx, "test_counter", "counter")
	fmt.Println(value2)

	service.SetValue(ctx, "test_counter2", "counter", "150")
	value3, _ := service.GetValue(ctx, "test_counter2", "counter")
	fmt.Println(value3)

	if err := service.SetValue(ctx, "test_counter_bad", "counter", "asd"); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 3.33
	// 50
	// 150
	// strconv.ParseInt: parsing "asd": invalid syntax

}

func ExampleService_GetModelValue() {
	repo := repository.NewMemStorage()
	service := NewService(&ServiceSettings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	valueFloat := 123.456
	valueCounter := int64(100)

	m := []*models.Metrics{
		{ID: "test_gauge", MType: "gauge", Value: &valueFloat},
		{ID: "test_counter", MType: "counter", Delta: &valueCounter}}
	service.SetModelValue(ctx, m)

	m2 := models.Metrics{ID: "test_gauge", MType: "gauge"}
	m3 := models.Metrics{ID: "test_counter", MType: "counter"}

	service.GetModelValue(ctx, &m2)
	service.GetModelValue(ctx, &m3)

	fmt.Println(*m2.Value)
	fmt.Println(*m3.Delta)

	// Output:
	// 123.456
	// 100
}

func ExampleService_GetAllValues() {
	repo := repository.NewMemStorage()
	service := NewService(&ServiceSettings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	valueFloat := 31.36
	valueCounter := int64(300)
	valueCounter2 := int64(23)

	m := []*models.Metrics{
		{ID: "test_gauge", MType: "gauge", Value: &valueFloat},
		{ID: "test_counter", MType: "counter", Delta: &valueCounter},
		{ID: "test_counter2", MType: "counter", Delta: &valueCounter2}}

	service.SetModelValue(ctx, m)

	sm := service.GetAllValues(ctx)
	for _, v := range sm.Counter {
		fmt.Println(v.Value)
	}
	for _, v := range sm.Gauge {
		fmt.Println(v.Value)
	}

	// Output:
	// 300
	// 23
	// 31.36
}
