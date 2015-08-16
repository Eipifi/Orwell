package fcli
import (
    "fmt"
    "bufio"
    "os"
    "strings"
)

type FSM struct {
    prompt string
    parsers map[string] InputParser
    handlers map[string] []*Handler
}

const EMPTY_FORMAT = ""
const VAR_PREFIX = "$"

func NewFSM(prompt string) *FSM {
    fsm := &FSM{}
    fsm.prompt = prompt
    fsm.parsers = make(map[string] InputParser)
    fsm.handlers = make(map[string] []*Handler)

    fsm.Parser("str", StringParser)
    fsm.Parser("uint64", Uint64Parser)
    fsm.Parser("U256", U256Parser)
    return fsm
}

func (f *FSM) Parser(name string, parser InputParser) {
    f.parsers[name] = parser
}

func (f *FSM) On(state_name string, format string, func_ptr interface{}) {
    new_handler, err := NewHandler(format, func_ptr)
    if err != nil {
        panic(err)
    }
     if handlers, ok := f.handlers[state_name]; ok {
         f.handlers[state_name] = append(handlers, new_handler)
     } else {
         f.handlers[state_name] = []*Handler{new_handler}
     }
}

func (f *FSM) Run(state_name string) error {
    for {
        handlers, ok := f.handlers[state_name]
        if ! ok {
            panic("Unknown state: " + state_name)
        }

        // Check if any handler specifies an empty format
        line := "Han shot first"
        for _, handler := range handlers {
            if handler.Format == EMPTY_FORMAT {
                line = EMPTY_FORMAT
            }
        }
        if line != EMPTY_FORMAT {
            line = f.readLine()
        }

        var result Result = nil
        found := false
        for _, handler := range handlers {
            if args, ok := f.parseFormat(handler.Format, line); ok {
                found = true
                result = handler.Call(args)
            }
        }

        if found {
            switch res := result.(type) {
                case NextStateResult: state_name = string(res)
                case ExitResult: return res.err
                case ErrorResult: if res.err != nil {
                    fmt.Printf("Error: %v \n", res.err)
                }
            }
        } else {
            if line != "" {
                fmt.Printf("Unexpected line \"%s\"\n", line)
            }
        }
    }
}

func (f *FSM) readLine() string {
    for {
        fmt.Printf(f.prompt)
        line, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err == nil {
            return strings.TrimSuffix(line, "\n")
        }
        fmt.Printf("Error: %v \n", err)
    }
}

func (f *FSM) parseFormat(format string, line string) ([]interface{}, bool) {
    values := make([]interface{}, 0)
    if format == EMPTY_FORMAT { return values, true }
    format_fields := strings.Fields(format)
    line_fields := strings.Fields(line)
    if len(line_fields) == 0 {
        line_fields = []string{""}
    }
    if len(format_fields) != len(line_fields) {
        return nil, false
    }
    for i := 0; i < len(format_fields); i += 1 {
        if strings.HasPrefix(format_fields[i], VAR_PREFIX) {
            parser_type := strings.TrimPrefix(format_fields[i], VAR_PREFIX)
            if parser, ok := f.parsers[parser_type]; ok {
                var_value, err := parser(line_fields[i])
                if err != nil {
                    return nil, false
                } else {
                    values = append(values, var_value)
                }
            } else {
                panic("Unknown parser type: " + parser_type)
            }
        } else {
            if format_fields[i] != line_fields[i] {
                return nil, false
            }
        }
    }
    return values, true
}