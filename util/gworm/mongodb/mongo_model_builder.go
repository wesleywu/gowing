package mongodb

// WhereBuilder holds multiple where conditions in a group.
type WhereBuilder struct {
	model       *Model        // A WhereBuilder should be bound to certain Model.
	whereHolder []WhereHolder // Condition strings for where operation.
}

// WhereHolder is the holder for where condition preparing.
type WhereHolder struct {
	Type     string        // Type of this holder.
	Operator int           // Operator for this holder.
	Where    interface{}   // Where parameter, which can commonly be type of string/map/struct.
	Args     []interface{} // Arguments for where parameter.
	Prefix   string        // Field prefix, eg: "user.", "order.".
}

// Builder creates and returns a WhereBuilder.
func (m *Model) Builder() *WhereBuilder {
	b := &WhereBuilder{
		model:       m,
		whereHolder: make([]WhereHolder, 0),
	}
	return b
}

// getBuilder creates and returns a cloned WhereBuilder of current WhereBuilder if `safe` is true,
// or else it returns the current WhereBuilder.
func (b *WhereBuilder) getBuilder() *WhereBuilder {
	return b.Clone()
}

// Clone clones and returns a WhereBuilder that is a copy of current one.
func (b *WhereBuilder) Clone() *WhereBuilder {
	newBuilder := b.model.Builder()
	newBuilder.whereHolder = make([]WhereHolder, len(b.whereHolder))
	copy(newBuilder.whereHolder, b.whereHolder)
	return newBuilder
}
