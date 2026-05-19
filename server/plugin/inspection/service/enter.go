package service

var Service = new(service)

type service struct {
	InspectionItem  inspectionItem
	Template        templateSvc
	ProductionOrder productionOrderSvc
	WorkOrder       workOrderSvc
}
