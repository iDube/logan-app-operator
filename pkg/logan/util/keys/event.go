package keys

const (
	// Boot event reason list

	// CreatedDeployment is the event reason for created deployment
	CreatedDeployment = "CreatedDeployment"
	// CreatedStatefulSet is the event reason for created statefulSet
	CreatedStatefulSet = "CreatedStatefulSet"
	// FailedCreateDeployment is the failed event reason for created deployment
	FailedCreateDeployment = "FailedCreateDeployment"
	// FailedCreateStatefulSet is the failed event reason for created statefulSet
	FailedCreateStatefulSet = "FailedCreateStatefulSet"
	// UpdatedDeployment is the event reason for updated deployment
	UpdatedDeployment = "UpdatedDeployment"
	// UpdatedStatefulSet is the event reason for updated statefulSet
	UpdatedStatefulSet = "UpdatedStatefulSet"
	// FailedUpdateDeployment is the failed event reason for updated deployment
	FailedUpdateDeployment = "FailedUpdateDeployment"
	// FailedUpdateStatefulSet is the failed event reason for updated statefulSet
	FailedUpdateStatefulSet = "FailedUpdateStatefulSet"
	// FailedGetDeployment is the failed event reason for got deployment
	FailedGetDeployment = "FailedGetDeployment"
	// FailedGetStatefulSet is the failed event reason for got statefulSet
	FailedGetStatefulSet = "FailedGetStatefulSet"

	// CreatedService is the event reason for created service
	CreatedService = "CreatedService"
	// FailedCreateService is the failed event reason for created service
	FailedCreateService = "FailedCreateService"
	// UpdatedService is the event reason for updated service
	UpdatedService = "UpdatedService"
	// FailedUpdateService is the failed event reason for updated service
	FailedUpdateService = "FailedUpdateService"
	// DeletedService is the event reason for deleted service
	DeletedService = "DeletedService"
	// FailedDeleteService is the failed event reason for deleted service
	FailedDeleteService = "FailedDeleteService"
	// FailedGetService is the failed event reason for got service
	FailedGetService = "FailedGetService"

	// UpdatedBootDefaulters is the event reason for updated boot defaulters
	UpdatedBootDefaulters = "UpdatedBootDefaulters"
	// FailedUpdateBootDefaulters is the failed event reason for updated boot defaulters
	FailedUpdateBootDefaulters = "FailedUpdateBootDefaulters"
	// UpdatedBootMeta is the event reason for updated boot meta
	UpdatedBootMeta = "UpdatedBootMeta"
	// FailedUpdateBootMeta is the failed event reason for updated boot meta
	FailedUpdateBootMeta = "FailedUpdateBootMeta"
)
