package mongodb

const (
	optionOmitNil             = optionOmitNilWhere | optionOmitNilData
	optionOmitEmpty           = optionOmitEmptyWhere | optionOmitEmptyData
	optionOmitNilDataInternal = optionOmitNilData | optionOmitNilDataList // this option is used internally only for ForDao feature.
	optionOmitEmptyWhere      = 1 << iota                                 // 8
	optionOmitEmptyData                                                   // 16
	optionOmitNilWhere                                                    // 32
	optionOmitNilData                                                     // 64
	optionOmitNilDataList                                                 // 128
)

// OmitNilData sets optionOmitNilData option for the model, which automatically filers
// the Data parameters for `nil` values.
func (m *Model) OmitNilData() *Model {
	m.option = m.option | optionOmitNilData
	return m
}
