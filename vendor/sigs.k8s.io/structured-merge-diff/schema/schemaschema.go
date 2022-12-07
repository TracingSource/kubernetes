package schema

// SchemaSchemaYAML is a schema against which you can validate other schemas.
// It will validate itself. It can be unmarshalled into a Schema type.
var SchemaSchemaYAML = `types:
- name: schema
  map:
    fields:
      - name: types
        type:
          list:
            elementRelationship: associative
            elementType:
              namedType: typeDef
            keys:
            - name
- name: typeDef
  map:
    fields:
    - name: name
      type:
        scalar: string
    - name: scalar
      type:
        scalar: string
    - name: map
      type:
        namedType: map
    - name: list
      type:
        namedType: list
    - name: untyped
      type:
        namedType: untyped
- name: typeRef
  map:
    fields:
    - name: namedType
      type:
        scalar: string
    - name: scalar
      type:
        scalar: string
    - name: map
      type:
        namedType: map
    - name: list
      type:
        namedType: list
    - name: untyped
      type:
        namedType: untyped
- name: scalar
  scalar: string
- name: map
  map:
    fields:
    - name: fields
      type:
        list:
          elementType:
            namedType: structField
          elementRelationship: associative
          keys: [ "name" ]
    - name: unions
      type:
        list:
          elementType:
            namedType: union
          elementRelationship: atomic
    - name: elementType
      type:
        namedType: typeRef
    - name: elementRelationship
      type:
        scalar: string
- name: unionField
  map:
    fields:
    - name: fieldName
      type:
        scalar: string
    - name: discriminatorValue
      type:
        scalar: string
- name: union
  map:
    fields:
    - name: discriminator
      type:
        scalar: string
    - name: deduceInvalidDiscriminator
      type:
        scalar: bool
    - name: fields
      type:
        list:
          elementRelationship: associative
          elementType:
            namedType: unionField
          keys:
          - fieldName
- name: structField
  map:
    fields:
    - name: name
      type:
        scalar: string
    - name: type
      type:
        namedType: typeRef
- name: list
  map:
    fields:
    - name: elementType
      type:
        namedType: typeRef
    - name: elementRelationship
      type:
        scalar: string
    - name: keys
      type:
        list:
          elementType:
            scalar: string
- name: untyped
  map:
    fields:
    - name: elementRelationship
      type:
        scalar: string
`
