package mongodb

type Option func(*mongodb)

// func Upsert(isUpsert bool) Option {
// 	return func(m *mongodb) {
// 		m.upsert = isUpsert
// 	}
// }

// func MaxPoolSize(max int) Option {
// 	return func(m *mongodb) {
// 		m.maxPoolSize = max
// 	}
// }
