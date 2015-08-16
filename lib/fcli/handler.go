package fcli
import (
    "reflect"
    "errors"
)

type Handler struct {
    Format string
    Ptr interface{}
}

func NewHandler(format string, f_ptr interface{}) (h *Handler, err error) {
    h = &Handler{format, f_ptr}
    f_type := reflect.TypeOf(h.Ptr)
    if f_type.Kind() != reflect.Func { return nil, errors.New("State handler must be a function") }
    if f_type.NumOut() == 1 {
        ret_type := f_type.Out(0)
        typeOfResult := reflect.TypeOf((*Result)(nil)).Elem()
        if ret_type == typeOfResult { return }
    }
    return nil, errors.New("State handler must return a Result type")
}

func (h *Handler) Call(args []interface{}) Result {
    arg_values := make([]reflect.Value, len(args))
    for i, arg := range args {
        arg_values[i] = reflect.ValueOf(arg)
    }
    call_result := reflect.ValueOf(h.Ptr).Call(arg_values)
    if len(call_result) != 1 { panic("Invalid number of returned values") }
    return_value := call_result[0].Interface()

    // If the returned value is nil, pass it on
    if return_value == nil { return nil }

    // If the returned value is an error, pass it on
    if result, ok := return_value.(error); ok {
        return ErrorResult{result}
    }

    // Else expect Result type
    if result, ok := return_value.(Result); ok {
        return result
    } else {
        panic("Failed to cast handler return value to Result")
    }
}
