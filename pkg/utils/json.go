package utils

import (
    "encoding/json"
)

func MustMarshal(v interface{}) []byte {
    b, _ := json.Marshal(v)
    return b
}

func MustUnmarshal(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}


