package mongodb

var (
	createdFiledNames = "createAt" // Default filed names of table for automatic-filled created datetime.
	updatedFiledNames = "updateAt" // Default filed names of table for automatic-filled updated datetime.
	deletedFiledNames = "deleteAt" // Default filed names of table for automatic-filled deleted datetime.
)

// getSoftFieldNameUpdate checks and returns the field name for record updating time.
// If there's no field name for storing updating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameUpdated() (field string) {
	return updatedFiledNames
}
