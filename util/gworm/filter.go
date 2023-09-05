package gworm

//func (f FilterRequest) AddAnypbValueCondition(_ context.Context, property string, value *anypb.Any, operator gwtypes.OperatorType, multi gwtypes.MultiType, wildcard gwtypes.WildcardType) FilterRequest {
//	condition := &gwtypes.Condition{
//		Operator: operator,
//		Multi:    multi,
//		Wildcard: wildcard,
//		Value:    value,
//	}
//	f.PropertyFilters = append(f.PropertyFilters, PropertyFilter{
//		Property:       property,
//		PropertyFilter: condition,
//	})
//	return f
//}
//
//func (f FilterRequest) AddAnyValueCondition(ctx context.Context, property string, value interface{}, operator gwtypes.OperatorType, multi gwtypes.MultiType, wildcard gwtypes.WildcardType, tag reflect.StructTag) FilterRequest {
//	if value == nil {
//		if operator == gwtypes.OperatorType_Null || operator == gwtypes.OperatorType_NotNull {
//			condition := &gwtypes.Condition{
//				Operator: operator,
//				Multi:    gwtypes.MultiType_Exact,
//				Wildcard: gwtypes.WildcardType_None,
//				Value:    nil,
//			}
//			f.PropertyFilters = append(f.PropertyFilters, PropertyFilter{
//				Property:       property,
//				PropertyFilter: condition,
//			})
//		}
//		return f
//	}
//	valueAny := gwwrapper.AnyInterface(ctx, value, tag)
//	return f.AddAnypbValueCondition(ctx, property, valueAny, operator, multi, wildcard)
//}
