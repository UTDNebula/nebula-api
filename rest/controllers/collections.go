package controllers

import "github.com/UTDNebula/nebula-api/rest/configs"

// InitializeCollections creates collection handles after the database connects.
func InitializeCollections() {
	astraCollection = configs.GetCollection("astra")
	DAGCollection = configs.GetCollection("DAG")
	budgetCollection = configs.GetCollection("budgets")
	cometCalendarCollection = configs.GetCollection("cometCalendar")
	courseCollection = configs.GetCollection("courses")
	discountCollection = configs.GetCollection("discounts")
	eventsCollection = configs.GetCollection("events")
	mazevoCollection = configs.GetCollection("mazevo")
	professorCollection = configs.GetCollection("professors")
	buildingCollection = configs.GetCollection("rooms")
	sectionCollection = configs.GetCollection("sections")
}
