package pricing

import (
    "errors"
    "fmt"
    "os"
    "sort"
    "github.com/captMcGoose/costguard/internal/terraform"
)

var (
    ErrUnsupportedResource = errors.New("unsupported resource")
    ErrMissingAttribute = errors.New("missing attribute")
)

type CostDriver struct {
    Address     string
    MonthlyCost float64
}

type CostSummary struct {
    TotalMonthly float64
    Drivers      []CostDriver
}

func EstimateMonthlyCost(change terraform.ResourceChange, attrs map[string]interface{}) (float64, error) {
    switch change.Type {
    case "aws_instance":
        rawType, ok := attrs["instance_type"]
        if !ok {
            return 0, ErrMissingAttribute
        }
        instanceType, ok := rawType.(string)
        if !ok || instanceType == "" {
            return 0, ErrMissingAttribute
        }
        price, ok := EC2Pricing[instanceType]
        if !ok {
            return 0, ErrUnsupportedResource
        }
        return price, nil
    case "aws_db_instance":
        rawClass, ok := attrs["instance_class"]
        if !ok {
            return 0, ErrMissingAttribute
        }
        instanceClass, ok := rawClass.(string)
        if !ok || instanceClass == "" {
            return 0, ErrMissingAttribute
        }
        price, ok := RDSPricing[instanceClass]
        if !ok {
            return 0, ErrUnsupportedResource
        }
        return price, nil
    case "aws_nat_gateway":
        return NATGatewayMonthly, nil
    case "aws_ebs_volume":
        rawSize, ok := attrs["size"]
        if !ok {
            return 0, ErrMissingAttribute
        }
        sizeFloat, err := toFloat64(rawSize)
        if err != nil {
            return 0, ErrMissingAttribute
        }
        return sizeFloat * EBSPricePerGB, nil
    default:
        return 0, ErrUnsupportedResource
    }
}

func CalculateCostDiff(changes []terraform.ResourceChange) CostSummary {
    summary := CostSummary{}
    drivers := make([]CostDriver, 0, len(changes))

    for _, ch := range changes {
        if ch.Action == "delete" {
            continue
        }

        cost, err := EstimateMonthlyCost(ch, ch.Attributes)
        if err != nil {
            fmt.Fprintf(os.Stderr, "skipping %s: %v\n", ch.Address, err)
            continue
        }
        if cost <= 0 {
            continue
        }

        summary.TotalMonthly += cost
        drivers = append(drivers, CostDriver{Address: ch.Address, MonthlyCost: cost})
    }

    sort.Slice(drivers, func(i, j int) bool {
        return drivers[i].MonthlyCost > drivers[j].MonthlyCost
    })
    summary.Drivers = drivers
    return summary
}

func toFloat64(v interface{}) (float64, error) {
    switch t := v.(type) {
    case float64:
        return t, nil
    case float32:
        return float64(t), nil
    case int:
        return float64(t), nil
    case int64:
        return float64(t), nil
    case uint64:
        return float64(t), nil
    case jsonNumber:
        return t.Float64()
    default:
        return 0, errors.New("unsupported numeric type")
    }
}

// jsonNumber is a minimal interface to support decoding numbers from json.Decoder.
type jsonNumber interface {
    Float64() (float64, error)
}
