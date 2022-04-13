package polygonws

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.Config{
    SortMapKeys:                   false,
    EscapeHTML:                    false,
    ValidateJsonRawMessage:        false,
    ObjectFieldMustBeSimpleString: true,
}.Froze()
