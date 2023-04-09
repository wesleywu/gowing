package mongodb

var (
	createdFiledNames = "createdAt" // Default filed names of table for automatic-filled created datetime.
	updatedFiledNames = "updatedAt" // Default filed names of table for automatic-filled updated datetime.
	deletedFiledNames = "deletedAt" // Default filed names of table for automatic-filled deleted datetime.
)

// getSoftFieldNameUpdate checks and returns the field name for record updating time.
// If there's no field name for storing updating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameUpdated() (field string) {
	// It checks whether this feature disabled.
	if m.db.GetConfig().TimeMaintainDisabled {
		return ""
	}
	config := m.db.GetConfig()
	if config.UpdatedAt != "" {
		return config.UpdatedAt
	}
	return updatedFiledNames
}
