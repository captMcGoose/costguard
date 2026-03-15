package pricing

import (
    "testing"

    "github.com/captMcGoose/costguard/internal/terraform"
)

func TestEstimateMonthlyCost(t *testing.T) {
    cases := []struct {
        name    string
        change  terraform.ResourceChange
        attrs   map[string]interface{}
        want    float64
        wantErr error
    }{
        {"aws_instance supported", terraform.ResourceChange{Type: "aws_instance"}, map[string]interface{}{"instance_type": "t3.medium"}, 30, nil},
        {"aws_db_instance supported", terraform.ResourceChange{Type: "aws_db_instance"}, map[string]interface{}{"instance_class": "db.m6i.large"}, 180, nil},
        {"aws_nat_gateway supported", terraform.ResourceChange{Type: "aws_nat_gateway"}, map[string]interface{}{}, 32, nil},
        {"aws_ebs_volume supported", terraform.ResourceChange{Type: "aws_ebs_volume"}, map[string]interface{}{"size": 100.0}, 8, nil},
        {"unsupported resource", terraform.ResourceChange{Type: "aws_lambda_function"}, map[string]interface{}{}, 0, ErrUnsupportedResource},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := EstimateMonthlyCost(tc.change, tc.attrs)
            if tc.wantErr != nil {
                if err == nil {
                    t.Fatalf("expected error %v, got nil", tc.wantErr)
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tc.want {
                t.Fatalf("expected %v, got %v", tc.want, got)
            }
        })
    }
}

func TestCalculateCostDiff(t *testing.T) {
    changes := []terraform.ResourceChange{
        {Address: "aws_db_instance.prod_db", Type: "aws_db_instance", Action: "create", Attributes: map[string]interface{}{"instance_class": "db.m6i.large"}},
        {Address: "aws_nat_gateway.main", Type: "aws_nat_gateway", Action: "create", Attributes: map[string]interface{}{}},
        {Address: "aws_ebs_volume.data", Type: "aws_ebs_volume", Action: "update", Attributes: map[string]interface{}{"size": 50.0}},
        {Address: "aws_instance.worker", Type: "aws_instance", Action: "replace", Attributes: map[string]interface{}{"instance_type": "t3.micro"}},
    }

    summary := CalculateCostDiff(changes, false)
    if summary.TotalMonthly <= 0 {
        t.Fatalf("expected positive total monthly cost, got %v", summary.TotalMonthly)
    }
    if len(summary.Drivers) == 0 {
        t.Fatal("expected non-empty drivers")
    }
    if summary.Drivers[0].MonthlyCost < summary.Drivers[len(summary.Drivers)-1].MonthlyCost {
        t.Fatal("drivers should be sorted descending")
    }
}
