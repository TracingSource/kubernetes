package validate

// SchemaValidatorOptions defines optional rules for schema validation
type SchemaValidatorOptions struct {
	EnableObjectArrayTypeCheck    bool
	EnableArrayMustHaveItemsCheck bool
}

// Option sets optional rules for schema validation
type Option func(*SchemaValidatorOptions)

// EnableObjectArrayTypeCheck activates the swagger rule: an items must be in type: array
func EnableObjectArrayTypeCheck(enable bool) Option {
	return func(svo *SchemaValidatorOptions) {
		svo.EnableObjectArrayTypeCheck = enable
	}
}

// EnableArrayMustHaveItemsCheck activates the swagger rule: an array must have items defined
func EnableArrayMustHaveItemsCheck(enable bool) Option {
	return func(svo *SchemaValidatorOptions) {
		svo.EnableArrayMustHaveItemsCheck = enable
	}
}

// SwaggerSchema activates swagger schema validation rules
func SwaggerSchema(enable bool) Option {
	return func(svo *SchemaValidatorOptions) {
		svo.EnableObjectArrayTypeCheck = enable
		svo.EnableArrayMustHaveItemsCheck = enable
	}
}

// Options returns current options
func (svo SchemaValidatorOptions) Options() []Option {
	return []Option{
		EnableObjectArrayTypeCheck(svo.EnableObjectArrayTypeCheck),
		EnableArrayMustHaveItemsCheck(svo.EnableArrayMustHaveItemsCheck),
	}
}
